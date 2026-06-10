#!/bin/bash

XRAY_CORE=./xray-core
OUT=./gen

mkdir -p $OUT

find $XRAY_CORE -name '*.proto' | xargs protoc \
    -I $XRAY_CORE \
    --go_out $OUT --go_opt paths=source_relative \
    --go-grpc_out $OUT --go-grpc_opt paths=source_relative \
    --grpc-gateway_out $OUT --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true
