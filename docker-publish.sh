#!/bin/bash -e
docker build -t pivotalservices/file-downloader-resource .
docker push pivotalservices/file-downloader-resource
