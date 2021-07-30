# Linkedsecret Installation

## Requirements

* Kubernetes 1.18 or newer cluster with RBAC (Role-Based Access Control) enabled is required.
* Helm 3
* Kubectl client installed and configured to access Kubernetes Cluster.

## Kubeideas Helm repo

Configure `kubeideas` Helm repository locally:

``` bash
helm repo add kubeideas https://kubeideas.github.io/linkedsecrets/
```

Search repository:

``` bash
helm search repo kubeideas
```

## Install Linkedsecrets Custom Resource Definitions (CRD)

Before install Linkedsecrets Helm chart, install manually Linkedsecrets CRD:

``` bash
kubectl apply -f https://github.com/kubeideas/linkedsecrets/releases/download/v0.7.0/security.kubeideas.io_linkedsecrets.yaml
```

## Enable GCP Secret Manager Access

* Create `Google Service account` with permission on role `Secret Manager Secret Accessor`.
* Create JSON key for the service account and save it locally.

**[IMPORTANT]** Avoid security issues and grant access only to secrets strictly relevant to your Kubernetes cluster.

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

**[IMPORTANT]** Avoid security issues and grant access only to secrets strictly relevant to your Kubernetes cluster.

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

**[IMPORTANT]** Avoid security issues and grant access only to secrets strictly relevant to your Kubernetes cluster.

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

Create IBM `Service Api Key` with permission on role `SecretsReader`.

``` bash
helm upgrade --install \
linkedsecrets \
--create-namespace=true \
--namespace=<LINKEDSECRETS-NAMESPACE> \
--set ibm.enabled=true \
--set ibm.ibmServiceApiKey="<IBM_SERVICE_API_KEY>" \
kubeideas/linkedsecrets
```

## Enable Mixed Cloud Secrets Solution Access

This example bellow enable Google Secret Manager and AWS Secrets Manager access, but any combination is allowed.

**[IMPORTANT]** Avoid security issues and grant access only to secrets strictly relevant to your Kubernetes cluster.

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

## Verifying installation

```bash
kubectl get pods -n <LINKEDSECRETS-NAMESPACE>
```

### Linkedsecret spec fields details

```bash
kubectl explain linkedsecret.spec
```

### Linkedsecret status fields details

```bash
kubectl explain linkedsecret.status
```

## Linkedsecrets commom fields

Commom fields are applicable for all supported Cloud providers:

``` yaml
apiVersion: security.kubeideas.io/v1
kind: LinkedSecret
metadata:
  name: <LINKEDSECRET-NAME>
spec:
  deployment: <DEPLOYMENT-NAME>
  keepSecretOnDelete: <true | false>
  secretName: <SECRET-NAME-CREATED-AND-MAINTAINED-BY-LINKEDSECRETS-ON-KUBERNETES>
  schedule: <"@every 10m" | ANY-OTHER-SYNCHRONIZATION-INTERVAL>
  suspended: <true | false>
```

## deployment

Set this field with deployment's name which secret is maintained by LinkedSecrets. If any change is detected, Linkedsecrets will rollout all deployment's pods automatically. This field can be omitted if you don't whant to use this feature.

## keepSecretOnDelete

Set this field to **`true`** if you want to keep secret intact after Linkedsecret object is deleted. This field can be omitted if you don't whant to use this feature.

This feature is particularly useful in upgrade situations.

## secretName

Set this field with Kubernetes Secret name you want. Linkedsecrets will create it with data retrieved from Cloud Secret.

## schedule

Linkedsecrets synchronization is executed based on schedule.
Pre-defined cron expressions and Classic cron expressions are accepted.

### Pre-defined Cron Expressions examples

| Expression       | Description                          |
|------------------|--------------------------------------|
| "@every 300s"    | Run every 5 minutes                  |
| "@every 10m"     | Run every 10 minutes                 |
| "@every 5m30s"   | Run every 5 minutes and 30 seconds   |
| "@hourly"        | Run once an hour, beginning of hour  |
| "@daily"         | Run once a day, midnight             |
|                  |                                      |

### Cron Expressions examples

| Expression       | Description                          |
|------------------|--------------------------------------|
| "*/20 * * * * *" | Run every 20 seconds                 |
| "0 */5 * * * *"  | Run every 5 minutes                  |
| "0 0 * * * *"    | Run once an hour, beginning of hour  |
| "0 0 0 * * *"    | Run once a day, midnight             |
|                  |                                      |

**[IMPORTANT]** Have in mind that Cloud will charge you on each secret created and access operations. Having said that, tune the schedule accordingly.

## suspended

Use this field any time you need to stop data synchronization between Kubernetes Secret and Cloud Secret.

## Cloud Provider specific instructions

Click [here](https://kubeideas.github.io/linkedsecrets/gcp) for GCP details and examples.

Click [here](https://kubeideas.github.io/linkedsecrets/aws) for AWS details and examples.

Click [here](https://kubeideas.github.io/linkedsecrets/azure) for Azure details and examples.

Click [here](https://kubeideas.github.io/linkedsecrets/ibm) for IBM details and examples.

## Uninstall Linkedsecrets

### **Remove linkedsecrets objects**

If you intend to keep applications secrets intact after remove Linkedsecrets objects, do not forget to enable option `"keepSecretOnDelete"` on all of them before.

``` bash
kubectl patch lns <NAME> --type='json' -p='[{"op": "replace", "path": "/spec/keepSecretOnDelete", "value":true}]' -n <NAMESPACE>
```

``` bash
kubectl delete lns --all --all-namespaces
```

### **Remove helm chart**

``` bash
helm -n <LINKEDSECRETS-NAMESPACE> delete linkedsecrets
```

### **Remove Custom Resource Definitions**

``` bash
kubectl delete -f https://github.com/kubeideas/linkedsecrets/releases/download/v0.7.0/security.kubeideas.io_linkedsecrets.yaml
```

### **Remove Linkedsecrets namespace**

``` bash
kubectl delete namespace <LINKEDSECRETS-NAMESPACE>
```
