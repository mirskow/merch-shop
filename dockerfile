FROM golang:1.23

WORKDIR /avito-shop
COPY . /avito-shop

RUN go build -o /build ./cmd \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]