# Azure Instructions and Examples

## Azure Keyvault Secrets data format

Linkedsecrets support `"PLAIN"` format and `"JSON"` format.

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

## Examples

Click [HERE](https://kubeideas.github.io/linkedsecrets/azure/examples.zip) and get them.

## References

[Azure Keyvault Secrets](https://docs.microsoft.com/en-us/azure/key-vault/)
