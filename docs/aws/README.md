# AWS Instructions and Examples

## AWS Secrets Manager data format

Linkedsecrets support `"PLAIN"` format and `"JSON"` format.

### PLAIN format

This format must use "=" to separate key/value. White spaces and white lines are allowed and will be skipped during payload parse. AWS console stores secrets data only in the **JSON** format and  you will have to use AWS CLI to store secret as binary to use **PLAIN** format.

Example:

Create file `[mysecret.txt]` with PLAIN text:

```bash
username = admin
password=teste123

host = myhost01
```

Create AWS secret with encoded file `[mysecret.txt]` using AWS CLI:

```bash
aws secretsmanager create-secret --name mysecret --secret-binary fileb://mysecret.txt
```

### JSON format

JSON secret creation can be done in AWS Console or using AWS CLI.

AWS CLI example:

Create file `[mysecret.txt]` with json text:

```bash
{
  "username" : "admin",
  "password" : "teste123",
  "host" : "myhost01"
}
```

Create AWS secret with encoded file `[mysecret.txt]` using AWS CLI:

```bash
aws secretsmanager create-secret --name mysecret --secret-binary fileb://mysecret.txt
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
    version: <AWSCURRENT | AWSPREVIOUS> 
  secretName: <SECRET-NAME-CREATED-AND-MAINTAINED-BY-LINKEDSECRETS-ON-KUBERNETES>
  schedule: <"@every 10m" | ANY-OTHER-SYNCHRONIZATION-INTERVAL>
  suspended: <true | false>
```

**[IMPORTANT]** Secret latest version will be used if field version is omitted.

## Examples

Click [HERE](https://kubeideas.github.io/linkedsecrets/aws/examples.zip) and get them.

## References

[AWS Secret Manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/getting-started.html)
