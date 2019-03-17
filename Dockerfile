FROM golang:1.12

ARG VERSION

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

CMD ["/bin/server"]
