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

-- name: GetMemberByEmail :one
SELECT id, email, full_name, nickname, position_id, committee_id, college, program,
       discord, interests, contact_number, fb_link, telegram, house_id, image_url
FROM members WHERE email = ?;

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
