# LSCS Central Auth API

The official *Authentication Microservice* of **La Salle Computer Society (LSCS)**

## Usage

This is an auth microservice, meant to be used by an application backend.

*Treat this like a "middleware" with extra features, only for authenticating LSCS Members and returning necessary data from them.*


## Auth Endpoints

### GET `/authenticate?provider=google`

- the main authentication endpoint that will redirect users to google authentication page

- example `response` (assuming already authenticated with Google):
```json
{ // success
  "email": "edwin_sadiarinjr@dlsu.edu.ph",
  "access_token": "examplejwttokenskibidiiwanttosleepnowits5inthemorning",
  "refresh_token": "refreshtokenexamplegimmetwicealbumpls",
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

- returns all LSCS members from database

### POST `/check-email`

- checks if the email exists in database (indicating if it is an LSCS member or not)

- example `request`:
```json
{
    "email": "edwin_sadiarinjr@dlsu.edu.ph"
}
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



## JWT

flow:
- client req to `/login` endpoint
-  `/login` endpoint -> redirects to google oauth2 login page
-  google
- hit backend
- generate JWT with Claims
-  response with HttpOnly Cookie
-  for every request, decode JWT using JWT asm/sym secret

needs:
- signing key
- signing key
- algo


UPDATE:
get data of officers from API (officersdb)
then store in central db here for auth
update User schema to match the officersdb
add endpoints for validating users:
- /validate - check user email exists in database
- /login - match the officersdb
- add redirect to /login?provider=google for everything /login
  - invalidate query param if not google
- /refresh-token - should be logged in first


# Database Schema

