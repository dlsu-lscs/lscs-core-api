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
