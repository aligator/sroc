FROM golang:alpine AS builder

RUN apk --no-cache add ca-certificates git

WORKDIR /app
COPY  . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app", "-listen=0.0.0.0:4242"]
