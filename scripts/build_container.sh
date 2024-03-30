#!/bin/bash

# Set the Docker image tag
TAG="go-pong-server"

# Build the Docker image using docker-buildx
docker buildx build -t $TAG .

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Docker image build successful!"
else
    echo "Docker image build failed."
fi