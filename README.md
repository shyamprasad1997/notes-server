# notes-server
A server with login signup add notes view notes and delete notes

### Create docker image
Run ```docker build --tag=notes-server --build-arg port=8080 .```

### Run docker image
Run ```docker run -p 8080:8080 notes-server```

## Run tests
First run ```go get github.com/vektra/mockery/v2```
Then run ```go test ./...```

## Run server without docker
Run ```go build && ./notes-server```