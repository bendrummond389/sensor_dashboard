#!/bin/bash

docker buildx create --use

# Build and push MQTT broker
docker buildx build --platform linux/arm64/v8 -t bendrummond389/mqtt-broker:3.14 -f ./mqtt/Dockerfile.mqtt ./mqtt --push

echo "MQTT Broker Docker image built and pushed successfully!"

# Build and push MQTT server
docker buildx build --platform linux/arm64/v8 -t bendrummond389/go-mqtt-server:3.14 -f ./server/Dockerfile.server ./server --push

echo "Go MQTT Server Docker image built and pushed successfully!"

