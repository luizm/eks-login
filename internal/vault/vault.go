package vault

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

func getGithubToken(githubTokenPath string) ([]byte, error) {
	return ioutil.ReadFile(githubTokenPath)
}

func timeNow() int64 {
	t := time.Now()
	return t.Unix()
}

func getVaultTokenGitHub(vaultAddr, githubTokenPath string) (string, error) {
	gitHubToken, _ := getGithubToken(githubTokenPath)
	options := map[string]interface{}{
		"token": strings.TrimSpace(string(gitHubToken)),
	}
	config := &api.Config{
		Address: vaultAddr,
	}
	client, err := api.NewClient(config)
	if err != nil {
		return "", err
	}
	secret, err := client.Logical().Write("/auth/github/login", options)
	if err != nil {
		return "", err
	}
	return secret.Auth.ClientToken, nil
}

// LeaseIsValid check if credentials still valid
func LeaseIsValid() bool {
	creationTime, _ := strconv.ParseInt(os.Getenv("CREATION_TIME"), 10, 64)
	ttl, _ := strconv.ParseInt(os.Getenv("TTL"), 10, 64)
	if timeNow() > (creationTime + ttl) {
		return false
	}
	return true
}

//FetchAwsCredsFromVault get the aws credencials from vault
func FetchAwsCredsFromVault(clusterName, vaultAddr, vaultPath, githubTokenPath string) (string, error) {
	config := &api.Config{
		Address: vaultAddr,
	}
	vaultToken, err := getVaultTokenGitHub(vaultAddr, githubTokenPath)
	if err != nil {
		return "", err
	}
	client, err := api.NewClient(config)
	if err != nil {
		return "", err
	}
	client.SetToken(string(vaultToken))
	cl := client.Logical()
	secret, err := cl.Read(vaultPath)
	if err != nil {
		return "", err
	}
	content := fmt.Sprintf("AWS_ACCESS_KEY_ID = %s \n"+
		"AWS_SECRET_ACCESS_KEY = %s \n"+
		"AWS_SESSION_TOKEN = %s \n"+
		"CREATION_TIME = %d \n"+
		"TTL = %d \n",
		secret.Data["access_key"],
		secret.Data["secret_key"],
		secret.Data["security_token"],
		timeNow(),
		secret.LeaseDuration)
	return content, nil
}
