#!/bin/bash

# Stop and remove any existing container
docker stop go-pong-server
docker rm go-pong-server

# Run the container
docker run -d --name go-pong-server -p 3000:3000 --entrypoint "./scripts/test.sh" go-pong-server