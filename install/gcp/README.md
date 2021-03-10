# Linkedsecrets installation

Before installing Linkedsecrets operator it is necessary create a `Google Service account` with the following details:

* role `Secret Manager Secret Accessor` permission.
* Create Json key file and save with name `gcp-credentials.json` in this directory

**[IMPORTANT]** Have in mind to grant access only to secrets strictly relevants to your Kubernetes cluster project.

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

## Google Secret format 
Linkedsecret support `"plain"` format and `"json"` format.

### Plain format
This format will use "=" to separate key/value. White spaces and white lines are allowed and will be skipped during payload parse.

Example:
```bash
username = admin
password=teste123

host = myhost01
```

### Json format
This format support a simple key/value json.

Example:
```bash
{
  "username" : "admin",
  "password" : "teste123",
  "host" : "myhost01"
}
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

**[IMPORTANT]** Have in mind that Google cloud will charge you based on secret access. Having said that, tune the schedule accordingly.

## References
[Google Secret Manager](https://cloud.google.com/secret-manager/docs/configuring-secret-manager)