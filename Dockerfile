##
## Build
##

FROM golang:1.16-buster AS build

WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -o /app

EXPOSE 8080

#ENV HTTP_PORT=8090

ENTRYPOINT [ "./app" ]
