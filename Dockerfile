FROM golang:1.17-alpine
WORKDIR /thrive-webserver
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]