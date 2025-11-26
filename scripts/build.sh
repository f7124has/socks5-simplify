#!/bin/sh
docker rmi -f "tcpforwarder:latest" || true
docker build -t "tcpforwarder:latest" .