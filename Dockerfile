
FROM golang:1.18

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -o main .

EXPOSE 8080

# Pass environment variables to the container
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=Nuhman@456
ENV DB_NAME=postgres

# Command to run the executable
CMD ["./main"]
