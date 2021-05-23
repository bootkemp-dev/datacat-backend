# Datacat Backend

# Endpoints

## Authentication
### Register:
Description:
```
Registers a new user and inserts the data into the database.
```
Method: POST
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
if data validation fails: HTTP 401
if server fails: HTTP 500
if everything ok: HTTP 200
```

