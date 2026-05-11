-- +goose Up
-- +goose StatementBegin

UPDATE members SET image_url = REPLACE(image_url, 's3.api.dlsu-lscs.org/lscs-core/', 'lscs-core.web.dlsu-lscs.org/')
WHERE image_url IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE members SET image_url = REPLACE(image_url, 'lscs-core.web.dlsu-lscs.org/', 's3.api.dlsu-lscs.org/lscs-core/')
WHERE image_url IS NOT NULL;

-- +goose StatementEnd
