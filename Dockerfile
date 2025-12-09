###############################################
# Builder Stage
###############################################
FROM golang:1.24.4 AS builder

ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# 清空其他源，只使用阿里源
RUN rm -f /etc/apt/sources.list.d/* \
 && echo 'deb http://mirrors.aliyun.com/debian bookworm main contrib non-free' > /etc/apt/sources.list \
 && echo 'deb http://mirrors.aliyun.com/debian bookworm-updates main contrib non-free' >> /etc/apt/sources.list \
 && echo 'deb http://mirrors.aliyun.com/debian-security bookworm-security main contrib non-free' >> /etc/apt/sources.list \
 && echo 'Acquire::Check-Valid-Until "false";' > /etc/apt/apt.conf.d/10no-check-valid-until

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    build-essential \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY . .
RUN go mod tidy
RUN go build -o auditor-backend .

###############################################
# Runtime Stage
###############################################
FROM debian:stable-slim AS runtime

# 清空其他源，只使用阿里源
RUN rm -f /etc/apt/sources.list.d/* \
 && echo 'deb http://mirrors.aliyun.com/debian bookworm main contrib non-free' > /etc/apt/sources.list \
 && echo 'deb http://mirrors.aliyun.com/debian bookworm-updates main contrib non-free' >> /etc/apt/sources.list \
 && echo 'deb http://mirrors.aliyun.com/debian-security bookworm-security main contrib non-free' >> /etc/apt/sources.list \
 && echo 'Acquire::Check-Valid-Until "false";' > /etc/apt/apt.conf.d/10no-check-valid-until

RUN apt-get update && apt-get install -y \
    librdkafka1 \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -m appuser
USER appuser

WORKDIR /app
COPY --from=builder /src/auditor-backend /app/auditor-backend

ENV TZ=Asia/Shanghai
ENV CONFIG_PATH=/data/conf/config.yaml
ENV GIN_MODE=release

VOLUME /data/conf

EXPOSE 8080
CMD ["./auditor-backend"]
