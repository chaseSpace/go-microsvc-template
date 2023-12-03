#!/usr/bin/env bash
set -e

SVC=$1
MICRO_SVC_ENV=$2
TAG=$3

# 配置区
IMAGE_REPO=leigg # 替换为你的镜像仓库
# -----

if [ -z "$SVC" ] || [ -z "$MICRO_SVC_ENV" ]; then
  echo "Wrong Input!"
  echo "Example: sh build_image.sh SVC ENV [tag]"
  exit 1
  fi

if [ ! -z "$TAG" ]; then
  TAG=":$TAG"
  fi

docker build --build-arg SVC=$SVC --build-arg MICRO_SVC_ENV=$MICRO_SVC_ENV . -t $IMAGE_REPO/go-svc-$SVC$TAG
docker push $IMAGE_REPO/go-svc-$SVC$TAG