FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
COPY .env ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /main
COPY --from=builder /app/.env ./  
EXPOSE 8080

ENTRYPOINT ["/main"]

FROM builder as tester
WORKDIR /app
CMD ["go", "test", "./limiter", "-v"]