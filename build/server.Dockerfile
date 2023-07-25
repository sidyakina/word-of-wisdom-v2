FROM golang:1.20.6

WORKDIR server

COPY cmd/server ./cmd/server
COPY internal/general ./internal/general
COPY internal/server ./internal/server
COPY pkg ./pkg
COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go build -o /word-of-wisdom-server ./cmd/server

CMD ["/word-of-wisdom-server"]