# 仓库根目录
REPODIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
# 构建输出目录
BUILDDIR := $(REPODIR)/dist

# 服务名称（与docker-compose.yml保持一致）
ES_CONTAINER := elasticsearch
KIBANA_CONTAINER := kibana
LOGSTASH_CONTAINER := logstash
FILEBEAT_CONTAINER := filebeat

#服务用户
LOGSTASH_USER:=logstash_writer
ELASTIC_USER=elastic
ROLE_NAME := logstash_writer

#net
ELASTICSEARCH_PORT=9200
ELASTICSEARCH_HOST := localhost

# 统一密码
ES_PASSWORD := changeme123

# compose文件路径（关键：指向es子目录）
ES_COMPOSE_FILE := ./filebeat/docker-compose.yaml

.PHONY: tests
tests:
	@echo "beginning to run test...."


.PHONY: build
build:
	@echo "Cleaning up and downloading modules..."
	go mod tidy
	@echo "Building for Linux amd64..."
	mkdir -p $(BUILDDIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILDDIR)/app $(REPODIR)
	@echo "Build completed: $(BUILDDIR)/app"

#deploy everything
.PHONY:deploy
deploy:
	docker compose up -d

# only app
.PHONY:docker
docker:
	docker-compose -f docker/production.yml up -d

.PHONY:es
es:
	@echo "可用子目标：start/stop"

.PHONY:es.start
es.start:
	docker-compose -f $(ES_COMPOSE_FILE) up -d
	@echo "ES服务已启动（Elasticsearch 初始化约30秒）"

.PHONY:es.stop
es.stop:
	docker-compose -f $(ES_COMPOSE_FILE) down
	@echo "ES服务已停止"
.PHONY:es.createUser
es.createUser:
	@echo "进入 Elasticsearch 容器，创建 $(LOGSTASH_USER) 用户..."
	# 通过 docker exec 进入容器内部执行命令
	docker exec -it $(ES_CONTAINER) bash -c '\
		curl -X POST "http://$(ELASTICSEARCH_HOST):$(ELASTICSEARCH_PORT)/_security/user/$(LOGSTASH_USER)" \
		  -u $(ELASTIC_USER):$(ES_PASSWORD) \
		  -H "Content-Type: application/json" \
		  -d "{ \
		    \"password\": \"$(ES_PASSWORD)\", \
		    \"roles\": [\"$(ROLE_NAME)\"], \
		    \"full_name\": \"Logstash Writer User\" \
		  }" \
	'
	@echo "$(LOGSTASH_USER) 用户创建完成（若提示 'created':true 则成功）"

.PHONY:es_password
es_password:
	@echo "可用子目标:reset/verify"

.PHONY: es_password.reset
es_password.reset:
	@echo "正在重置 kibana_system 密码..."
	# 先输入y确认操作，再输入两次密码（解决确认步骤问题）
	docker exec -it $(ES_CONTAINER) bash -c 'echo -e "y\n$(ES_PASSWORD)\n$(ES_PASSWORD)" | ./bin/elasticsearch-reset-password -u kibana_system -i'
	@echo "正在重置 logstash_system 密码..."
	docker exec -it $(ES_CONTAINER) bash -c 'echo -e "y\n$(ES_PASSWORD)\n$(ES_PASSWORD)" | ./bin/elasticsearch-reset-password -u $(LOGSTASH_USER) -i'
	@echo "正在重置 elastic 密码..."
	docker exec -it $(ES_CONTAINER) bash -c 'echo -e "y\n$(ES_PASSWORD)\n$(ES_PASSWORD)" | ./bin/elasticsearch-reset-password -u elastic -i'
	@echo "所有用户密码已重置为：$(ES_PASSWORD)"


.PHONY: es_password.verify
es_password.verify:
	@echo "验证 kibana_system 密码..."
	docker exec -it $(ES_CONTAINER) curl -u kibana_system:$(ES_PASSWORD) http://localhost:9200/_cluster/health || echo "kibana_system 密码错误"
	@echo "\n验证 logstash_system 密码..."
	docker exec -it $(ES_CONTAINER) curl -u $(LOGSTASH_USER):$(ES_PASSWORD) http://localhost:9200/_cluster/health || echo "$(LOGSTASH_USER) 密码错误"
	@echo "\n验证 elastic 密码..."
	docker exec -it $(ES_CONTAINER) curl -u elastic:$(ES_PASSWORD) http://localhost:9200/_cluster/health || echo "elastic 密码错误"


run:
	cd $(REPODIR) && ./dist/app
