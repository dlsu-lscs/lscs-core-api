-- name: GetMemberInfo :one
SELECT
  m.id, m.email, m.full_name, m.nickname, m.image_url,
  c.committee_id, c.committee_name,
  d.division_id, d.division_name,
  p.position_id, p.position_name,
  h.name as house_name,
  m.contact_number, m.college, m.program,
  m.interests, m.discord, m.fb_link, m.telegram
FROM members m
LEFT JOIN committees c ON m.committee_id = c.committee_id
LEFT JOIN divisions d ON c.division_id = d.division_id
LEFT JOIN positions p ON m.position_id = p.position_id
LEFT JOIN houses h ON m.house_id = h.id
WHERE m.email = ?;

-- name: GetMemberInfoById :one
SELECT
  m.id, m.email, m.full_name, m.nickname, m.image_url,
  c.committee_id, c.committee_name,
  d.division_id, d.division_name,
  p.position_id, p.position_name,
  h.name as house_name,
  m.contact_number, m.college, m.program,
  m.interests, m.discord, m.fb_link, m.telegram
FROM members m
LEFT JOIN committees c ON m.committee_id = c.committee_id
LEFT JOIN divisions d ON c.division_id = d.division_id
LEFT JOIN positions p ON m.position_id = p.position_id
LEFT JOIN houses h ON m.house_id = h.id
WHERE m.id = ?;

-- name: ListMembers :many
SELECT
    m.id,
    m.full_name,
    m.nickname,
    m.email,
    m.telegram,
    m.position_id,
    m.committee_id,
    m.college,
    m.program,
    m.discord,
    m.interests,
    m.contact_number,
    m.fb_link,
    m.image_url,
    h.name as house_name
FROM members m
LEFT JOIN houses h ON m.house_id = h.id
WHERE (sqlc.arg(position_ids) = '' OR FIND_IN_SET(m.position_id, sqlc.arg(position_ids)))
  AND (sqlc.arg(committee_ids) = '' OR FIND_IN_SET(m.committee_id, sqlc.arg(committee_ids)))
ORDER BY m.email;

-- name: CheckEmailIfMember :one
SELECT email FROM members WHERE email = ?;

-- name: CheckIdIfMember :one
SELECT id FROM members WHERE id = ?;

-- name: GetAllCommittees :many
SELECT c.committee_id, c.committee_name, c.committee_head, c.division_id FROM committees c;

-- name: GetAllDivisions :many
SELECT d.division_id, d.division_name, d.division_head FROM divisions d;

-- name: StoreAPIKey :exec
INSERT INTO api_keys (
    member_email,
    api_key_hash,
    project,
    allowed_origin,
    is_dev,
    is_admin,
    expires_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
);

-- name: GetAPIKeyInfo :one
SELECT api_key_id, member_email, api_key_hash, project, allowed_origin, is_dev, is_admin, created_at, expires_at FROM api_keys WHERE api_key_hash = ?;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys WHERE member_email = ? LIMIT 1;

-- name: GetAllAPIKeyHashes :many
SELECT api_key_hash FROM api_keys;

-- name: GetAPIKeyInfoWithEmail :one
SELECT api_key_id, member_email, api_key_hash, project, allowed_origin, is_dev, is_admin, created_at, expires_at FROM api_keys WHERE member_email = ?;

-- name: GetEmailsInAPIKey :many
SELECT member_email FROM api_keys;

-- name: CheckAllowedOriginExists :one
SELECT EXISTS(SELECT 1 FROM api_keys WHERE allowed_origin = ? AND is_dev = false);

-- name: ListAPIKeysByEmail :many
SELECT api_key_id, member_email, project, allowed_origin, is_dev, is_admin, created_at, expires_at FROM api_keys WHERE member_email = ? ORDER BY created_at DESC;

-- name: DeleteAPIKeyById :exec
DELETE FROM api_keys WHERE api_key_id = ? AND member_email = ?;

-- Session queries for web UI authentication

-- name: CreateSession :exec
INSERT INTO sessions (id, member_id, expires_at, user_agent, ip_address)
VALUES (?, ?, ?, ?, ?);

-- name: GetSession :one
SELECT id, member_id, created_at, expires_at, last_activity, user_agent, ip_address
FROM sessions WHERE id = ? AND expires_at > NOW();

-- name: GetSessionWithMember :one
SELECT 
    s.id, s.member_id, s.created_at, s.expires_at, s.last_activity, s.user_agent, s.ip_address,
    m.email, m.full_name
FROM sessions s
JOIN members m ON s.member_id = m.id
WHERE s.id = ? AND s.expires_at > NOW();

-- name: UpdateSessionActivity :exec
UPDATE sessions SET last_activity = NOW() WHERE id = ?;

-- name: ExtendSession :exec
UPDATE sessions SET expires_at = ?, last_activity = NOW() WHERE id = ?;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = ?;

-- name: DeleteAllSessionsForMember :exec
DELETE FROM sessions WHERE member_id = ?;

-- name: CleanupExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < NOW();

-- name: GetMemberByEmail :one
SELECT id, email, full_name, nickname, position_id, committee_id, college, program,
       discord, interests, contact_number, fb_link, telegram, house_id, image_url
FROM members WHERE email = ?;

-- Member profile update queries

-- name: UpdateMemberSelf :exec
UPDATE members SET
    nickname = ?,
    telegram = ?,
    discord = ?,
    interests = ?,
    contact_number = ?,
    fb_link = ?,
    image_url = ?
WHERE id = ?;

-- name: UpdateMemberById :exec
UPDATE members SET
    full_name = ?,
    nickname = ?,
    email = ?,
    position_id = ?,
    committee_id = ?,
    college = ?,
    program = ?,
    house_id = ?,
    telegram = ?,
    discord = ?,
    interests = ?,
    contact_number = ?,
    fb_link = ?,
    image_url = ?
WHERE id = ?;

-- RBAC: Role queries

-- name: GetAllRoles :many
SELECT id, name, description FROM roles ORDER BY id;

-- name: GetRoleById :one
SELECT id, name, description FROM roles WHERE id = ?;

-- name: GetMemberRoles :many
SELECT r.id, r.name, r.description, mr.granted_by, mr.granted_at
FROM member_roles mr
JOIN roles r ON mr.role_id = r.id
WHERE mr.member_id = ?;

-- name: HasRole :one
SELECT EXISTS(SELECT 1 FROM member_roles WHERE member_id = ? AND role_id = ?);

-- name: GrantRole :exec
INSERT INTO member_roles (member_id, role_id, granted_by) VALUES (?, ?, ?);

-- name: RevokeRole :exec
DELETE FROM member_roles WHERE member_id = ? AND role_id = ?;

-- name: GetMembersWithRole :many
SELECT m.id, m.email, m.full_name, m.position_id, m.committee_id, mr.granted_at
FROM member_roles mr
JOIN members m ON mr.member_id = m.id
WHERE mr.role_id = ?
ORDER BY mr.granted_at DESC;

-- name: IsAdmin :one
SELECT EXISTS(SELECT 1 FROM member_roles WHERE member_id = ? AND role_id = 'ADMIN');

-- name: GetMemberAuthInfo :one
-- lightweight query for authorization checks (no image_url dependency)
SELECT id, position_id, committee_id FROM members WHERE email = ?;
