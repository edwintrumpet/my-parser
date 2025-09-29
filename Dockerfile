FROM golang:1.24.7-alpine3.22 AS builder

WORKDIR /usr/src/app
COPY ["go.mod", "./"]
RUN go mod download
COPY . .
RUN go build -o ./bin/app .

FROM alpine

WORKDIR /usr/src/app
COPY --from=builder ["/usr/src/app/bin/app", "/usr/src/app/"]

CMD ["./app"]
