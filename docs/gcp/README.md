# GCP Instructions and Examples

## Google Secret Manager data format

Linkedsecrets support `"PLAIN"` format and `"JSON"` format.

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
  rolloutRestartDeploy: <DEPLOYMENT-NAME>
  keepSecretOnDelete: <true | false>
  provider: Google
  providerSecretFormat: <JSON | PLAIN>
  providerOptions:
    project: <GCP-PROJECT-ID>
    secret: <GCP-SECRET-NAME>
    version: <latest | "1" | "2" | ...>  
  secretName: <KUBERNETES-SECRET-NAME-CREATED-AND-MAINTAINED-BY-LINKEDSECRETS>
  schedule: <"@every 10m" | ANY-OTHER-SYNCHRONIZATION-INTERVAL>
  suspended: <true | false>
```

**[IMPORTANT]** Secret latest version will be used if field version is omitted.

## Examples

Click [Here](https://kubeideas.github.io/linkedsecrets/examples.zip) and get them.

## References

[Google Secret Manager](https://cloud.google.com/secret-manager/docs/configuring-secret-manager)
