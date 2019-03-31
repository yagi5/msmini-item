FROM golang:1.12

ARG VERSION

WORKDIR /go/src/github.com/yagi5/msmini-item

RUN go get github.com/golang/dep/cmd/dep
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -v -vendor-only

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go install -v \
    -ldflags="-w -s -X main.version=${VERSION}" \
    github.com/yagi5/msmini-item/cmd/item

FROM alpine:latest

RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/item /bin/item
# TODO embed into binary
COPY --from=0 /go/src/github.com/yagi5/msmini-item/items.csv /bin/items.csv

CMD ["/bin/item"]
