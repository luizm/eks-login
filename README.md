### Description

I would like do use the vault to get temporary [AWS credencial](https://www.vaultproject.io/docs/secrets/aws/index.html) and access the EKS service.

The problem is, the AWS credentials no valid for more than 12 hours, so, this script will automate this process.

**Notes:**

- The `aws cli` is necessary yet
- The github auth is the only method supported to auth into vault
- If the AWS credential is valid eks-login does not be create another one

<b>

### How to install 

Using homebrew:

```sh
brew install luizm/tap/eks-login
```

For Linux platform, download the binary directly on the [release page](https://github.com/luizm/eks-login/releases) 

<b>

### How to use

Edit the kubeconfig

```
vi ~/.kube/config
```

In the right context, configure the `eks-login` as command instead of `aws cli` or `aws-iam-authenticator`

Example:

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
      - -vault-path
      - <aws/creds/k8s-sandbox>
      - -github-token-path
      - ~/.github-token
```
