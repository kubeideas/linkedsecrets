# IBM Instructions and Examples

## IBM Secret type

Linkedsecrets support `"Arbitrary Secrets"` only.

This kind of secret support ratation but not versioning.

## IBM Secret data format

Linkedsecrets support `"PLAIN"` and `"JSON"` formats.

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

Follow bellow all spec fields supported by Linkedsecrets when using IBM Secret Manager:

``` yaml
apiVersion: security.kubeideas.io/v1
kind: LinkedSecret
metadata:
  name: <LINKEDSECRET-NAME>
spec:
  deployment: <DEPLOYMENT-NAME>
  keepSecretOnDelete: <true | false>
  provider: IBM
  providerDataFormat: <JSON | PLAIN>
  providerOptions:
    secretManagerInstanceId: <SECRET-MANAGER-INSTANCE-UUID>
    secretId: <SECRET-UUID>
    region: <SECRET-MANAGER-REGION>
  secretName: <SECRET-NAME-CREATED-AND-MAINTAINED-BY-LINKEDSECRETS-ON-KUBERNETES>
  schedule: <"@every 10m" | ANY-OTHER-SYNCHRONIZATION-INTERVAL>
  suspended: <true | false>
```

## Examples

Click [Here](https://kubeideas.github.io/linkedsecrets/examples.zip) and get them.

## References

[IBM Secret Manager](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-getting-started)

[IBM Secret Manager API](https://cloud.ibm.com/apidocs/secrets-manager?code=go#create-secret)
