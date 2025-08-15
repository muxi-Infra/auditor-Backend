#!/bin/bash

# 设置默认值
IMAGE_NAME="muxi-auditor-backend"
VERSION=${1:-latest}
PLATFORM=${2:-linux/amd64}

echo "Building Docker image: ${IMAGE_NAME}:${VERSION}"
echo "Platform: ${PLATFORM}"

# 构建镜像
docker build \
    --platform ${PLATFORM} \
    --build-arg VERSION=${VERSION} \
    --build-arg TARGETARCH=$(echo ${PLATFORM} | cut -d'/' -f2) \
    --build-arg TARGETOS=$(echo ${PLATFORM} | cut -d'/' -f1) \
    -t ${IMAGE_NAME}:${VERSION} \
    -t ${IMAGE_NAME}:latest \
    .

echo "Build completed successfully!"
echo "Image: ${IMAGE_NAME}:${VERSION}"
echo "Image: ${IMAGE_NAME}:latest"
