FROM golang:1.13.4

WORKDIR /go/src/github.com/communitybridge/easycla-api
COPY . .
RUN apt-get update && apt-get install sudo -y
RUN make setup
RUN make clean swagger swagger-validate deps fmt build test lint

ENTRYPOINT ["./bin/cla-api"]
