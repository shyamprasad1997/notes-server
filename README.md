# notes-server
A server with login signup add notes view notes and delete notes

### Create docker image
```docker build --tag=notes-server --build-arg port=8080 .```

### Run docker image
```docker run -p 8080:8080 notes-server```