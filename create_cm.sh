#!/usr/bin/env bash

number=$1
size=$2

content=$( dd if=/dev/urandom bs=1K count=$2 | base64 )

for ((n=0;n<$1;n++))
do

  cat <<EOF | kubectl create -f -
apiVersion: v1
data:
  poop: "$content"
kind: ConfigMap
metadata:
  name: thingy-sample-$n
EOF

done
