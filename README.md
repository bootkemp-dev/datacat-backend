# Datacat Backend

## Endpoints

## Monitoring

### Add Job

```
Inserts new job to the database, adds it to the job pool and runs
```

Method: POST <br/>
Protected: YES - login required <br/>
Endpoint:

```
/api/v1/protected/add-job
```

Payload:

```
{
    "name": "test-name",
    "url": "https://test.com",
    "frequency": 300
}
```

Response:

```
if server fails: HTTP 500
if does not pass validation: HTTP 400
if not authenticated: HTTP 401
if everything ok: HTTP 200
{
    "name": "new job name"
    "frequency": 1000,
    "status": up
}
```
