FROM golang:1.11.5-stretch

WORKDIR /go/src/test/peterjasc-golang-test

COPY . .

ENV PORT 8080
EXPOSE ${PORT}

RUN go get -u github.com/golang/dep/cmd/dep \
&& dep ensure \
&& go test -v ./... \
&& mkdir -p /go/bin \
&& go build -o /go/bin/app main.go \
&& cp /go/bin/app /app

CMD ["/app"]