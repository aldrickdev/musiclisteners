FROM golang:1.22.4 AS builder
WORKDIR /src
COPY src /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /src/build/app cmd/app/main.go

FROM scratch AS runner
WORKDIR /app
COPY --from=builder /src/build/app .
ENTRYPOINT ["./app"]
