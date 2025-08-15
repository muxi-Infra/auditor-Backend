# Docker 部署指南

## 快速开始

### 1. 构建镜像

```bash
# 使用构建脚本
chmod +x docker/build.sh
./docker/build.sh v1.0.0

# 或直接使用 docker build
docker build -t muxi-auditor-backend:latest .
```

### 2. 开发环境运行

```bash
# 启动完整环境（包含 MySQL 和 Redis）
docker-compose up -d

# 仅启动应用（需要外部数据库）
docker-compose -f docker/production.yml up -d
```

### 3. 生产环境部署

```bash
# 使用生产配置
docker-compose -f docker/production.yml up -d

# 设置版本
VERSION=v1.0.0 docker-compose -f docker/production.yml up -d
```

## 配置说明

### 环境变量

- `VERSION`: 应用版本号
- `CONFIG_PATH`: 配置文件路径
- `GIN_MODE`: Gin 运行模式（release/debug）

### 端口映射

- 应用端口: `8080:8080`
- MySQL 端口: `3306:3306` (仅开发环境)
- Redis 端口: `6379:6379` (仅开发环境)

### 数据卷

- 配置文件: `./config:/app/config:ro`
- 日志文件: `./logs:/app/logs`

## 多架构构建

```bash
# AMD64
./docker/build.sh v1.0.0 linux/amd64

# ARM64
./docker/build.sh v1.0.0 linux/arm64

# 多架构构建
docker buildx build --platform linux/amd64,linux/arm64 -t muxi-auditor-backend:latest .
```

## 运行

本项目默认不启用健康检查，如需启用请自行在 `Dockerfile` 或 `docker-compose.yml` 中添加。

## 日志查看

```bash
# 查看应用日志
docker logs muxi-auditor-backend

# 实时日志
docker logs -f muxi-auditor-backend
```

## 故障排除

1. **端口冲突**: 确保 8080 端口未被占用
2. **权限问题**: 确保 logs 目录有写权限
3. **数据库连接**: 检查配置文件中的数据库连接信息
4. **内存不足**: 调整 docker-compose 中的资源限制

## 安全建议

1. 生产环境使用非 root 用户运行
2. 定期更新基础镜像
3. 使用 secrets 管理敏感信息
4. 启用镜像签名验证
