# Build Stage
FROM golang:1.21.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cyberpecker-api


# Final Stage
FROM scratch

COPY --from=builder /app/cyberpecker-api /

EXPOSE 8000

CMD [ "/cyberpecker-api" ]