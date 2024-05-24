#!/bin/bash

# Build the Docker image
docker build -t myapp:latest .

# Check the build status
if [ $? -eq 0 ]; then
  echo "Build successful!"
else
  echo "Build failed!"
fi

