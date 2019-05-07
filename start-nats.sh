#!/bin/sh

docker service rm nats-cluster-node-1
docker service rm nats-cluster-node-2

docker service create --network my-services --name nats-cluster-node-1 nats:1.0.0 -cluster nats://0.0.0.0:6222 -DV
docker service create --network my-services --name nats-cluster-node-2 nats:1.0.0 -cluster nats://0.0.0.0:6222 -routes nats://nats-cluster-node-1:6222 -DV