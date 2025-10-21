FROM golang:1.24-bookworm AS builder

COPY . /go/src/github.com/mintel/k8s-yaml-splitter
WORKDIR /go/src/github.com/mintel/k8s-yaml-splitter
RUN go mod tidy && make build

FROM scratch
COPY --from=builder /go/src/github.com/mintel/k8s-yaml-splitter/bin/k8s-yaml-splitter-* /
ENTRYPOINT ["/k8s-yaml-splitter"]
CMD ["-h"]
