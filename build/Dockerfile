FROM onosproject/golang-build:v1.0 as build

ENV GO111MODULE=on
ENV XAPPNAME=ransim-api-demo

COPY Makefile go.mod go.sum /go/src/github.com/rimedo-labs/ransim-api-demo/
COPY cmd/ /go/src/github.com/rimedo-labs/ransim-api-demo/cmd/
COPY vendor/ /go/src/github.com/rimedo-labs/ransim-api-demo/vendor/
COPY pkg/ /go/src/github.com/rimedo-labs/ransim-api-demo/pkg/
COPY build/build-tools /go/src/github.com/rimedo-labs/ransim-api-demo/build/build-tools/

RUN cd /go/src/github.com/rimedo-labs/ransim-api-demo && GOFLAGS=-mod=vendor make build

FROM alpine:3.11
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/rimedo-labs/ransim-api-demo/build/_output/ransim-api-demo /usr/local/bin/ransim-api-demo

ENTRYPOINT ["ransim-api-demo"]