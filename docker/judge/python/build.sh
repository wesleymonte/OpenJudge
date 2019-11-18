#!/bin/sh

if [[ "$#" -ne 1 ]]; then
  echo "Usage: $0 <docker tag>"
  exit 1
fi

tag=$1

sudo docker build -t wesleymonte/judge-python:${tag} -f docker/judge/python/Dockerfile .