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

Accepted formats:
* Pre-defined cron expressions. Ex: "@every 10minutes"
* 6 field format cron expressions. Ex: "*/20 * * * * *"

**[IMPORTANT]** Have in mind that Google cloud will charge you based on secret access. Having said that, tune the schedule accordingly.

## References
[Google Secret Manager](https://cloud.google.com/secret-manager/docs/configuring-secret-manager)