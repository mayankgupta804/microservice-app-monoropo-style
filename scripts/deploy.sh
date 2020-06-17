#!/bin/bash

echo "Spinning up one docker-machine instance..."

docker-machine create node-1

echo "Initializing Swarm mode..."

docker-machine ssh node-1 -- docker swarm init --advertise-addr $(docker-machine ip node-1)

echo "Deploying the Incident service..."

docker stack deploy --compose-file=../docker-compose.yml sq

echo "Get the IP address..."
eval $(docker-machine env node-1)
docker-machine ip $(docker service ps -f "desired-state=running" --format "{{.Node}}" sq_web)