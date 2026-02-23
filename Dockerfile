FROM golang:1.25.7-bookworm AS builder

COPY . /go/src/github.com/ohauer/k8s-yaml-splitter
WORKDIR /go/src/github.com/ohauer/k8s-yaml-splitter
RUN go mod tidy && make build

FROM scratch
COPY --from=builder /go/src/github.com/ohauer/k8s-yaml-splitter/bin/k8s-yaml-splitter-linux-amd64 /k8s-yaml-splitter
ENTRYPOINT ["/k8s-yaml-splitter"]
CMD ["-h"]
