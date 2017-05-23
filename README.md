## Url Check

### Objective:
Given a host and request path (optional query string), return a result that notifies the user whether or not the host and request path is safe to load.   

### Storage:
Initially, storage will be handed with a MongoDB backend, fronted by Memcached with a long TTL.  This allows for multiple copies of the service to run in parallel, while limiting stress on the backend as cache handles initial load (assuming objects in cache).   

Storage is seeded initially through data downloaded from http://www.phishtank.com/, and processed via a quick python script that uses `urlparse` to break the url into it's various pieces.

### Deployment:
Service will be deployed within a docker container with configuration for deploying dependant services as well.  Scaling of resources will be dependant on further configuration, taking into account CPU and memory reservations with an auto-scaling group.   

### Design Decisions:

#### Programming Language:
I opted to use GO as the programming language of choice for this project.  While I'm more comfortable hacking out code in python, GO feels like a better choice for deploying small micro-services within a container environment.  GO also benefits from simple API handling built-in without making use of third-party frameworks such as django.   

#### Database:
I opted to not bother starting with in-memory storage.  While possible, it would require having a method of ensuring that the service had access to a list of URLs that could be loaded into memory upon startup, along with a method of synchronizing updates across instances as the service scales.  Instead, I opted to go with a simple database that could store the data required along with a caching mechanism to reduce load.  

### Building:
Execute `make build` to build the project.

### Deploying:
Execute `make run` to run the project.

### Known Issues

1.  Encoding errors in the request cause Gorilla MUX errors that do not conform to the standard responses from other parts of the API.
2.  For consistency, each start of the service (`make run`) will wipe data and start fresh.
3.  Database connection and error handling code is rudimentary.  Not production ready.

### Testing   

Requests can be performed using any HTTP client, though most testing to date has been performed via `curl`   
```shell
curl -s -XGET http://localhost:8001/urlinfo/1/wirtualnyanalityk.pl:80/%2Fadministrator%2Fcomponents%2Fcom_content%2Felements%2Findex.htm
```

```shell
curl -s -XPUT http://localhost:8001/urlinfo/1/wirtualnyanalityk.pl:80/%2Fadministrator%2Fcomponents%2Fcom_content%2Felements%2Findex.htm
```

Note that the request path and query string are expected to be encoded, where the hostname and port are not.   

Responses are all in JSON (barring those that do not make it past the router, which is a known issue.)   
All responses will contain a `status` object and a `data` object.  The `status` object is meant to reflect the HTTP status code and message.  The contents of the `data` object are dependant on the API endpoint being queried.

   ```javascript
   { "status": { "code": 200, "message": "OK" }, "data": { "safe": true } }
   ```

Any errors that occur, will return an error object for data.   

   ```javascript
   { "status": {"code": 400, "message": "Bad Request" }, "data": { "error": "Hostname does not appear to be a valid format" } }
   ```

If the response has no data beyond the HTTP status code, then it will be an empty object.

   ```javascript
   { "status": {"code": 500, "message": "Internal Server Error" }, "data": { } }
   ```
