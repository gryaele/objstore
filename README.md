# In Memory Object Store
HTTP service that is able to store objects organized by buckets

## API description:
### Upload object to the service
- Request:  PUT /objects/{bucket}/{objectID}
- Response:  Status: 201 Created {  "id": "<objectID>" }

### Download an object from the service
- Request: GET /objects/{bucket}/{objectID}
- Response if object is found:  Status: 200 OK {object data}
- Response if object is not found:  Status 404 Not Found

### Delete an object from the service
- Request:  DELETE /objects/{bucket}/{objectID}
- Response if object found: Status: 200 OK
- Response if object not found:  Status: 404 Not Found


## Run tests
```go test ./...```

## How to run

```go run .```

### Create an object

```curl -v -X PUT  -d 'some file content' localhost:8080/objects/bucket1/object1```

### Get object by id
```curl -v -X GET localhost:8080/objects/bucket1/object1```

### Delete an object
```curl -v -X DELETE localhost:8080/objects/bucket1/object1```