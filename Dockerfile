FROM alpine:3.3
MAINTAINER Dustin Blackman

ENV GOROOT /usr/lib/go
ENV GOPATH /gopath
ENV GOBIN /gopath/bin
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

COPY . /gopath/src/github.com/dustinblackman/speakerbot

RUN apk add --update ffmpeg opus opus-dev bash git make pkgconfig build-base && \
  apk add go --update-cache --repository http://dl-cdn.alpinelinux.org/alpine/edge/community/ --allow-untrusted && \
  cd /gopath/src/github.com/dustinblackman/speakerbot && \
  make && \
  make build && \
  mkdir /app && \
  mv ./speakerbot /app/ && \
  apk del go git make pkgconfig opus-dev build-base && \
  rm -rf /usr/share/man /tmp/* /var/tmp/* /var/cache/apk/* /gopath

WORKDIR /app
CMD ["./speakerbot"]
