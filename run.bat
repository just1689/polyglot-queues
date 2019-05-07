#!/bin/sh

# Tear down
docker service rm Numero1
docker service rm Numero2

# Build container
docker build -t pq/pq:pq .


# Deploy stack
docker service create     --name Numero1     --network my-services     pq/pq:pq --name Numero1
docker service create     --name Numero2     --network my-services     pq/pq:pq --name Numero2