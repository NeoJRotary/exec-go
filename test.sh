#!/bin/bash
set -e

echo "Run All Tests in Docker"
docker build -t neojrotary/exec-go/test -f test.Dockerfile .
docker rmi neojrotary/exec-go/test
