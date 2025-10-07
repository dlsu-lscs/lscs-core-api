# LSCS Core API

The official *Core API Microservice* of **La Salle Computer Society (LSCS)**

This is a core microservice, meant to be used by an application backend or a frontend client for authenticating LSCS Members and returning necessary data from them.

## Usage

> [!IMPORTANT]
> NOTE: **only RND members can request an API key (associated with their DLSU email)** - to prevent unauthorized access
> 
> --> All endpoints now requires an API key (in the `Authorization` request headers)


## Auth Endpoints

- `[UPDATE 20251007]`: **no longer needs to login to Google**

### POST `/request-key`

- requires `email` in the request body

- `request`:
```bash
curl -X POST https://core.api.dlsu-lscs.org/request-key \
  -H "Content-Type: application/json" \
  -d '{"email": "edwin_sadiarinjr@dlsu.edu.ph"}'
```

- `response`:
```json
{ // success
    "api_key": "somethingsomethingstringssdgasdfkgjdsf",
    "email": "edwin_sadiarinjr@dlsu.edu.ph"
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
