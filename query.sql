-- name: GetMemberInfo :one
SELECT 
  m.id, m.email, m.full_name, m.nickname, 
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
  m.id, m.email, m.full_name, m.nickname, 
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
    h.name as house_name
FROM members m
LEFT JOIN houses h ON m.house_id = h.id
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