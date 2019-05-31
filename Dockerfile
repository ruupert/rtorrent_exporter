FROM golang:1.13

ENV GO111MODULE=on

EXPOSE 9135

WORKDIR /go/src/app
COPY . .

RUN go install -v ./...

ENTRYPOINT ["rtorrent_exporter"]
