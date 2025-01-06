# LSCS Central Auth API

The official *Authentication Microservice* of **La Salle Computer Society (LSCS)**

This is an auth microservice, meant to be used by an application backend.

_**Treat this as a service that simply returns a JSON payload, used only for authenticating LSCS Members and returning necessary data from them.**_

## Usage

> [!IMPORTANT]
> 1. **only RND members can request an API key (associated with their DLSU email)** - to prevent unauthorized access
> 2. **1 API key per DLSU email** - to prevent impersonation and duplicate API keys 


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


## Member Information Routes

### GET `/members`

- returns all LSCS members from database (*yes*)

### POST `/member`

- returns `email`, `full_name`, `committee_name`, `position_name`, and `division_name` of the LSCS member 
- example `request`:
```bash
curl -X POST http://localhost:42069/member \
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
  "committee_name": "Research and Development",
  "division_name": "Internals",
  "email": "edwin_sadiarinjr@dlsu.edu.ph",
  "full_name": "Edwin Sadiarin Jr.",
  "position_name": "Committee Trainee"
}

{ // fail
  "error": "Email is not an LSCS member"
}
```

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


## Admin Routes

### POST `/refresh-token`

- used for requesting new access tokens using existing refresh-token

