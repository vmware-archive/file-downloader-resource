FROM golang:alpine as builder
COPY . /go/src/github.com/calebwashburn/file-downloader
RUN apk update && apk add bash git unzip curl
ENV CGO_ENABLED 0
RUN go get -a -t github.com/onsi/gomega
RUN go build -o /assets/in github.com/calebwashburn/file-downloader/in
RUN go build -o /assets/out github.com/calebwashburn/file-downloader/out
RUN go build -o /assets/check github.com/calebwashburn/file-downloader/check
WORKDIR /go/src/github.com/calebwashburn/file-downloader
RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM alpine:edge AS resource
RUN apk add --no-cache bash tzdata ca-certificates git jq openssh
RUN git config --global user.email "git@localhost"
RUN git config --global user.name "git"
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*
