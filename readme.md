
# Golang Hash API
A small application that can take a plain text password, hash it and keep some basic metrics about the process

## Starting the server
In the root directory run:
```sh 
go run cmd/main.go
```

## Using the API
The applications offers 4 endpoints:
### Hash Password
Takes a plain text password, hashes it, base64 encodes it and stores for later use
```
curl --data "password=angryMonkey" http://localhost:8080/hash
```
Example Response:
```json
1
```
### Get Password Hash
Fetch the password after it has been stored
```
curl http://localhost:8080/hash/1
```
Example Response:
```
"7iaw3Ur350mqGo7jwQrpkj9hiYB3Lkc/iBml1JQODbJ6wYX4oOHV+E+IvIh/1nsUNzLDBMxfqa2Ob1f1ACio/w=="
```
## Get Statistics
Get basic statistics about the `Password Hash` endpoint
```
curl http://localhost:8080/stats
```
Example Response:
```json
{
    "total":2,
    "average":0.0914775
}
```

## Shutdown Application
Gracefully shutsdown the application and lets any inflight processes finish
```
curl http://localhost:8080/shutdown
```
Example Response:
```json
OK
```
## Testing 
Included in the `tests/` directory is a `./load_test.sh` script for running many hash requests 
against the `Hash Password` endpoint. 

Along with the load test, a postman Library is also available in the `tests/postman` directory. 

## Project Notes
The Go language is new to me.  I used documentation and many tutorials to figure out how 
the mux server, middleware, and session storage work

[Effective Go](https://golang.org/doc/effective_go)

[Middleware Patterns](https://drstearns.github.io/tutorials/gomiddleware)

[Building Middleware](https://www.alexedwards.net/blog/making-and-using-middleware)

[Sync Maps](https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c)

[Golang Gotchas](https://yourbasic.org/golang/#go-gotchas)

[Golang File Structure](https://github.com/golang-standards/project-layout)

[Golang Standard Library](https://pkg.go.dev/std#stdlib)