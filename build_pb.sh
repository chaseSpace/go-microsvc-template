#!/usr/bin/env bash


__dir=$1  # 指定生成的pb子目录，比如user 对应 microsvc/svc/user


case $__dir in
all)
    protoc -I ./proto/ --go_out ./protocol/ --go_opt paths=source_relative \
      proto/svc/*/*.proto proto/svc/*.proto
    ;;
*)
    protoc -I ./proto/ --go_out ./protocol/ --go_opt paths=source_relative \
      proto/svc/$__dir/*.proto proto/svc/*.proto
  ;;
esac
