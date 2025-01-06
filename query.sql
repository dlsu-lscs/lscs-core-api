-- name: GetMemberInfo :one
SELECT m.email, m.full_name, c.committee_name, d.division_name, p.position_name 
FROM members m
JOIN committees c ON m.committee_id = c.committee_id
JOIN divisions d ON c.division_id = d.division_id
JOIN positions p ON m.position_id = p.position_id
WHERE m.email = ?;

-- name: GetFullMemberInfo :one
SELECT 
    m.id, m.email, m.full_name, m.nickname,
    c.committee_id, c.committee_name,
    d.division_id, d.division_name,
    p.position_id, p.position_name 
FROM members m
JOIN committees c ON m.committee_id = c.committee_id
JOIN divisions d ON c.division_id = d.division_id
JOIN positions p ON m.position_id = p.position_id
WHERE m.email = ?;

-- name: ListMembers :many
SELECT * FROM members ORDER BY email;

-- name: CheckEmailIfMember :one
SELECT email FROM members WHERE email = ?;

-- name: GetAllCommittees :many
SELECT * FROM committees;

-- name: StoreAPIKey :exec
INSERT INTO api_keys (member_email, api_key_hash, expires_at) VALUES (?, ?, ?);

-- name: GetAPIKeyInfo :one
SELECT * FROM api_keys WHERE api_key_hash = ?;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys WHERE member_email = ?;

-- name: GetAllAPIKeyHashes :many
SELECT api_key_hash FROM api_keys;
