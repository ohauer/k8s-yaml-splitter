FROM golang:1.27.4-bookworm@sha256:4fa7c968ba93f1635670522dbb01a2c78150443e1385746924566079f36570ea AS builder

COPY . /go/src/github.com/ohauer/k8s-yaml-splitter
WORKDIR /go/src/github.com/ohauer/k8s-yaml-splitter
RUN go mod tidy && make build

FROM scratch
COPY --from=builder /go/src/github.com/ohauer/k8s-yaml-splitter/bin/k8s-yaml-splitter-linux-amd64 /k8s-yaml-splitter
USER 65534:65534
ENTRYPOINT ["/k8s-yaml-splitter"]
CMD ["-h"]
