# Linkedsecrets installation

## Requirements

* AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY with permissions to read secrets on AWS Secret manager

**[IMPORTANT]** Have in mind to grant access only to secrets strictly relevants to your Kubernetes cluster project.

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

## AWS Secret format 
Linkedsecret support `"plain"` format and `"json"` format.

### Plain format
This format will use "=" to separate key/value. White spaces and white lines are allowed and will be skipped during payload parse. As AWS console stores data only in the `"SecretString"`, you will have to use AWS CLI to store secret as binary.

Example:

Create file `[mysecret.txt]` with plain text:
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


### Json format
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

## Schedule
Linkedsecret supports synchronization based on schedule. 
Pre-defined cron expressions and Classic cron expressions are accepted.

### Pre-defined Cron Expressions examples:
| Expression       | Description                          |
|------------------|--------------------------------------|
| "@every 300s"    | Run every 5 minutes                  |
| "@every 10m"     | Run every 10 minutes                 | 
| "@every 5m30s"   | Run every 5 minutes and 30 seconds   |
| "@hourly"        | Run once an hour, beginning of hour  |
| "@daily"         | Run once a day, midnight             |
|                  |                                      |

### Cron Expressions examples:

| Expression       | Description                          |
|------------------|--------------------------------------|
| "*/20 * * * * *" | Run every 20 seconds                 |
| "0 */5 * * * *"  | Run every 5 minutes                  |
| "0 0 * * * *"    | Run once an hour, beginning of hour  | 
| "0 0 0 * * *"    | Run once a day, midnight             |
|                  |                                      |

**[IMPORTANT]** Have in mind that AWS cloud will charge you based on secret access. Having said that, tune the schedule accordingly.

## References
[AWS Secret Manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/getting-started.html)