#!/bin/bash

## check arguments
if [ $# -ne 2 ]
 then
   echo ""
   echo "Invalid AWS Credentials"
   echo "Ex: create_secret.sh <AWS_ACCESS_KEY_ID> <AWS_SECRET_ACCESS_KEY>"
   echo ""
   exit 1
fi

## AWS  crecentials
AWS_ACCESS_KEY_ID=$1
AWS_SECRET_ACCESS_KEY=$2

## Creating namespace ##
kubectl create namespace linkedsecrets-system

## Creating secret for aws credentials ##
kubectl create secret generic aws-credentials \
--from-literal=AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
--from-literal=AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
--namespace  linkedsecrets-system
