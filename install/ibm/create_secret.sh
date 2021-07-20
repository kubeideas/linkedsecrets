#!/bin/bash

## check arguments
if [ $# -ne 1 ]
 then
   echo ""
   echo "Invalid IBM Service API key"
   echo "Ex: create_secret.sh <SERVICE_API_KEY>"
   echo ""
   exit 1
fi

## IBM  crecentials
IBM_SERVICE_API_KEY=$1


## Creating namespace ##
kubectl create namespace linkedsecrets-system

## Creating secret for ibm credentials ##
kubectl create secret generic ibm-credentials \
--from-literal=IBM_SERVICE_API_KEY=$IBM_SERVICE_API_KEY \
--namespace  linkedsecrets-system
