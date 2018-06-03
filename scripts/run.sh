#!/usr/bin/env bash

echo "*** Everything must be stoped ***"
docker kill $(docker container ls -q)

#echo "*** Ensure everything is clean and ready ***"
#sleep 1s
#docker-compose rm -f

#echo "*** Remove all unused volumes ***"
#docker volume prune -f

#echo "*** Remove all stopped containers ***"
#docker container prune -f

echo "*** Build and run! ***"
sleep 1s
docker-compose build
docker-compose up -d