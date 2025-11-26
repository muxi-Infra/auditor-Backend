#FROM golang:1.24.4 AS builder
#
#COPY . /src
#
#WORKDIR /src
#RUN go mod tidy
#
#RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 go build -o ./bin/auditor-backend .
#
## 为运行镜像准备所需的目录结构
#RUN mkdir -p /src/runtime/app/logs
#
#FROM gcr.io/distroless/base-debian12:nonroot
#
#ENV TZ=Asia/Shanghai
#
#COPY --from=builder --chown=nonroot:nonroot /src/bin/auditor-backend /app/auditor-backend
#COPY --from=builder --chown=nonroot:nonroot /src/runtime/app/logs /app/logs
#
#WORKDIR /app
#
#EXPOSE 8080
#
#VOLUME /data/conf
#
#ENV CONFIG_PATH=/data/conf/config.yaml
#ENV GIN_MODE=release
#
#CMD ["./auditor-backend"]
FROM golang:1.24.4 AS builder

# 安装 librdkafka（运行构建时需要）
RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    pkg-config

COPY . /src
WORKDIR /src
RUN go mod tidy

# 必须启用 CGO
RUN GOPROXY=https://goproxy.cn CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o ./bin/auditor-backend .

RUN mkdir -p /src/runtime/app/logs

# =====================
#   Runtime 镜像
# =====================
FROM gcr.io/distroless/base-debian12:nonroot

ENV TZ=Asia/Shanghai

# 运行时 Deps：必须复制 librdkafka.so
COPY --from=builder --chown=nonroot:nonroot /usr/lib/x86_64-linux-gnu/librdkafka.so* /usr/lib/

COPY --from=builder --chown=nonroot:nonroot /src/bin/auditor-backend /app/auditor-backend
COPY --from=builder --chown=nonroot:nonroot /src/runtime/app/logs /app/logs

WORKDIR /app

EXPOSE 8080

VOLUME /data/conf

ENV CONFIG_PATH=/data/conf/config.yaml
ENV GIN_MODE=release

CMD ["./auditor-backend"]
