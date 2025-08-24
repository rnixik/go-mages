#!/usr/bin/env bash

echo $(( $(cat version) + 1 )) > version
git pull && docker-compose -p mages -f deployments/docker-compose.prod.yml up -d --build
