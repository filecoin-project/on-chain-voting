#!/bin/bash
set -e

IMAGE_NAME="pv-backend-mainnet"

if [ ! -f ".env" ]; then
    echo "Error: .env does not exist."
    exit 1
fi

PORT=$(awk -F'=' '/^PORT=/ {gsub(/[^0-9]/, "", $2); print $2}' .env)

if [ -z "$PORT" ]; then
  echo "Error: Can not get port from .env."
  exit 1
fi

echo "project: $IMAGE_NAME, port: $PORT"

(
    if ! git show-ref --verify --quiet "refs/heads/main"; then
        echo "Error: Branch main does not exist."
        exit 1
    fi

    git checkout main
    git pull origin main
)

docker build -t $IMAGE_NAME .

if docker ps -a --format '{{.Names}}' | grep -wq "$IMAGE_NAME"; then
    echo "Stopping container: $IMAGE_NAME..."
    docker stop $IMAGE_NAME
    docker rm $IMAGE_NAME

else
    echo "Container $IMAGE_NAME does not exist or is already stopped."
fi

docker run --name $IMAGE_NAME  -p $PORT:$PORT -d $IMAGE_NAME