# Vault token renovator

## Basic usage

### Create Secret with configuration

```json
{
  "tokens": [
    {
      "name": "token00",
      "token": "3SYG............."
    },
    {
      "name": "token01",
      "token": "3SYF............."
    }
  ]
}
```

```bash
kubectl create secret generic renovator-config --from-file=config.json=local-file.json -n production
```

### Create CronJob resource

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: vault-renovator
spec:
  schedule: "0 8 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: renovator
              image: ackee/renovator:latest
              volumeMounts:
                - name: renovator-config-volume
                  mountPath: /etc/renovator/
                  readOnly: true
              env:
                - name: VAULT_ADDRESS
                  value: https://your.vault.co.uk
                - name: INSECURE
                  value: 'false'
                - name: DEBUG
                  value: 'false'
                - name: TTL_THRESHOLD
                  value: '15206400'
                - name: TTL_INCREMENT
                  value: '5184000'
                - name: CONFIG_FILE_PATH
                  value: '/etc/renovator/config.json'
                - name: SLACK_WEBHOOK_URL
                  value: 'https://hooks.slack.com/services/....'
          volumes:
            - name: renovator-config-volume
              secret:
                secretName: renovator-config
          restartPolicy: OnFailure
```

```bash
kubectl apply -f cronjob.yaml -n production
```