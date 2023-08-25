FROM  golang:1.21.0 as builder
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace

COPY ./ ./

RUN go build -o gecko-server ./cmd/server

FROM ubuntu:22.04
ENV LANG=en_US.utf8

RUN set -x \
    && mkdir -p /data \
    && yum install git -y

WORKDIR /data

COPY --from=builder /workspace/gecko-server ./
COPY --from=builder /workspace/config/ ./config

CMD ["/data/gecko-server"]