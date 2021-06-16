# Monitoring

## Add Job

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

## Get Jobs
```
Gets users jobs from the database, both active and inactive
```

Method: GET <br/>
Protected: YES - login required <br/>
Endpoint:
```
api/v1/protected/jobs
```

Response:
```
if server fails: HTTP 500
if not authenticated: HTTP 401
if everything ok: HTTP 200
{
    "success": true,
    "jobs": {
        
    }
}
```

## Get Job Status
```
Returns current status of the job with id given as a param
```

Method: GET <br/>
Protected: YES - login required <br/>
Endpoint:
```
api/v1/protected/job/:id/status
```

Response:
```
if server fails: HTTP 500
if id is not valid: HTTP 406
if not authenticated: HTTP 401
if job with provided ID does not exist: HTTP 404
if everything ok: HTTP 200
{
   "success": true,
   "status": "job-status"
}
```

## Get Job Active
```
Returns active attribute  of the job with id given as a param
```

Method: GET <br/>
Protected: YES - login required <br/>
Endpoint:
```
api/v1/protected/job/:id/active 
```

Response:
```
if server fails: HTTP 500
if id is not valid: HTTP 406
if not authenticated: HTTP 401
if job with provided ID does not exist: HTTP 404
if everything ok: HTTP 200
{
   "success": true,
   "active": true/false
}
```

## Pause Job
```
Stops a running job, but does not remove it from the pool and the database
```

Method: POST <br/>
Protected: YES - login required <br/>
Endpoint:
```
api/v1/protected/job/:id/pause
```

Response:
```
if server fails: HTTP 500
if id is not valid: HTTP 406
if not authenticated: HTTP 401
if job with provided ID does not exist: HTTP 404
if job is not running: HTTP 400
if everything ok: HTTP 200
```
## Restart Job
```
Restarts the job with provided id
```
Method: POST <br/>
Protected: YES - login required <br/>
Endpoint:
```
api/v1/protected/job/:id/restart
```

Response:
```
if server fails: HTTP 500
if id is not valid: HTTP 406
if not authenticated: HTTP 401
if job with provided ID does not exist: HTTP 404
if everything ok: HTTP 200
```

## Delete Job
```
Stops the job if running, removes it from the pool and deletes it from the database
```
Method: DELETE <br/>
Protected: YES - login required <br/>
Endpoint:
```
api/v1/protected/job/:id
```

Response:
```
if server fails: HTTP 500
if id is not valid: HTTP 406
if not authenticated: HTTP 401
if job with provided ID does not exist: HTTP 404
if everything ok: HTTP 200
```


