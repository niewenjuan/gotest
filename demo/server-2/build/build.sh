#!/bin/bash
set -e
set -x

IMAGE=$1
TAG=$2

BINARY_NAME="app"
PKG_NAME="service.tar.gz"

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)

cd $ROOT_PATH
if [ -f $BINARY_NAME ]; then
    rm $BINARY_NAME
fi

if [ -f $PKG_NAME ]; then
    rm $PKG_NAME
fi

go build -a -o "app"
tar -zcvf $PKG_NAME conf $BINARY_NAME start.sh

docker build -t $IMAGE:$TAG .
docker save -o $IMAGE.tar $IMAGE:$TAG
docker rmi -f $IMAGE:$TAG
ddd=$(docker images|grep $IMAGE|awk '{print $3}'|awk 'NR==2{print}');
echo $ddd
rm -rf app*
rm -rf $PKG_NAME

echo "Build success!"

