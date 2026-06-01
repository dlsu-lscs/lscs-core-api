package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/config"
	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"

	_ "github.com/dlsu-lscs/lscs-core-api/internal/helpers" // for swagger type definitions
)

const (
	googleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL = "https://oauth2.googleapis.com/token"
	googleUserURL  = "https://www.googleapis.com/oauth2/v2/userinfo"
	sessionCookie  = "session_id"
)

// OAuthHandler handles OAuth authentication for web UI
type OAuthHandler struct {
	cfg            *config.Config
	sessionService SessionService
	dbService      database.Service
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(cfg *config.Config, sessionService SessionService, dbService database.Service) *OAuthHandler {
	return &OAuthHandler{
		cfg:            cfg,
		sessionService: sessionService,
		dbService:      dbService,
	}
}

// googleUserInfo represents the response from Google's userinfo endpoint
type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// googleTokenResponse represents the response from Google's token endpoint
type googleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token,omitempty"`
}

// MeResponse represents the response for /auth/me endpoint
type MeResponse struct {
	ID            int32  `json:"id" example:"12212345"`
	Email         string `json:"email" example:"john.doe@dlsu.edu.ph"`
	FullName      string `json:"full_name" example:"John Doe"`
	Nickname      string `json:"nickname,omitempty" example:"Johnny"`
	PositionID    string `json:"position_id,omitempty" example:"VP"`
	CommitteeID   string `json:"committee_id,omitempty" example:"RND"`
	College       string `json:"college,omitempty" example:"CCS"`
	Program       string `json:"program,omitempty" example:"BSCS"`
	Discord       string `json:"discord,omitempty" example:"john#1234"`
	Interests     string `json:"interests,omitempty" example:"coding, gaming"`
	ContactNumber string `json:"contact_number,omitempty" example:"+639123456789"`
	FbLink        string `json:"fb_link,omitempty" example:"https://fb.com/johndoe"`
	Telegram      string `json:"telegram,omitempty" example:"@johndoe"`
	HouseID       int32  `json:"house_id,omitempty" example:"1"`
}

// GoogleLoginHandler initiates Google OAuth flow
// @Summary Initiate Google OAuth Login
// @Description Redirects to Google OAuth consent screen for web UI login
// @Tags auth
// @Param remember query bool false "Remember me for 30 days"
// @Param redirect query string false "URL to redirect after login"
// @Success 302 "Redirect to Google OAuth"
// @Router /auth/google/login [get]
func (h *OAuthHandler) GoogleLoginHandler(c echo.Context) error {
	rememberMe := c.QueryParam("remember") == "true"
	redirectURL := c.QueryParam("redirect")

	// build state parameter (encodes remember me and redirect URL)
	state := fmt.Sprintf("%t|%s", rememberMe, redirectURL)

	params := url.Values{}
	params.Set("client_id", h.cfg.GoogleClientID)
	params.Set("redirect_uri", h.cfg.OAuthRedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("state", state)
	params.Set("access_type", "online")
	params.Set("prompt", "select_account")

	authURL := googleAuthURL + "?" + params.Encode()
	return c.Redirect(http.StatusFound, authURL)
}

// GoogleCallbackHandler handles the OAuth callback from Google.
// @Summary Handle Google OAuth Callback
// @Description Processes Google OAuth callback, creates session, and redirects to frontend
// @Tags auth
// @Param code query string true "Authorization code from Google"
// @Param state query string false "State parameter containing remember me flag"
// @Success 302 "Redirect to frontend with session cookie"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized - not an LSCS member"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Router /auth/google/callback [get]
func (h *OAuthHandler) GoogleCallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	errorParam := c.QueryParam("error")

	if errorParam != "" {
		log.Error().Str("error", errorParam).Msg("OAuth error from Google")
		return c.Redirect(http.StatusFound, h.cfg.FrontendURL()+"/login?error=oauth_denied")
	}

	if code == "" {
		return c.Redirect(http.StatusFound, h.cfg.FrontendURL()+"/login?error=no_code")
	}

	// parse state parameter
	rememberMe := false
	redirectPath := ""
	if state != "" {
		parts := strings.SplitN(state, "|", 2)
		if len(parts) >= 1 {
			rememberMe = parts[0] == "true"
		}
		if len(parts) >= 2 {
			redirectPath = parts[1]
		}
	}

	// exchange code for token
	tokenResp, err := h.exchangeCodeForToken(c.Request().Context(), code)
	if err != nil {
		log.Error().Err(err).Msg("failed to exchange code for token")
		return c.Redirect(http.StatusFound, h.cfg.FrontendURL()+"/login?error=token_exchange")
	}

	// get user info from Google
	userInfo, err := h.getUserInfo(c.Request().Context(), tokenResp.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user info")
		return c.Redirect(http.StatusFound, h.cfg.FrontendURL()+"/login?error=user_info")
	}

	// check if user is an LSCS member
	q := repository.New(h.dbService.GetConnection())
	member, err := q.GetMemberByEmail(c.Request().Context(), userInfo.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Str("email", userInfo.Email).Msg("non-member attempted login")
			return c.Redirect(http.StatusFound, h.cfg.FrontendURL()+"/login?error=not_member")
		}
		log.Error().Err(err).Msg("failed to check member status")
		return c.Redirect(http.StatusFound, h.cfg.FrontendURL()+"/login?error=db_error")
	}

	// get client info for session
	userAgent := c.Request().UserAgent()
	ipAddress := c.RealIP()

	// create session
	session, err := h.sessionService.CreateSession(c.Request().Context(), member.ID, rememberMe, userAgent, ipAddress)
	if err != nil {
		log.Error().Err(err).Msg("failed to create session")
		return c.Redirect(http.StatusFound, h.cfg.FrontendURL()+"/login?error=session_create")
	}

	// set session cookie
	h.setSessionCookie(c, session.ID, rememberMe)

	log.Info().
		Int32("member_id", member.ID).
		Str("email", userInfo.Email).
		Bool("remember_me", rememberMe).
		Msg("user logged in")

	// redirect to frontend
	redirectTo := h.cfg.FrontendURL()
	if redirectPath != "" {
		redirectTo = redirectTo + redirectPath
	}

	return c.Redirect(http.StatusFound, redirectTo)
}

// LogoutHandler logs out the user by deleting the session
// @Summary Logout
// @Description Deletes the current session and clears the session cookie
// @Tags auth
// @Success 200 {object} map[string]string "Logged out successfully"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Router /auth/logout [post]
func (h *OAuthHandler) LogoutHandler(c echo.Context) error {
	cookie, err := c.Cookie(sessionCookie)
	if err != nil || cookie.Value == "" {
		return c.JSON(http.StatusOK, map[string]string{"message": "Already logged out"})
	}

	// delete session from database
	if err := h.sessionService.DeleteSession(c.Request().Context(), cookie.Value); err != nil {
		log.Error().Err(err).Msg("failed to delete session")
		// continue to clear cookie even if db delete fails
	}

	// clear cookie
	h.clearSessionCookie(c)

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// MeHandler returns the current user's info
// @Summary Get Current User
// @Description Returns the authenticated user's profile information
// @Tags auth
// @Produce json
// @Success 200 {object} MeResponse "User profile"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Router /auth/me [get]
func (h *OAuthHandler) MeHandler(c echo.Context) error {
	// user_id and user_email are set by session middleware
	memberID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	q := repository.New(h.dbService.GetConnection())
	member, err := q.GetMemberInfoById(c.Request().Context(), memberID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Member not found"})
		}
		log.Error().Err(err).Int32("member_id", memberID).Msg("failed to get member info")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := MeResponse{
		ID:       member.ID,
		Email:    member.Email,
		FullName: member.FullName,
	}

	if member.Nickname.Valid {
		response.Nickname = member.Nickname.String
	}
	if member.PositionID.Valid {
		response.PositionID = member.PositionID.String
	}
	if member.CommitteeID.Valid {
		response.CommitteeID = member.CommitteeID.String
	}
	if member.College.Valid {
		response.College = member.College.String
	}
	if member.Program.Valid {
		response.Program = member.Program.String
	}
	if member.Discord.Valid {
		response.Discord = member.Discord.String
	}
	if member.Interests.Valid {
		response.Interests = member.Interests.String
	}
	if member.ContactNumber.Valid {
		response.ContactNumber = member.ContactNumber.String
	}
	if member.FbLink.Valid {
		response.FbLink = member.FbLink.String
	}
	if member.Telegram.Valid {
		response.Telegram = member.Telegram.String
	}

	return c.JSON(http.StatusOK, response)
}

// exchangeCodeForToken exchanges the authorization code for an access token
func (h *OAuthHandler) exchangeCodeForToken(ctx context.Context, code string) (*googleTokenResponse, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", h.cfg.GoogleClientID)
	data.Set("client_secret", h.cfg.GoogleClientSecret)
	data.Set("redirect_uri", h.cfg.OAuthRedirectURL)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, googleTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp googleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// getUserInfo fetches user info from Google using the access token
func (h *OAuthHandler) getUserInfo(ctx context.Context, accessToken string) (*googleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute userinfo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userinfo request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo response: %w", err)
	}

	return &userInfo, nil
}

// setSessionCookie sets the session cookie with appropriate settings
func (h *OAuthHandler) setSessionCookie(c echo.Context, sessionID string, rememberMe bool) {
	maxAge := h.cfg.SessionDuration
	if rememberMe {
		maxAge = h.cfg.SessionRememberDuration
	}

	cookie := &http.Cookie{
		Name:     sessionCookie,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cfg.IsProduction(), // HTTPS only in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
	c.SetCookie(cookie)
}

// clearSessionCookie removes the session cookie
func (h *OAuthHandler) clearSessionCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cfg.IsProduction(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // delete cookie
	}
	c.SetCookie(cookie)
}
