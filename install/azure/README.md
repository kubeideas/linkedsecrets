# Linkedsecrets installation

## Requirements

* **AZURE_TENANT_ID**, **AZURE_CLIENT_ID** and **AZURE_CLIENT_SECRET** with permissions to get and list secrets on Azure Keyvault.

**[IMPORTANT]** Avoid security issues and grant access only to secrets strictly relevant to your Kubernetes cluster project.

## Namespace and Azure credentials secret

```bash
./create_secret.sh
```

## CRD's and controller

```bash
kubectl apply -f install-linkedsecret-azure.yaml
```

## Verifying installation

```bash
kubectl get pods -n linkedsecrets-system
```

### Linkedsecret spec fields details

```bash
kubectl explain linkedsecret.spec
```

### Linkedsecret status fields details

```bash
kubectl explain linkedsecret.status
```

## Azure Keyvault Secrets data format

Linkedsecret support `"PLAIN"` format and `"JSON"` format.

### PLAIN format

This format must use "=" to separate key/value. White spaces and white lines are allowed and will be skipped during payload parse.
PLAIN secret creation can be done in Azure Console or using Azure CLI.

Example:

Create resource group:

```bash
az group create --name "kubernetes" --location "EastUS" 
```

Create keyvault:

```bash

az keyvault create --name "lnsvault" --resource-group "kubernetes" --location "EastUS" 
```

Create file `[mysecret.txt]` with PLAIN text:

```bash
username = admin
password=teste123

host = myhost01
```

Create a secret with file `[mysecret.txt]` :

```bash
az keyvault secret set --vault-name "lnsvault" --name "mysecret" --file "./mysecret.txt"
```

### JSON format

JSON secret creation can be done in Azure Console or using Azure CLI.

```bash
az group create --name "kubernetes" --location "EastUS" 
```

Create keyvault:

```bash

az keyvault create --name "lnsvault" --resource-group "kubernetes" --location "EastUS" 
```

Create file `[mysecret.txt]` with json text:

```bash
{
  "username" : "admin",
  "password" : "teste123",
  "host" : "myhost01"
}
```

Create a secret with encoded file `[mysecret.txt]`:

```bash
az keyvault secret set --vault-name "lnsvault" --name "mysecret" --file "./mysecret.txt"
```

## Linkedsecrets Spec Fields

Follow bellow all spec fields supported by Linkedsecrets when using Azure Keyvault Secrets:

``` yaml
apiVersion: security.kubeideas.io/v1
kind: LinkedSecret
metadata:
  name: <LINKEDSECRET-NAME>
spec:
  deployment: <DEPLOYMENT-NAME>
  keepSecretOnDelete: <true | false>
  provider: Azure
  providerDataFormat: <JSON | PLAIN>
  providerOptions:
    keyvault: <AZURE-KEYVAULT-NAME>
    secret: <AZURE-SECRET-NAME>
    version: <AWZURE-SECRET-VERSION-ID> 
  secretName: <SECRET-NAME-CREATED-AND-MAINTAINED-BY-LINKEDSECRETS-ON-KUBERNETES>
  schedule: <"@every 10m" | ANY-OTHER-SYNCHRONIZATION-INTERVAL>
  suspended: <true | false>
```

**[IMPORTANT]** Secret latest version will be used if field version is omitted.

### Deployment Field

Set this field with deployment name which use the secret maintained by LinkedSecrets. If any change is detected, all deployment pods will be restarted in order to use the new secret data. This field can be omitted if you don't whant to use this feature.

### keepSecretOnDelete Field

Set this field to **`true`** if you want to keep secret after linkedsecret has been deleted. This field can be omitted if you don't whant to use this feature.

This feature is particularly useful in upgrade situations.

### SecretName Field

This is used by Linkedsecrets to create a Kubernetes Secret with data retrieved from Secrets provider.

### Schedule Field

Linkedsecret supports synchronization based on schedule.
Pre-defined cron expressions and Classic cron expressions are accepted.

#### Pre-defined Cron Expressions examples

| Expression       | Description                          |
|------------------|--------------------------------------|
| "@every 300s"    | Run every 5 minutes                  |
| "@every 10m"     | Run every 10 minutes                 |
| "@every 5m30s"   | Run every 5 minutes and 30 seconds   |
| "@hourly"        | Run once an hour, beginning of hour  |
| "@daily"         | Run once a day, midnight             |
|                  |                                      |

#### Cron Expressions examples

| Expression       | Description                          |
|------------------|--------------------------------------|
| "*/20 * * * * *" | Run every 20 seconds                 |
| "0 */5 * * * *"  | Run every 5 minutes                  |
| "0 0 * * * *"    | Run once an hour, beginning of hour  |
| "0 0 0 * * *"    | Run once a day, midnight             |
|                  |                                      |

**[IMPORTANT]** Have in mind that Azure cloud will charge you on each secret created and access operations. Having said that, tune the schedule accordingly.

### Suspended Field

Use this field any time you need to stop data synchronizatin between Kubernetes Secret and Secrets Provider.

## References

[Azure Keyvault Secrets](https://docs.microsoft.com/en-us/azure/key-vault/)
