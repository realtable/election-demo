FROM golang:1.23.0 AS build

WORKDIR /usr/src

COPY go.mod go.sum .
RUN go mod download

COPY *.go .
COPY handler handler
COPY util util
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM gcr.io/distroless/base-debian12

COPY --from=build /usr/src/app ./app
COPY certs ./certs

CMD ["./app"]
