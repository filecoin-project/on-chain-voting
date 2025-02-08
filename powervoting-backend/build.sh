#!/bin/bash
set -e

IMAGE_NAME="power-voting-backend"

if [ ! -f "configuration.yaml" ]; then
    echo "Error: configuration.yaml does not exist."
    exit 1
fi

PORT=$(awk '/^server:/{flag=1;next} /^  port:/{if(flag) print $2; flag=0}' "configuration.yaml" | tr -d ':')

if [ -z "$PORT" ]; then
  echo "Error: Can not get port from configuration."
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

docker run --name $IMAGE_NAME -v ./configuration.yaml:/dist/configuration.yaml -v ./proof.ucan:/dist/proof.ucan -p $PORT:$PORT -d $IMAGE_NAME