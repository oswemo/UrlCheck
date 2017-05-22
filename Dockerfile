FROM golang:1.8.1

ADD . build
RUN cd /go/build && make

WORKDIR /go/build

CMD ./bin/urlcheck
