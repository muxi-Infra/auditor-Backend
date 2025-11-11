FROM golang:1.24.4 AS builder

COPY . /src

WORKDIR /src
RUN go mod tidy

RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 go build -o ./bin/auditor-backend .

# 为运行镜像准备所需的目录结构
RUN mkdir -p /src/runtime/app/logs

FROM gcr.io/distroless/base-debian12:nonroot

ENV TZ=Asia/Shanghai

COPY --from=builder --chown=nonroot:nonroot /src/bin/auditor-backend /app/auditor-backend
COPY --from=builder --chown=nonroot:nonroot /src/runtime/app/logs /app/logs

WORKDIR /app

EXPOSE 8080

VOLUME /data/conf

ENV CONFIG_PATH=/data/conf/config.yaml
ENV GIN_MODE=release

CMD ["./auditor-backend"]
