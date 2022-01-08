# syntax=docker/dockerfile:1

#
# Build
#
FROM golang:1.17 AS build
ARG BUILD_ARCH=amd64
ARG BUILD_OS=linux

ADD . /go/src/azure-permissions-checker
WORKDIR /go/src/azure-permissions-checker

RUN BUILD_OS=$BUILD_OS BUILD_ARCH=$BUILD_ARCH make azpermissions

#
# Deploy
#
FROM gcr.io/distroless/static
ARG BUILD_ARCH=amd64

WORKDIR /
COPY --from=build /go/src/azure-permissions-checker/bin/$BUILD_ARCH/azpermissions /azpermissions

ENTRYPOINT [ "/azpermissions" ]
