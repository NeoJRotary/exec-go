# test linux at my poor windows
FROM golang:stretch
WORKDIR /go/src/github.com/NeoJRotary/exec-go
COPY . .
RUN go test