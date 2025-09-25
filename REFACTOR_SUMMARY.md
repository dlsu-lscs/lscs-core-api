# Refactor Summary

This document summarizes the refactoring changes made to the LSCS Core API.

## 1. API Renaming

*   The API has been renamed from **LSCS Central Auth API** to **LSCS Core API**.
*   The Go module name has been updated to `github.com/dlsu-lscs/lscs-core-api`.
*   All import paths have been updated to reflect the new module name.
*   The `README.md` file has been updated with the new name and description.

## 2. Frontend Compatibility

*   Added CORS (Cross-Origin Resource Sharing) middleware to allow requests from different origins, such as web frontends.
*   The CORS middleware is configured to allow all origins, but can be restricted to specific domains for production.

## 3. Authentication

*   The authentication mechanism has been refactored to use JWT (JSON Web Tokens).
*   The old API key generation and revocation system has been removed.
*   A new endpoint, `/auth/google/callback`, has been added to handle authentication via Google Sign-In.
*   This endpoint receives a user's email, verifies if they are an LSCS member, and returns a JWT if they are.
*   A new JWT authentication middleware has been created to protect routes. This middleware checks for a valid JWT in the `Authorization` header.

## 4. Code Cleanup

*   Removed the unused `internal/middlewares/admin.go` file.
*   Removed the `cmd/api/main.go.bak` backup file.

## 5. Documentation

*   The `GEMINI.md` file has been updated to reflect all the refactoring changes, including the new project overview, build and run instructions, and development conventions.
