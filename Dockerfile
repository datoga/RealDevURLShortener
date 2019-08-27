FROM golang AS builder

ENV GO111MODULE=on

RUN mkdir -p /go/src/URLShortener
WORKDIR /go/src/URLShortener

ADD . /go/src/URLShortener

WORKDIR /go/src/URLShortener

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s -extldflags "-static"' .

FROM alpine

RUN mkdir /app 
WORKDIR /app
COPY --from=builder /go/src/URLShortener/RealDevURLShortener .

CMD ["/app/RealDevURLShortener"]
