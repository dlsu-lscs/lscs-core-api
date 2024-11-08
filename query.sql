-- name: GetMemberInfo :one
SELECT m.email, m.full_name, c.committee_name, d.division_name, p.position_name 
FROM members m
JOIN committees c ON m.committee_id = c.committee_id
JOIN divisions d ON c.division_id = d.division_id
JOIN positions p ON m.position_id = p.position_id
WHERE m.email = ?;

-- name: ListMembers :many
SELECT * FROM members ORDER BY email;

-- name: CheckEmailIfMember :one
SELECT email FROM members WHERE email = ?;
