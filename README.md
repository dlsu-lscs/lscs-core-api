# LSCS Core API

The official *Core API Microservice* of **La Salle Computer Society (LSCS)**

This is a core microservice, meant to be used by an application backend or a frontend client for authenticating LSCS Members and returning necessary data from them.

## Usage

> [!IMPORTANT]
> NOTE: **only RND members can request an API key (associated with their DLSU email)** - to prevent unauthorized access
> 
> --> All endpoints now requires an API key (in the `Authorization` request headers)


## Auth Endpoints

### POST `/request-key`

- **Requires Google Authentication.** You must provide a valid Google ID token in the `Authorization: Bearer <ID_TOKEN>` header.
- The request body now specifies the key's properties.

- `request`:
```bash
curl -X POST https://core.api.dlsu-lscs.org/request-key \
  -H "Authorization: Bearer <GOOGLE_ID_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
        "project": "My Awesome Project",
        "allowed_origin": "https://my-awesome-project.com",
        "is_dev": false,
        "is_admin": false
      }'
```

- **Request Body Fields:**
    *   `project` (string, required): A name for your project.
    *   `allowed_origin` (string, optional): The URL where the key will be used. Required for production keys. Must start with `http://localhost` for dev keys if provided.
    *   `is_dev` (boolean, optional): Set to `true` for a development key (for `localhost`). Defaults to `false`.
    *   `is_admin` (boolean, optional): Set to `true` to create an admin key (unrestricted). Defaults to `false`.

- `response`:
```json
{ // success
    "api_key": "a_very_long_and_secure_api_key_string",
    "email": "user_from_token@dlsu.edu.ph"
}

{ // fail
  "error": "User is not a member of Research and Development"
}
```

### POST `/revoke-key`

- for revoking/deleting key
- requires `email` and `pepper` in the request body

- `request`:
```bash
curl -X POST https://core.api.dlsu-lscs.org/revoke-key \
  -H "Content-Type: application/json" \
  -d '{"email": "edwin_sadiarinjr@dlsu.edu.ph", "pepper": "<CONTACT_ADMIN_DEVELOPER_TO_REVOKE_KEY>"}'
```

- `response`:
```
API key for <email> is successfully revoked
```


## Member Endpoints

- all routes: requires `Authorization: Bearer <API-KEY>` in the request headers

### GET `/members`

- returns all LSCS members from database (*yes*)
- requires `Authorization: Bearer <API-KEY>` in the request headers

- `request`:
```bash
curl -X GET https://core.api.dlsu-lscs.org/members \
  -H "Authorization: Bearer <API-KEY>"
```

- `response`:
```json
[
   {
     "id": 12312312,
     "full_name": "Hehe E. Hihi",
     "nickname": "Huhi",
     "email": "hehe_hihi@dlsu.edu.ph",
     "telegram": "",
     "position_id": "MEM",
     "committee_id": "MEM",
     "college": "CCS",
     "program": "BS-Org",
     "discord": ""
   },
   {
     "id": 11111110,
     "full_name": "Peter Parker",
     "nickname": "Peter",
     "email": "peter_parker@dlsu.edu.ph",
     "telegram": "@something",
     "position_id": "MEM",
     "committee_id": "MEM",
     "college": "CLA",
     "program": "POM-LGL",
     "discord": ""
   }
]
```

### GET `/committees`

- returns all committees
- requires `Authorization: Bearer <API-KEY>` in the request headers

- `request`:
```bash
curl -X GET https://core.api.dlsu-lscs.org/committees \
  -H "Authorization: Bearer <API-KEY>"
```

- `response`:
```json
{
    "committees": [...] 
}
```

### POST `/member`

- returns `email`, `full_name`, `committee_name`, `position_name`, `division_name`, `committee_id`, and `division_id` of the LSCS member 
- requires `Authorization: Bearer <API-KEY>` in the request headers
- requires `email` in the request body

- `request`:
```bash
curl -X POST https://core.api.dlsu-lscs.org/member \
  -H "Authorization: Bearer <API-KEY>" \
  -H "Content-Type: application/json" \
  -d '{"email": "edwin_sadiarinjr@dlsu.edu.ph"}'
```

- `response`:
```json
{ // success
  "committee_id": "RND",
  "committee_name": "Research and Development",
  "division_id": "INT",
  "division_name": "Internals",
  "email": "edwin_sadiarinjr@dlsu.edu.ph",
  "full_name": "Edwin Sadiarin Jr.",
  "position_name": "Committee Trainee"
}

{ // fail
  "error": "Email is not an LSCS member"
}
```

### POST `/member-id`

- returns `id`, `email`, `full_name`, `committee_name`, `position_name`, `division_name`, `committee_id`, and `division_id` of the LSCS member 
- requires `Authorization: Bearer <API-KEY>` in the request headers
- requires `id` in the request body

- `request`:
```bash
curl -X POST https://core.api.dlsu-lscs.org/member-id \
  -H "Authorization: Bearer <API-KEY>" \
  -H "Content-Type: application/json" \
  -d '{"id": 12323004}'
```

- `response`:
```json
{ // success
  "id": 12323004,
  "committee_id": "RND",
  "committee_name": "Research and Development",
  "division_id": "INT",
  "division_name": "Internals",
  "email": "edwin_sadiarinjr@dlsu.edu.ph",
  "full_name": "Edwin Sadiarin Jr.",
  "position_name": "Committee Trainee"
}

{ // fail
  "error": "ID is not an LSCS member"
}
```

### POST `/check-email`

- checks if the email exists in database (indicating if it is an LSCS member or not)
- requires `Authorization: Bearer <API-KEY>` in the request headers
- requires `email` in the request body

- `request`:
```bash
curl -X POST https://core.api.dlsu-lscs.org/check-email \
  -H "Authorization: Bearer <API-KEY>" \
  -H "Content-Type: application/json" \
  -d '{"email": "edwin_sadiarinjr@dlsu.edu.ph"}'
```

- `response`:
```json
{ // success
  "email": "edwin_sadiarinjr@dlsu.edu.ph",
  "state": "present",
  "success": "Email is an LSCS member"
}

{ // fail
  "email": "test@dlsu.edu.ph",
  "error": "Not an LSCS member",
  "state": "absent"
}
```

### POST `/check-id`

- checks if the provided id exists in database (indicating if it is an LSCS member or not)
- requires `Authorization: Bearer <API-KEY>` in the request headers
- requires `id` in the request body

> [!IMPORTANT]
> **MAKE SURE to send the `id` as an int (in the request body)**

- `request`:
```bash
curl -X POST https://core.api.dlsu-lscs.org/check-id \
  -H "Authorization: Bearer <API-KEY>" \
  -H "Content-Type: application/json" \
  -d '{"id": 12323004}'
```

- `response`:
```json
{ // success
    "id": 12323004,
    "state": "present",
    "success": "ID is an LSCS member"
}

{ // fail
    "id": 1231434214,
    "state": "absent"
    "error": "Not an LSCS member",
}
```

## Contributing

### (for Maintainers & Admins) Creating a Release

To create a new release, you need to push a new tag to the repository. The tag must follow the semantic versioning format (e.g., `v1.2.3`).

1.  **Create a new tag:**
    ```bash
    git tag v1.2.4
    ```

2.  **Push the tag to the repository:**
    ```bash
    git push origin v1.2.4
    ```

Pushing a new tag will trigger the `release` workflow, which will automatically:
- Build the binaries for different operating systems.
- Create a new release on GitHub.
- Upload the binaries as release assets.
- Include a link to the corresponding Docker image in the release notes.
