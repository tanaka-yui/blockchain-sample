#!/bin/bash

BASE_DIR=$(cd $(dirname $0); pwd)

COMMAND=${1:-"help"}

block_node1() {
  go run $BASE_DIR/app/blockchain_node_server/main.go --addr=0.0.0.0:5001
}

block_node2() {
  go run $BASE_DIR/app/blockchain_node_server/main.go --addr=0.0.0.0:5002
}

block_node3() {
  go run $BASE_DIR/app/blockchain_node_server/main.go --addr=0.0.0.0:5002
}

wallet1() {
  go run $BASE_DIR/app/wallet_server/main.go --addr=0.0.0.0:8080
}

wallet2() {
  go run $BASE_DIR/app/wallet_server/main.go --addr=0.0.0.0:8081
}

test() {
  go test -v ./test/...
}

bench() {
  go test -bench . -benchmem ./test/...
}

case "${COMMAND}" in
  block_node1) block_node1 ;;
  block_node2) block_node2 ;;
  block_node3) block_node3 ;;
  wallet1) wallet1 ;;
  wallet2) wallet2 ;;
  test) test ;;
  bench) bench ;;
esac
