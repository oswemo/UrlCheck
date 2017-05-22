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
