#!/bin/bash

set -e

KAFKA_BROKER="muxi-kafka"
KAFKA_PORT=19092

echo " "
echo "==========================================="
echo "🚀 Kafka init script started"
echo "==========================================="
echo " "

# ------------------------------------------
# 等待 Kafka 可用
# ------------------------------------------
echo "⏳ Waiting for Kafka to be reachable at $KAFKA_BROKER:$KAFKA_PORT..."

# ping + nc 检查网络连通性
while ! ping -c 1 $KAFKA_BROKER >/dev/null 2>&1; do
    echo "❌ Cannot ping $KAFKA_BROKER, retrying..."
    sleep 2
done

# 检查端口是否打开
while ! nc -z $KAFKA_BROKER $KAFKA_PORT >/dev/null 2>&1; do
    echo "❌ Kafka port $KAFKA_PORT not open yet, retrying..."
    sleep 2
done

echo "✅ Kafka is reachable!"
echo " "

# ------------------------------------------
# 等待 Kafka 完全 ready
# ------------------------------------------
echo "⏳ Waiting for Kafka to be ready to accept commands..."

while ! /opt/kafka/bin/kafka-topics.sh --bootstrap-server $KAFKA_BROKER:$KAFKA_PORT --list >/dev/null 2>&1; do
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

    /opt/kafka/bin/kafka-topics.sh \
        --create \
        --if-not-exists \
        --bootstrap-server "$KAFKA_BROKER:$KAFKA_PORT" \
        --topic "$name" \
        --partitions "$partitions" \
        --replication-factor "$replicas"

    /opt/kafka/bin/kafka-configs.sh \
        --alter \
        --bootstrap-server "$KAFKA_BROKER:$KAFKA_PORT" \
        --entity-type topics \
        --entity-name "$name" \
        --add-config retention.ms="$retention",cleanup.policy="$policy"

    echo "✅ Topic created or already exists: $name"
    echo " "
done

echo "🎉 All topics created successfully!"
echo "==========================================="
echo " "
