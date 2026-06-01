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
