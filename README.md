[![CircleCI](https://circleci.com/gh/giantswarm/{APP-NAME}-app.svg?style=shield)](https://circleci.com/gh/giantswarm/{APP-NAME}-app)

# k8s-jwt-to-vault-token

Docker image which uses service account JWT to request vault token and save into secret.

## Usage

`k8s-jwt-to-vault-token` is used as init container for pods which require vault access.


```yaml
      initContainers:
      - args:
        - --vault-address=<vault-address>
        - --vault-role=<vault-consumer>
        - --vault-token-secret-name=<vault-consumer>-vault-token
        - --vault-token-secret-namespace=giantswarm
        image: quay.io/giantswarm/k8s-jwt-to-vault-token:0.1.0
        imagePullPolicy: Always
        name: ensure-vault-token
...
      containers:
      - image: <vault-consumer>
        env:
        - name: VAULT_TOKEN
          valueFrom:
            secretKeyRef:
              key: token
              name: <vault-consumer>-vault-token
```

## How it works?

1. Read Kubernetes service account JWT.
2. Log in vault with JWT and get vault token in response.
3. Write vault token into Kubernetes secret.
4. Consume vault token in main container.