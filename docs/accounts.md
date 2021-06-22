# Accounts

## Reset Password
```
Generates a one time token of previously set length and it's expiration date. Then updates row in the users table on given username and sends the link to reset the password on users email.
```

Method: POST <br/>
Protected: NO
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

## Update Password

