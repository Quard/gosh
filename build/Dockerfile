FROM golang:1.12

WORKDIR /app

COPY . /app

RUN go build

EXPOSE 5000

CMD ["./gosh", "-storage", "redis"]