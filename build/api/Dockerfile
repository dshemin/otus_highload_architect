FROM golang:1.20-alpine as builder

ADD . /code

WORKDIR /code

RUN go mod download

RUN go build -o api cmd/api/main.go

FROM scratch

COPY --from=builder /code/api /api

ENTRYPOINT ["/api"]
