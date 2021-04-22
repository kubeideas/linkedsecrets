# Linkedsecrets installation

## Requirements

* **AWS_ACCESS_KEY_ID** and **AWS_SECRET_ACCESS_KEY** with permissions to read secrets on AWS Secret manager.

**[IMPORTANT]** Avoid security issues and grant access only to secrets strictly relevant to your Kubernetes cluster project.

## Namespace and AWS credentials secret

```bash
./create_secret.sh
```

## CRD's and controller

```bash
kubectl apply -f install-linkedsecret-aws.yaml
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

## AWS Secrets Manager data format

Linkedsecret support `"PLAIN"` format and `"JSON"` format.

### PLAIN format

This format must use "=" to separate key/value. White spaces and white lines are allowed and will be skipped during payload parse. AWS console stores secrets data only in the **JSON** format and  you will have to use AWS CLI to store secret as binary to use **PLAIN** format.

Example:

Create file `[mysecret.txt]` with PLAIN text:

```bash
username = admin
password=teste123

host = myhost01
```

Encode file `[mysecret.txt]`:

```bash
cat mysecret.txt |base64 > mysecret-encoded.txt
```

Create a secret with encoded file `[mysecret-encoded.txt]` using AWS CLI:

```bash
aws secretsmanager create-secret --name mysecret --secret-binary=fileb://mysecret-encoded.txt
```

### JSON format

JSON secret creation can be done in AWS Console or using AWS CLI.

AWS CLI example:

Create file `[mysecret.txt]` with plain text:

```bash
{
  "username" : "admin",
  "password" : "teste123",
  "host" : "myhost01"
}
```

Create a secret with encoded file `[mysecret.txt]` using AWS CLI:

```bash
aws secretsmanager create-secret --name mysecret --secret-binary=file://mysecret.txt
```

## Linkedsecrets Spec Fields

Follow bellow all spec fields supported by Linkedsecrets when using AWS Secrets Manager:

``` yaml
apiVersion: security.kubeideas.io/v1
kind: LinkedSecret
metadata:
  name: <LINKEDSECRET-NAME>
spec:
  deployment: <DEPLOYMENT-NAME>
  keepSecretOnDelete: <true | false>
  provider: AWS
  providerDataFormat: <JSON | PLAIN>
  providerOptions:
    secret: <AWS-SECRET-NAME>
    region: <AWS-SECRET-RESOURCE-REGION>
    version: <AWSCURRENT  | ANY-OTHER-VERSION> 
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
| "0 0 0 * * *"    | Run once a day, midnight             |
|                  |                                      |

**[IMPORTANT]** Have in mind that AWS cloud will charge you based on secret access. Having said that, tune the schedule accordingly.

### Suspended Field

Use this field any time you need to stop data synchronizatin between Kubernetes Secret and Secrets Provider.

## References

[AWS Secret Manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/getting-started.html)
