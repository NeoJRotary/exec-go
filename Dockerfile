# test linux at my poor windows
FROM golang
WORKDIR /go/src/app
COPY . .
CMD go test