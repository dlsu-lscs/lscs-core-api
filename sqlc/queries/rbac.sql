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
