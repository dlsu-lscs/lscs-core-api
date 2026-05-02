# LSCS Core API

The official _Core API Microservice_ of **La Salle Computer Society (LSCS)**

This is a core microservice, meant to be used by an application backend or a frontend client for authenticating LSCS Members and returning necessary data from them.

## Usage

> [!IMPORTANT]
>
> - NOTE: **only RND members can request an API key (associated with their DLSU email)** - to prevent unauthorized access
> - API-key endpoints require `Authorization: Bearer <API-KEY>`
> - Web UI endpoints require the `session_id` cookie (set after Google OAuth)
> - Only RND AVP+ can request API keys

HTTP request samples are in `requests.http` (use with VS Code REST Client or JetBrains HTTP client).

## Public Endpoints

### GET `/`

- health check
- `request`:

```http
GET {{baseUrl}}/
```

### GET `/docs/*`

- Swagger UI
- `request`:

```http
GET {{baseUrl}}/docs/
```

## Auth Endpoints

OAuth login/callback are public. Logout and profile endpoints use the `session_id` cookie.

### GET `/auth/google/login`

- initiates Google OAuth flow for web UI login
- `request`:

```http
GET {{baseUrl}}/auth/google/login?remember=true&redirect=/dashboard
```

### GET `/auth/google/callback`

- OAuth callback endpoint (normally called by Google after login)
- `request`:

```http
GET {{baseUrl}}/auth/google/callback?code=<GOOGLE_CODE>&state=true%7C%2Fdashboard
```

### POST `/auth/logout`

- clears the current session
- `request`:

```http
POST {{baseUrl}}/auth/logout
Cookie: session_id={{sessionId}}
```

### GET `/auth/me`

- returns the authenticated member profile
- `request`:

```http
GET {{baseUrl}}/auth/me
Cookie: session_id={{sessionId}}
```

### PUT `/auth/me`

- updates the authenticated member's profile (self-editable fields only)
- `request`:

```http
PUT {{baseUrl}}/auth/me
Content-Type: application/json
Cookie: session_id={{sessionId}}

{"nickname":"JD","telegram":"@jd","discord":"jd#0001"}
```

## API Key Management (Web UI - Session Cookie)

### POST `/request-key`

- request an API key (RND AVP+ only)
- `request`:

```http
POST {{baseUrl}}/request-key
Content-Type: application/json
Cookie: session_id={{sessionId}}

{
  "project": "My Awesome Project",
  "allowed_origin": "https://my-awesome-project.com",
  "is_dev": false,
  "is_admin": false
}
```

- **Request Body Fields:**
    - `project` (string, optional): A name for your project.
    - `allowed_origin` (string, optional): The URL where the key will be used. Required for production keys. Must start with `http://localhost` for dev keys if provided.
    - `is_dev` (boolean, optional): Set to `true` for a development key (for `localhost`). Defaults to `false`.
    - `is_admin` (boolean, optional): Set to `true` to create an admin key (unrestricted). Defaults to `false`.

- `response`:

```json
{ // success
    "api_key": "a_very_long_and_secure_api_key_string",
    "email": "user_from_token@dlsu.edu.ph",
    "expires_at": "2027-01-27T15:04:05Z"
}

{ // fail
  "error": "User has unauthorized position or committee"
}
```

### GET `/api-keys`

- lists API keys for the authenticated user
- `request`:

```http
GET {{baseUrl}}/api-keys
Cookie: session_id={{sessionId}}
```

### DELETE `/api-keys/{id}`

- deletes/revokes an API key by ID (must belong to the authenticated user)
- `request`:

```http
DELETE {{baseUrl}}/api-keys/123
Cookie: session_id={{sessionId}}
```

## Member Endpoints (API Key - Bearer JWT)

- all routes require `Authorization: Bearer <API-KEY>` in the request headers

### GET `/members`

- returns members from database
- optional query params:
    - `position`: one or more `position_id` values (comma-separated or repeated)
    - `committee`: one or more `committee_id` values (comma-separated or repeated)
- `request`:

```http
GET {{baseUrl}}/members?position=MEM,EVP&committee=RND&committee=HRD
Authorization: Bearer {{apiKey}}
```

- `response`:

```json
[
    {
        "id": 12312312,
        "full_name": "Hehe E. Hihi",
        "nickname": null,
        "email": "hehe_hihi@dlsu.edu.ph",
        "telegram": null,
        "position_id": "MEM",
        "committee_id": "RND",
        "college": "CCS",
        "program": "BSCS",
        "discord": null,
        "interests": null,
        "contact_number": null,
        "fb_link": null,
        "image_url": null,
        "house_name": "Gell-Mann"
    }
]
```

### GET `/committees`

- returns all committees
- `request`:

```http
GET {{baseUrl}}/committees
Authorization: Bearer {{apiKey}}
```

- `response`:

```json
{
    "committees": [
        {
            "committee_id": "RND",
            "committee_name": "Research and Development",
            "committee_head": 12312312,
            "division_id": "INT"
        }
    ]
}
```

### POST `/member`

- returns detailed member info by email
- `request`:

```http
POST {{baseUrl}}/member
Authorization: Bearer {{apiKey}}
Content-Type: application/json

{"email": "edwin_sadiarinjr@dlsu.edu.ph"}
```

- `response`:

```json
{
    "id": 12323004,
    "committee_id": "RND",
    "committee_name": "Research and Development",
    "division_id": "INT",
    "division_name": "Internals",
    "email": "edwin_sadiarinjr@dlsu.edu.ph",
    "full_name": "Edwin Sadiarin Jr.",
    "position_id": "MEM",
    "position_name": "Committee Trainee",
    "house_name": "Gell-Mann"
}
```

### POST `/member-id`

- returns detailed member info by ID
- `request`:

```http
POST {{baseUrl}}/member-id
Authorization: Bearer {{apiKey}}
Content-Type: application/json

{"id": 12323004}
```

- `response`:

```json
{
    "id": 12323004,
    "committee_id": "RND",
    "committee_name": "Research and Development",
    "division_id": "INT",
    "division_name": "Internals",
    "email": "edwin_sadiarinjr@dlsu.edu.ph",
    "full_name": "Edwin Sadiarin Jr.",
    "position_id": "MEM",
    "position_name": "Committee Trainee",
    "house_name": "Gell-Mann"
}
```

### POST `/check-email`

- checks if the email exists in database (indicating if it is an LSCS member or not)
- `request`:

```http
POST {{baseUrl}}/check-email
Authorization: Bearer {{apiKey}}
Content-Type: application/json

{"email": "edwin_sadiarinjr@dlsu.edu.ph"}
```

- `response`:

```json
{
    // success
    "email": "edwin_sadiarinjr@dlsu.edu.ph",
    "state": "present",
    "success": "Email is an LSCS member"
}
```

### POST `/check-id`

- checks if the provided id exists in database (indicating if it is an LSCS member or not)
- requires `id` in the request body

> [!IMPORTANT]
> **MAKE SURE to send the `id` as an int (in the request body)**

- `request`:

```http
POST {{baseUrl}}/check-id
Authorization: Bearer {{apiKey}}
Content-Type: application/json

{"id": 12323004}
```

- `response`:

```json
{
    // success
    "id": 12323004,
    "state": "present",
    "success": "ID is an LSCS member"
}
```

## Member Endpoints (Web UI - Session Cookie)

### GET `/auth/members/{id}`

- returns detailed member info by ID (session-protected)
- `request`:

```http
GET {{baseUrl}}/auth/members/12323004
Cookie: session_id={{sessionId}}
```

### PUT `/auth/members/{id}`

- updates a member profile (requires authorization to edit the target member)
- `request`:

```http
PUT {{baseUrl}}/auth/members/12323004
Content-Type: application/json
Cookie: session_id={{sessionId}}

{"full_name":"Updated Name","position_id":"MEM","committee_id":"RND"}
```

## Upload Endpoints (Web UI - Session Cookie)

### POST `/upload/profile-image`

- generates a pre-signed upload URL
- `request`:

```http
POST {{baseUrl}}/upload/profile-image
Content-Type: application/json
Cookie: session_id={{sessionId}}

{"content_type":"image/png"}
```

### POST `/upload/profile-image/complete`

- confirms upload complete and updates `image_url`
- `request`:

```http
POST {{baseUrl}}/upload/profile-image/complete
Content-Type: application/json
Cookie: session_id={{sessionId}}

{"object_key":"profile-images/12323004/abc.png"}
```

### DELETE `/upload/profile-image`

- deletes the current profile image
- `request`:

```http
DELETE {{baseUrl}}/upload/profile-image
Cookie: session_id={{sessionId}}
```

## Contributing

### Deployment

This project uses a monorepo structure with separate CI/CD pipelines for the API and Web services.

#### Architecture

- **API Service**: Go/Echo backend running on port 8080
- **Web Service**: Next.js frontend running on port 3000
- **Deployment Platform**: Dokploy (self-hosted VPS)

#### CI/CD Pipeline

The GitHub Actions workflows handle testing, building, and deployment:

| Workflow                 | Trigger                                                    | Purpose                           |
| ------------------------ | ---------------------------------------------------------- | --------------------------------- |
| `001-test.yml`           | Push to Go or Web files                                    | Run tests and linting             |
| `002-build-push-api.yml` | Changes to `**/*.go`, `go.mod`, `go.sum`, `Dockerfile.api` | Build & push API Docker image     |
| `003-build-push-web.yml` | Changes to `web/**`, `web/Dockerfile`                      | Build & push Web Docker image     |
| `004-deploy-api.yml`     | After `002-build-push-api.yml` completes                   | Trigger API deployment in Dokploy |
| `005-deploy-web.yml`     | After `003-build-push-web.yml` completes                   | Trigger Web deployment in Dokploy |

#### Docker Images

- **API Image**: `ghcr.io/<org>/lscs-core-api-api:<tag>`
- **Web Image**: `ghcr.io/<org>/lscs-core-api-web:<tag>`

#### Dokploy Setup

1. Create two applications in Dokploy:
    - **API**: Points to `Dockerfile.api`, port 8080
    - **Web**: Points to `web/Dockerfile`, port 3000

2. Configure environment variables in Dokploy for each application

3. Obtain webhook URLs and tokens from Dokploy, then add to GitHub secrets:
    - `DOKPLOY_API_WEBHOOK_URL`
    - `DOKPLOY_API_TOKEN`
    - `DOKPLOY_WEB_WEBHOOK_URL`
    - `DOKPLOY_WEB_TOKEN`

#### Selective Deployment

Changes are automatically isolated:

- Go code changes → Only rebuilds and deploys API
- Web code changes → Only rebuilds and deploys Web
- Configuration changes → Rebuilds and deploys both

#### Security Scanning

Both Docker images are scanned using Trivy for vulnerabilities. Critical and High severity issues are reported in the GitHub Security tab.

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

TMP
Given the state of the @PLAN.md , formulate todos for phase 3. Check if sessions table is created (i think it is) in migrations.

Don't delete existing functionality, especially the queries in @query.sql since that is used in the business logic @internal/member/ . Check @PLAN.md @AGENTS.md for more context.

Also i noticed the in @schema.sql the session table is not there. Can we first refactor the sqlc-related files, so the folder structure to be like @sqlc/queries/ also edit the @sqlc.yaml for this structure. Now how about the schema? how should i move the @schema.sql ? Can i point it to migrations/ as schema? Formulate a plan first for this
