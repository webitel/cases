# !/bin/sh
#set -x

src=pkg/proto/cases
dst=api

# Ensure target dir exists
mkdir -p $dst

protoc -I pkg/proto \
  --go_opt=paths=source_relative --go_out=$dst \
  --go-grpc_out=require_unimplemented_servers=false,paths=source_relative:$dst \
  $src/*.proto

res=$? # Last command execution
if [ $res -ne 0 ]; then
  >&2 echo "[ERR]: protoc: failed with exit code ${res}."
  exit $res
fi

echo "[INFO]: protoc: successfully generated Proto files."
