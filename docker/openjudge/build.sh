#!/bin/sh

if [[ "$#" -ne 1 ]]; then
  echo "Usage: $0 <docker tag>"
  exit 1
fi

tag=$1

go build -o openjudge

sudo docker build -t wesleymonte/openjudge:${tag} -f docker/openjudge/Dockerfile .