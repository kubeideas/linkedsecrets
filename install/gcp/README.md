# Linkedsecrets installation

Before installing Linkedsecrets operator it is necessary to create a `Google Service account` with the following details:

* Role `Secret Manager Secret Accessor` permission.
* Create JSON key file and name it as `gcp-credentials.json` in this directory.

**[IMPORTANT]** Avoid security issues and grant access only to secrets strictly relevant to your Kubernetes cluster project.

## Namespace and GCP credentials secret

```bash
./create_secret.sh
```

## CRD's and controller

```bash
kubectl apply -f install-linkedsecret-gcp.yaml
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

## Google Secret Manager data format

Linkedsecret support `"PLAIN"` format and `"JSON"` format.

### PLAIN format

This format must use "=" to separate key/value. White spaces and white lines are allowed and will be skipped during payload parse.

Example:

```bash
username = admin
password=teste123

host = myhost01
```

### JSON format

This format support a simple key/value JSON.

Example:

```bash
{
  "username" : "admin",
  "password" : "teste123",
  "host" : "myhost01"
}
```

## Linkedsecrets Spec Fields

Follow bellow all spec fields supported by Linkedsecrets when using Google Secret Manager:

``` yaml
apiVersion: security.kubeideas.io/v1
kind: LinkedSecret
metadata:
  name: <LINKEDSECRET-NAME>
spec:
  deployment: <DEPLOYMENT-NAME>
  keepSecretOnDelete: <true | false>
  provider: Google
  providerDataFormat: <JSON | PLAIN>
  providerOptions:
    project: <GCP-PROJECT-ID>
    secret: <GCP-SECRET-NAME>
    version: <latest | ANY-OTHER-VERSION>  
  secretName: <SECRET-NAME-CREATED-AND-MAINTAINED-BY-LINKEDSECRETS>
  schedule: <"@every 10m" | ANY-OTHER-SYNCHRONIZATION-INTERVAL>
  suspended: <true | false>
```

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
| "0 0 * * * *"    | Run once an hour, beginning of hour  |
| "0 0 0 * * *"    | Run once a day, midnight             |
|                  |                                      |

**[IMPORTANT]** Have in mind that Google cloud will charge you based on secret access. Having said that, tune the schedule accordingly.

### Suspended Field

Use this field any time you need to stop data synchronizatin between Kubernetes Secret and Secrets Provider.

## References

[Google Secret Manager](https://cloud.google.com/secret-manager/docs/configuring-secret-manager)
