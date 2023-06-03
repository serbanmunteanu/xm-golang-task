FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./executable main.go


FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/executable .
COPY --from=builder /app/swagger/swagger.json .
COPY --from=builder /app/config.yml .

EXPOSE 8080
CMD ["./executable"]