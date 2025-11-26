#!/bin/bash

set -e

KAFKA_BROKER="kafka:9092" # 如果容器名不是kafka,请改为对应的名字。

echo " "
echo "==========================================="
echo "🚀 Kafka init script started"
echo "==========================================="
echo " "

# ------------------------------------------
# 等待 Kafka 可用
# ------------------------------------------
echo "⏳ Waiting for Kafka to be ready at $KAFKA_BROKER..."

while ! kafka-topics.sh --bootstrap-server $KAFKA_BROKER --list >/dev/null 2>&1; do
    echo "🔄 Kafka is not ready yet. Retrying..."
    sleep 2
done

echo "✅ Kafka is ready!"
echo " "

# ------------------------------------------
# 要创建的 Topic 列表
# 格式：topic_name partitions replication retention_ms cleanup_policy
# ------------------------------------------
TOPICS=(
  "event_stream 2 1 86400000 delete" # 在此处添加你要创建的主题
)

# ------------------------------------------
# 创建 Topic
# ------------------------------------------

echo "📌 Starting to create Kafka topics..."
echo " "

for topic in "${TOPICS[@]}"; do
    read -r name partitions replicas retention policy <<< "$topic"

    echo "➡ Creating topic: $name"
    echo "   Partitions: $partitions, Replicas: $replicas"
    echo "   retention.ms: $retention, cleanup.policy: $policy"

    kafka-topics.sh \
        --create \
        --if-not-exists \
        --bootstrap-server "$KAFKA_BROKER" \
        --topic "$name" \
        --partitions "$partitions" \
        --replication-factor "$replicas"

    # 设置高级参数
    kafka-configs.sh \
        --alter \
        --bootstrap-server "$KAFKA_BROKER" \
        --entity-type topics \
        --entity-name "$name" \
        --add-config retention.ms="$retention",cleanup.policy="$policy"

    echo "✅ Topic created or already exists: $name"
    echo " "
done

echo "🎉 All topics created successfully!"
echo "==========================================="
echo " "
