FROM golang:1.17-alpine as base

RUN apk update && apk add make git

WORKDIR /go/src/todo-app-server

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN make go-build

FROM alpine

RUN apk --no-cache add iptables ca-certificates && \
    update-ca-certificates curl

RUN mkdir /app

COPY --from=base /go/src/todo-app-server/todo-app-server /app

WORKDIR /app

EXPOSE 8080 8081
ENTRYPOINT ["/app/todo-app-server"]

