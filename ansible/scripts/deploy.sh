#!/bin/bash

readonly IMAGE_NAME=wesleymonte/openjudge
readonly TAG=latest
readonly CONTAINER_NAME=openjudge

if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    docker stop $CONTAINER_NAME
fi

if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    docker rm $CONTAINER_NAME
fi

# shellcheck disable=SC2046
docker run -itd -v $(pwd)/.env:/service/.env \
                -v /var/run/docker.sock:/var/run/docker.sock \
                -v /usr/bin/docker:/usr/bin/docker \
                --net=host \
                --name $CONTAINER_NAME $IMAGE_NAME:$TAG