FROM docker.io/golang:1.22 as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o sanic ./...

# -----

FROM scratch as runner

WORKDIR /

COPY --from=builder /usr/src/app/sanic /sanic
COPY --from=builder /usr/src/app/static /static

EXPOSE 8080

ENTRYPOINT ["/sanic"]
