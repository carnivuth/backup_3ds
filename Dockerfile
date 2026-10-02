FROM golang:1.27.1-alpine3.24 AS deps
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM golang:1.27.1-alpine3.24 AS build
WORKDIR /app
COPY --from=deps $GOPATH/pkg/mod $GOPATH/pkg/mod
COPY . .
RUN go build -o console_backupper

FROM alpine:3.24 AS console_backupper

COPY --from=build /app/console_backupper /console_backupper
COPY ./templates /templates
COPY ./static /static
EXPOSE 8080
ENTRYPOINT [ "/console_backupper" ]
