#!/bin/bash

## check arguments
if [ $# -ne 3 ]
 then
   echo ""
   echo "Invalid Azure Credentials"
   echo "Ex: create_secret.sh <AZURE_TENANT_ID> <AZURE_CLIENT_ID> <AZURE_CLIENT_SECRET>"
   echo ""
   exit 1
fi

## Azure  crecentials
AZURE_TENANT_ID=$1
AZURE_CLIENT_ID=$2
AZURE_CLIENT_SECRET=$3

## Creating namespace ##
kubectl create namespace linkedsecrets-system

## Creating secret for Azure credentials ##
kubectl create secret generic azure-credentials \
--from-literal=AZURE_TENANT_ID=$AZURE_TENANT_ID \
--from-literal=AZURE_CLIENT_ID=$AZURE_CLIENT_ID \
--from-literal=AZURE_CLIENT_SECRET=$AZURE_CLIENT_SECRET \
--namespace  linkedsecrets-system
