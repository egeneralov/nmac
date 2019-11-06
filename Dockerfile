FROM golang:1.13-alpine

RUN apk add --no-cache ca-certificates
RUN adduser -D -g '' appuser

ENV \
  GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /go/src/github.com/egeneralov/nmac
ADD go.mod go.sum /go/src/github.com/egeneralov/nmac/
RUN go mod download

ADD . .

RUN \
  go build -v -installsuffix cgo -ldflags="-w -s" -o /go/bin/nmac-parse-api cmds/nmac-parse-api/main.go && \
  go build -v -installsuffix cgo -ldflags="-w -s" -o /go/bin/nmac-parse-page cmds/nmac-parse-page/main.go && \
  go build -v -installsuffix cgo -ldflags="-w -s" -o /go/bin/nmac-parse-all cmds/nmac-parse-all/main.go


FROM alpine:3.10

RUN apk add --no-cache ca-certificates
# RUN adduser -D -g '' appuser

# COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=0 /etc/passwd /etc/passwd
COPY --from=0 /go/bin /go/bin

USER nobody
# ENTRYPOINT ["/go/bin/nmac-parse-api"]
EXPOSE 8018
ENV PATH='/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin'
# CMD /go/bin/nmac-parse-api
