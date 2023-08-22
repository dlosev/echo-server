FROM golang:1.21.0 as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /echo-server

FROM --platform=linux/amd64 ubuntu:20.04
#FROM scratch

ENV PORT 8080

COPY --from=build /echo-server /echo-server

EXPOSE 8080

CMD ["/echo-server"]
