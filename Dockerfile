FROM golang:1.20.6 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /service


FROM alpine:latest
WORKDIR /app
COPY --from=build /service /app/service
ENV JWT_SECRET="unsecure"
ENV PEPPER="thisispepper"
ENV DB_URL="host=localhost port=5432 user=miauw password=password sslmode=disabled timezone=Europe/Berlin"
ENV REDIS_HOST="localhost:6379"
ENV REDIS_PASS=""
ENV RABBITMQ="amqp://guest:guest@localhost:5672"
CMD ["/app/service"]


