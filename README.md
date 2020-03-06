### Description

I would like do use the vault to get temporary [AWS credencial](https://www.vaultproject.io/docs/secrets/aws/index.html) and using it to access the EKS service.

The problem is, the STS AWS credentials no valid for more than 12 hours, so, this script will automate the process.

**Notes:**

- The `aws cli` is necessary yet
- The github auth is the only method supported to auth into vault

### How to use

1 - Download the binary from github page or on OsX:

```
brew install
```

In the correct context into kubeconfig file, edit de command block and set the `eks-login` including some arguments:

Example:

```
- name: cluster-name
  user:
    exec:
      command: eks-login
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - -cluster-name
      - <CLUSTER_NAME>
      - -vault-addr
      - <https://VAULT_ENDPOINT>
      - -vault-path
      - <PATH_TO_GET_THE_CREDENDIALS>
      - -github-token-path
      - <GITHUB_TOKEN_PATH>
```
