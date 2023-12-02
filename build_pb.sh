#!/usr/bin/env bash
set -e

# 此脚本支持在linux、mac环境下执行
# - mac环境请在根目录下新建 tool_mac目录（不要提交到git），并自行下载tool/目录下的二进制文件，脚本会自动调用

_pb_dir=$1 # 指定生成的pb子目录，比如user 对应 microsvc/svc/user
_protoc_path=""

if [[ $(uname) == "Linux" ]]; then
    echo "Running on Linux"
    _protoc_path='./tool/protoc_v24'
elif [[ $(uname) == "Darwin" ]]; then
    echo "Running on macOS"
    _protoc_path='./tool_mac/protoc_v24'
elif [[ $(uname) == *MINGW* ]]; then
    echo "Running on Windows"
    _protoc_path='./tool_win/protoc_v24'
else
    echo "Unknown OS"
    exit 1
fi

chmod +x $_protoc_path/*

PATH=$PATH:$_protoc_path

OUTPUT_DIR="./protocol/"
mkdir -p $OUTPUT_DIR

case $_pb_dir in
clear)
  rm -rf ./protocol/*
  ;;
*)
  if [[ -n $_pb_dir ]]; then
    echo "regenerate proto/svc/$_pb_dir/..."
    rm -rf ./protocol/$_pb_dir/*

    $_protoc_path/protoc -I ./proto/ -I ./proto/include/ \
      --go_out=$OUTPUT_DIR \
      --go_opt paths=source_relative \
      --go-grpc_out=$OUTPUT_DIR \
      --go-grpc_opt=paths=source_relative \
      --go-grpc_opt=require_unimplemented_servers=false \
      proto/svc/$_pb_dir/*.proto proto/svc/*.proto
#      --grpc-gateway_out=$OUTPUT_DIR \
#      --grpc-gateway_opt logtostderr=true \
#      --grpc-gateway_opt paths=source_relative \
#      --grpc-gateway_opt generate_unbound_methods=true \

  else
    echo "regenerate all proto files..."
    rm -rf ./protocol/*

    $_protoc_path/protoc -I ./proto/ -I ./proto/include/ \
      --go_out=$OUTPUT_DIR \
      --go_opt paths=source_relative \
      --go-grpc_out=$OUTPUT_DIR \
      --go-grpc_opt=paths=source_relative \
      --go-grpc_opt=require_unimplemented_servers=false \
      proto/svc/*/*.proto proto/svc/*.proto
#      --grpc-gateway_out=$OUTPUT_DIR \
#      --grpc-gateway_opt logtostderr=true \
#      --grpc-gateway_opt paths=source_relative \
#      --grpc-gateway_opt generate_unbound_methods=true \

  fi
  ;;
esac

echo done.
