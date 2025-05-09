FROM golang:1.24.3-bookworm

COPY . /go/src/github.com/mintel/k8s-yaml-splitter
WORKDIR /go/src/github.com/mintel/k8s-yaml-splitter
RUN make vendor && make

FROM scratch
COPY --from=0 /go/src/github.com/mintel/k8s-yaml-splitter/bin/k8s-yaml-splitter* /
ENTRYPOINT ["/k8s-yaml-splitter"]
CMD ["--help"]

