FROM golang:1.6.2

WORKDIR /
COPY ../edurtc-biz /edurtc-biz

RUN go get
RUN go build edurtc-biz.go

EXPOSE 8080
ENTRYPOINT ["/edurtc-biz/edurtc-biz"]