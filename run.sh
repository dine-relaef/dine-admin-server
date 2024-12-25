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
    docker-compose up -d --build
else
    echo "Invalid argument. Please pass 'dev', 'test', or 'prod' as an argument."
    exit 1
fi
