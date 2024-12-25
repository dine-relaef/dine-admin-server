#!/bin/bash

# Check if the first argument is passed
if [ -z "$1" ]; then
    echo "Please pass 'dev', 'test', or 'prod' as an argument."
    exit 1
fi

# Process environment argument
if [ "$1" == "dev" ]; then
    echo "Starting in development mode..."
    docker-compose -f docker-compose.dev.yml up --build
elif [ "$1" == "prod" ]; then
    echo "Starting in production mode..."
    docker-compose -f docker-compose.prod.yml up -d --build
elif [ "$1" == "test" ]; then
    echo "Starting in test mode..."

    # Step 1: Build the Docker images
    echo "Building Docker images..."
    docker-compose build

    # Step 2: Check if a container is using port 8080
    container_id=$(docker ps -q --filter 'publish=8080')

    # Step 3: If a container is found, stop and remove it
    if [ -n "$container_id" ]; then
        echo "Stopping container with ID: $container_id"
        docker stop $container_id
        echo "Removing container with ID: $container_id"
        docker rm $container_id
    else
        echo "No container using port 8080 found."
    fi

    # Step 4: Start the containers in detached mode
    echo "Starting containers in detached mode..."
    docker-compose up -d --build
else
    echo "Invalid argument. Please pass 'dev', 'test', or 'prod' as an argument."
    exit 1
fi
