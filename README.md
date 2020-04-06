### Description

I would like do use the [hashicorp vault](https://www.vaultproject.io/docs/secrets/aws/index.html) to get temporary [AWS Credencial](https://www.vaultproject.io/docs/secrets/aws/index.html) and using it to access the EKS service.

The problem is, the STS AWS credentials no valid for more than 12 hours, so, this script will automate the process.

**Auth methods supported:**

- github

### How to use

1 - Download the binary from github page or on OsX:

```
brew install luizm/tap/eks-login
```

In the correct context into kubeconfig file, edit the `command` block and use `eks-login` instead of `aws cli` or `aws-iam-authenticator`

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
```
