FROM  golang:1.21.0 as builder
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace

COPY ./ ./

RUN --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go \
    go build -o gecko-server ./cmd/server

FROM sanmuyan/ubuntu:22.04-base-20250427

RUN apt install git -y

COPY --from=builder /workspace/gecko-server ./

CMD ["/opt/app/gecko-server"]