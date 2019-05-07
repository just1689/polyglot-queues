#!/bin/sh

# Tear down
docker stack rm stack

# Deploy stack
docker stack deploy stack --compose-file docker-compose.yml
