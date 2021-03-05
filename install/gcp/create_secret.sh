#!/bin/bash

# check for gcp-credentials.json file 
if [ ! -f ./gcp-credentials.json ]
  then
    echo -e "\n\n"
    echo "File gcp-credentials.json with GCP credentials must exists in this directory."
    echo "Please create a Google Service Account and grant \"Secret Manager Secret Accessor\" role properly."
    echo "Finnaly generate Service Account json file key and save with name \"gcp-credentials.json\" in this directory."
    echo -e "\n\n"
    exit 1
fi

## Creating namespace ##
kubectl create namespace linkedsecrets-system

## Creating secret with json crecentials ##
kubectl create secret generic gcp-credentials --from-file=gcp-credentials.json --namespace  linkedsecrets-system
