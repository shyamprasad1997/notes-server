FROM golang:1.20

ARG port
COPY . ./src/notes-server
WORKDIR /go/src/notes-server
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notes-server
RUN echo $port
EXPOSE $port

CMD ["./notes-server"]
