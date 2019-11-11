### How to use

Edit the kubeconfig

```
vi ~/.kube/config
```

Configure the `eks-login` as command: 
```
- name: cluster-name
  user:
    exec:
      command: eks-login
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - -cluster-name
      - <cluster-name>
      - -vault-addr
      - <https://your.vault.domain>
      - -vault-key
      - <aws/creds/k8s-sandbox>
```

**Notes:**
- You will need to be logged in vault
- If the AWS credential is valid eks-login does not be create another one
