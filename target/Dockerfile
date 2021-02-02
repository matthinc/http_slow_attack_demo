
FROM golang:1.15.0-buster

WORKDIR /opt

COPY ./main.go .

RUN go build -o target main.go

CMD ["/opt/target"] 
