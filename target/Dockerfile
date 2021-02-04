
FROM golang:1.15.0-buster

WORKDIR /opt

COPY ./main.go .
COPY ./kitten.webp .

RUN go build -o target main.go

CMD ["/opt/target"] 
