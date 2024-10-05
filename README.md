# LSCS Central Auth API

## TODOs

- [x] feat: testing db sqlite (for dev)
- [ ] feat: port database to MongoDB or PostgreSQL (for production)
- [ ] feat: add token generation
- [ ] chore: update to es6
- [ ] build: dockerize

## Route Endpoints

Everything will be redirected to `/` after successful login

`/login?provider=google`
- endpoint for logging in to google
- if successful: will receive token via 
- **update: just check/validate if user is already in the database (meaning they are an LSCS member)**

`/auth/google/callback`
- google callback

`/invalidate`
- for logging out

`/refresh`
- for refresh tokens request


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
