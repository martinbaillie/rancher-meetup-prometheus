FROM golang:alpine as builder
MAINTAINER Martin Baillie <martin.t.baillie@gmail.com>
ARG HOST
ARG PROJECT
ARG LDFLAGS

WORKDIR /go/src/${HOST}/martinbaillie/${PROJECT}
COPY ${PWD} .
RUN apk --no-cache add git make && \
        make dep-update && \
        go get -u github.com/jteeuwen/go-bindata/... && \
        go generate && CGO_ENABLED=0 GOARCH=amd64 GOOS=linux \
        go build -a -installsuffix cgo -o app -ldflags "${LDFLAGS}"

FROM alpine
ARG HOST
ARG PROJECT

WORKDIR /
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/${HOST}/martinbaillie/${PROJECT}/app .
ENTRYPOINT ["/app"]
