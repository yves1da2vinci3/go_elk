FROM golang:1.17-alpine

COPY go.mod go.sum ./
WORKDIR /github.com/yves1da2vinci3/go_elk
RUN rm $GOPATH/go.mod
RUN go mod download
COPY . /github.com/yves1da2vinci3/go_elk
RUN go build -o /usr/bin/go_elk github.com/yves1da2vinci3/go_elastic/cmd/api

EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/go_elk"]