#!/usr/bin/env bash
set -e

# 此脚本支持在linux、mac环境下执行
# - mac环境请在根目录下新建 tool_mac目录（不要提交到git），并自行下载tool/目录下的二进制文件，脚本会自动调用

_pb_dir=$1 # 指定生成的pb子目录，比如user 对应 microsvc/svc/user
_bin_path='./tool/protoc_v24'
if [[ $(arch) = 'arm64' ]]; then
  echo '*** within mac env ***'
  _bin_path='./tool_mac/protoc_v24'
fi

chmod +x $_bin_path/*

PATH=$PATH:$_bin_path

OUTPUT_DIR="./protocol/"

case $_pb_dir in
clear)
  rm -rf ./protocol/*
  ;;
*)
  if [[ -n $_pb_dir ]]; then
    echo "regenerate proto/svc/$_pb_dir/..."
    rm -rf ./protocol/$_pb_dir/*

    $_bin_path/protoc -I ./proto/ --go_out=$OUTPUT_DIR --go_opt=paths=source_relative \
      --go-grpc_out=$OUTPUT_DIR --go-grpc_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false \
      proto/svc/$_pb_dir/*.proto proto/svc/*.proto
  else
    echo "regenerate all proto files..."
    rm -rf ./protocol/*

    $_bin_path/protoc -I ./proto/ --go_out=$OUTPUT_DIR --go_opt paths=source_relative \
      --go-grpc_out=$OUTPUT_DIR --go-grpc_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false \
      proto/svc/*/*.proto proto/svc/*.proto
  fi
  ;;
esac

echo done.
