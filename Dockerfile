# base image build stage
FROM golang:1.17.3-alpine3.14 as builder-base
RUN apk add --no-cache \
    bash \
    git \
    build-base

# build stage
FROM builder-base AS builder

WORKDIR /build
COPY . /build

#instead, we use go mod vendor and commit the vendor folder to git repo.
#RUN go mod download
RUN make build

# final stage
FROM alpine:3.14
RUN apk add --no-cache bash ca-certificates tzdata
COPY --from=builder /build/bin/* /opt/the-automatic-manager/

WORKDIR /opt/the-automatic-manager/
