#!/bin/bash

readonly IMAGE_NAME=wesleymonte/openjudge
readonly TAG=latest
readonly CONTAINER_NAME=openjudge

# shellcheck disable=SC2046
docker run -itd -v $(pwd)/.env:/service/.env \
                -v /var/run/docker.sock:/var/run/docker.sock \
                -v /usr/bin/docker:/usr/bin/docker \
                --net=host \
                --name $CONTAINER_NAME $IMAGE_NAME:$TAG