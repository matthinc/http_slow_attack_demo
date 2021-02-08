
FROM golang:1.15.0-buster

WORKDIR /opt

COPY ./main.go .
COPY ./www ./www

RUN go build -o target main.go

CMD ["/opt/target"] 
