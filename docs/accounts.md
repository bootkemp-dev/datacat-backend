# Accounts

## Reset Password
```
Generates a one time token of previously set length and it's expiration date. Then updates row in the users table on given username and sends the link to reset the password on users email.
```

Method: POST <br/>
Protected: NO <br/>
Endpoint:
```
/api/v1/accounts/reset_password?username=some-username
```

Response:
```
if username query not set: HTTP 400
if username does not exist in the database: HTTP 404
if server fails: HTTP 500
if everything OK: HTTP 200
```
## Reset Password Token Validation
```
Checks if provided token belongs to the right username, and checks if it's valid
```
Method: GET <br/>
Protected: NO <br/>
Endpoint:
```
/api/v1/accounts/reset_password?username=some-username&token=XkMymLZLPacCehDk7aHPbnm5xAQDqS
```

Response:
```
if username or token not set: HTTP 400
if token not found in the database: HTTP 401
if server fails HTTP 500
if token is expired: HTTP 406
if everything OK: HTTP 200
```

## Update Password
```
Updates users password in the database
```
Method: PUT <br/>
Protected: NO <br/>
Endpoint:
```
/api/v1/accounts/update_password?username=some-username&token=XkMymLZLPacCehDk7aHPbnm5xAQDqS
```
Request JSON:
```
{
    "password1": "new-password-1",
    "password2": "new-password-1"
}
```

Response:
```
if username or token not set: HTTP 400
if passwords from the request do not match: HTTP 400
if token not found in the database: HTTP 401
if server fails HTTP 500
if token is expired: HTTP 406
if everything OK: HTTP 200
```

