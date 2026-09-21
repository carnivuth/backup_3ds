FROM golang:1.27.1-alpine3.24 AS build
WORKDIR /app
COPY . .
RUN go build -o console_backupper

FROM alpine:3.24 AS console_backupper

COPY --from=build /app/console_backupper /console_backupper
COPY ./templates /templates
EXPOSE 8080
ENTRYPOINT [ "/console_backupper" ]
