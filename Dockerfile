FROM docker.io/golang:alpine as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o sanic ./...

# -----

FROM scratch as runner

WORKDIR /

COPY --from=builder /usr/src/app/sanic /sanic

EXPOSE 8080
EXPOSE 8443

ENTRYPOINT ["/sanic"]
