# LinkedSecrets Helm chart

This document describes how to install and remove Linkedsecrets CRD and helm chart.

## Requirements

* Kubernetes 1.18 or newer cluster with RBAC (Role-Based Access Control) enabled is required.
* Helm 3
* Kubectl client installed and configured to access Kubernetes Cluster.

## Install Linkedsecrets Custom Resource Definitions (CRD)

Before install Linkedsecrets Helm chart, install manually Linkedsecrets CRD:

``` bash
kubectl apply -f https://github.com/kubeideas/linkedsecrets/releases/download/v0.8.4/security.kubeideas.io_linkedsecrets.yaml
```

## Enable GCP Secret Manager Access

Create a `Google Service account` with the following details:

* Role `Secret Manager Secret Accessor` permission.
* Create JSON key file and save it locally.

**[IMPORTANT]** Avoid security issues and grant access only to secrets strictly relevant to your Kubernetes cluster project.

``` bash
helm upgrade --install \
linkedsecrets \
--create-namespace=true \
--namespace=<LINKEDSECRETS-NAMESPACE> \
--set gcp.enabled=true \
--set-file gcp.credentialFile="path/<GCP_CREDENTIALS_FILE>.json" \
kubeideas/linkedsecrets
```

## Enable AWS Secrets manager Access

Create AWS user with permissions to read secrets.

**[IMPORTANT]** Avoid security issues and grant access only to secrets strictly relevant to your Kubernetes cluster project.

``` bash
helm upgrade --install \
linkedsecrets \
--create-namespace=true \
--namespace=<LINKEDSECRETS-NAMESPACE> \
--set aws.enabled=true \
--set aws.awsAccessKeyId="<AWS_ACCESS_KEY_ID>" \
--set aws.awsSecretAccessKey="<AWS_SECRET_ACCESS_KEY>" \
kubeideas/linkedsecrets
```

## Enable Azure Keyvault Access

Register an App on Azure Active directory with permissions to get and list secrets.

**[IMPORTANT]** Avoid security issues and grant access only to secrets strictly relevant to your Kubernetes cluster project.

``` bash
helm upgrade --install \
linkedsecrets \
--create-namespace=true \
--namespace=<LINKEDSECRETS-NAMESPACE> \
--set azure.enabled=true \
--set azure.azureTenantId="<AZURE_TENANT_ID>" \
--set azure.azureClientId="<AZURE_CLIENT_ID>" \
--set azure.azureClientSecret="<AZURE_CLIENT_SECRET>" \
kubeideas/linkedsecrets
```

## Enable IBM Secrets Manager Access

Create an `IBM Service Api Key` with the following details:

* Role `SecretsReader`.

``` bash
helm upgrade --install \
linkedsecrets \
--create-namespace=true \
--namespace=<LINKEDSECRETS-NAMESPACE> \
--set ibm.enabled=true \
--set ibm.ibmServiceApiKey="<IBM_SERVICE_API_KEY>" \
kubeideas/linkedsecrets
```

## Enable mixed Cloud solutions Access

This example bellow enables Google Secret Manager and AWS Secrets Manager access, but any combination is allowed.

``` bash
helm upgrade --install \
linkedsecrets \
--create-namespace=true \
--namespace=<LINKEDSECRETS-NAMESPACE> \
--set gcp.enabled=true \
--set-file gcp.credentialFile="path/<GCP_CREDENTIALS_FILE>.json" \
--set aws.enabled=true \
--set aws.awsAccessKeyId="<AWS_ACCESS_KEY_ID>" \
--set aws.awsSecretAccessKey="<AWS_SECRET_ACCESS_KEY>" \
kubeideas/linkedsecrets
```

## Uninstall Linkedsecrets

### Remove linkedsecrets objects

If you intend to keep applications secrets intact after remove Linkedsecrets objects, do not forget to enable option `"keepSecretOnDelete"` on all of them before.

``` bash
kubectl patch lns <NAME> --type='json' -p='[{"op": "replace", "path": "/spec/keepSecretOnDelete", "value":true}]' -n <NAMESPACE>
```

``` bash
Kubectl delete lns --all --all-namespaces
```

### Remove helm chart

``` bash
helm -n <LINKEDSECRETS-NAMESPACE> delete linkedsecrets
```

### Remove Custom Resource Definitions

``` bash
kubectl delete -f https://github.com/kubeideas/linkedsecrets/releases/download/v0.8.1/security.kubeideas.io_linkedsecrets.yaml
```

### Remove Linkedsecrets namespace

``` bash
kubectl delete namespace <LINKEDSECRETS-NAMESPACE>
```
