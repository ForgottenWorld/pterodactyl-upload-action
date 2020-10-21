FROM golang:alpine

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /go/app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
RUN go mod verify

# Copy the source from the current directory to the Working Directory inside the container
COPY entrypoint.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/entrypoint .

# Command to run the executable
ENTRYPOINT ["/go/bin/entrypoint"]
