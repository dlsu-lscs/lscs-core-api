# LSCS Central Auth API


we will manually add the LSCS Members as `users` to the database
- manualSaveUser()
- role will be "manually filled out"
- avatar_url will only be "automatically filled out" when they login to google

MAY SHEETS TAYO
- convert to SQL nalang gagawin
- no more registration

## TODOs

- [ ] add auto redirect to `/web` when testing with frontend 
- [ ] frontend: how receive jwt token when hit `/login`

## Route Endpoints

### Auth (MAIN)

Everything will be redirected to `/` after successful login

`/login?provider=google`
- endpoint for logging in to google
- if successful: will receive token via 
- **update: just check/validate if user is already in the database (meaning they are an LSCS member)**

`/auth/google/callback`
- google callback

`/logout`
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
