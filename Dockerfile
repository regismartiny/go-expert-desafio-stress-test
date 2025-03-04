FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o app

FROM scratch
WORKDIR /app
COPY --from=builder /app .
ENTRYPOINT [ "./app" ]
