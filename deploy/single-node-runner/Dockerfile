FROM golang:1.15-alpine AS build-env

ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

WORKDIR /go/src/github.com/onomyprotocol/market

COPY . .

RUN apk add --no-cache $PACKAGES && \
  make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
 COPY --from=build-env /go/bin/marketd /usr/bin/marketd

COPY deploy/single-node-runner/scripts /root/scripts
COPY deploy/single-node-runner/market_home /root/.market

CMD sh /root/scripts/run.sh