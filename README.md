# LSCS Central Auth API

The official *Authentication Microservice* of **La Salle Computer Society (LSCS)**

## Usage

This is an auth microservice, meant to be used by an application backend.

*Treat this as a service that simply returns a JSON payload, used only for authenticating LSCS Members and returning necessary data from them.*


## Auth Endpoints

### GET `/authenticate?provider=google`

- the main authentication endpoint that will redirect users to google authentication page
    - NOTE: only emails with `@dlsu.edu.ph` domain are able to authenticate.

- example `response` (assuming already authenticated with Google):
```json
{ // success
  "email": "edwin_sadiarinjr@dlsu.edu.ph",
  "access_token": "examplejwttokenskibidiiwanttosleepnowits5inthemorning",
  "refresh_token": "refreshtokenexamplegimmetwicealbumpls",
  "member_info":   {...},
  "google_info":   {...},
  "success": "Email is an LSCS member"
  "state": "present",
}

{ // fail
  "email": "test@dlsu.edu.ph",
  "error": "Not an LSCS member",
  "state": "absent"
}

```


### POST `/invalidate`

- for logging out


## Admin Routes

### POST `/refresh-token`

- used for requesting new access tokens using existing refresh-token

### GET `/members`

- returns all LSCS members from database (*yes*)

### POST `/check-email`

- checks if the email exists in database (indicating if it is an LSCS member or not)

- example `request`:
```bash
curl -X POST http://localhost:42069/check-email \
  -H "Content-Type: application/json" \
  -d '{"email": "edwin_sadiarinjr@dlsu.edu.ph"}'

# in JSON (request):
# {
#   "email": "edwin_sadiarinjr@dlsu.edu.ph",
# }
```
- example `response`:
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
