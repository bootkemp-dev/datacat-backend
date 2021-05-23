# Datacat Backend

# Endpoints

## Authentication
### Register:
Description:
```
Registers a new user and inserts the data into the database.
```
Method: POST <br/>
Protected: NO <br/>
Endpoint:
```
/api/v1/auth/register
```
Payload:
```
{
    "username": "test-username",
    "email": "test-email",
    "password1": "password1",
    "password2": "password2"
}
```
Response:
```
if data validation fails: HTTP 400
if server fails: HTTP 500
if everything ok: HTTP 200
```

### Login
Description:
```
Logs in the user, sets httpOnly cookie with the jwt token and user data
```
Method: POST <br/>
Protected: NO <br/>
Endpoint:
```
/api/v1/auth/login
```
Payload:
```
{
    "username": "test-username",
    "password": "test-password"
}
```
Response: 
```
if data validation fails: HTTP 400
if user does not exist: HTTP 406
if server fails: HTTP 500
if everything ok: HTTP 200
```

### Me
```
Returns information about currently logged user.
```
Method: GET <br/>
Protected: YES - login required <br/>
Payload: None
Response:
```
if server fails: HTTP 500
if everything ok HTTP 200 and payload with data:
{
    "id": 1,
    "username": "test-username",
    "email": "test@email.com",
    "last-login": "01.01.1990"
}
```
