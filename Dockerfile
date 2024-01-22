FROM docker.io/golang:1.20 as builder

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /sanic

# -----

FROM builder AS tester

RUN go test -v ./...

# -----

FROM scratch as runner

WORKDIR /

COPY --from=builder /sanic /sanic
COPY --from=builder /app/static /static

EXPOSE 8080

ENTRYPOINT ["/sanic"]
