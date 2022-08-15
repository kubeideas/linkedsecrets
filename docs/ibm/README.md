# IBM Instructions and Examples

## IBM Secret type

Linkedsecrets support `"Arbitrary Secrets"` only.

This kind of secret support rotation but not versioning.

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
  rolloutRestartDeploy: <DEPLOYMENT-NAME>
  keepSecretOnDelete: <true | false>
  provider: IBM
  providerSecretFormat: <JSON | PLAIN>
  providerOptions:
    secretManagerInstanceId: <SECRET-MANAGER-INSTANCE-UUID>
    secretId: <SECRET-UUID>
    region: <SECRET-MANAGER-REGION>
  secretName: <KUBERNETES-SECRET-NAME-CREATED-AND-MAINTAINED-BY-LINKEDSECRETS>
  schedule: <"@every 10m" | ANY-OTHER-SYNCHRONIZATION-INTERVAL>
  suspended: <true | false>
```

## Examples

Click [linkedsecret_json_example1.yaml](https://kubeideas.github.io/linkedsecrets/ibm/examples/linkedsecret_json_example1.yaml).

Click [linkedsecret_json_example2.yaml](https://kubeideas.github.io/linkedsecrets/ibm/examples/linkedsecret_json_example2.yaml).

Click [linkedsecret_plain_example1.yaml](https://kubeideas.github.io/linkedsecrets/ibm/examples/linkedsecret_plain_example1.yaml).

Click [linkedsecret_rollout_restart_deploy.yaml](https://kubeideas.github.io/linkedsecrets/ibm/examples/linkedsecret_rollout_restart_deploy.yaml).

## References

[IBM Secret Manager](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-getting-started)

[IBM Secret Manager API](https://cloud.ibm.com/apidocs/secrets-manager?code=go#create-secret)
