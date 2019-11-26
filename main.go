package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

var homeDir = os.Getenv("HOME")
var eksLoginDir = filepath.Join(homeDir, ".eks-login")

const version string = "v0.1.1"

func createFile(clusterName string, content string) error {
	if err := os.MkdirAll(eksLoginDir, 0700); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(eksLoginDir, clusterName), []byte(content), 0644); err != nil {
		return err
	}
	return nil
}

func getEKSToken(clusterName string) ([]byte, error) {
	loadEnv(clusterName)
	return exec.Command("aws", "eks", "get-token", "--cluster-name", clusterName).Output()
}

func getGithubToken(githubTokenPath string) ([]byte, error) {
	return ioutil.ReadFile(githubTokenPath)
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

func fetchAwsCredsFromVault(clusterName, vaultAddr, vaultPath, githubTokenPath string) error {
	config := &api.Config{
		Address: vaultAddr,
	}
	vaultToken, err := getVaultTokenGitHub(vaultAddr, githubTokenPath)
	if err != nil {
		return err
	}
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	client.SetToken(string(vaultToken))
	cl := client.Logical()
	secret, err := cl.Read(vaultPath)
	if err != nil {
		return err
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
	return createFile(clusterName, content)
}

func leaseIsValid() bool {
	creationTime, _ := strconv.ParseInt(os.Getenv("CREATION_TIME"), 10, 64)
	ttl, _ := strconv.ParseInt(os.Getenv("TTL"), 10, 64)
	if timeNow() > (creationTime + ttl) {
		return false
	}
	return true
}

func loadEnv(clusterName string) {
	godotenv.Load(filepath.Join(eksLoginDir, clusterName))
}

func timeNow() int64 {
	t := time.Now()
	return t.Unix()
}

func main() {
	clusterName := flag.String("cluster-name", "k8s-sandbox", "EKS cluster name, you can see this name in EKS console")
	vaultAddr := flag.String("vault-addr", "", "The vault address, example: https://your.vault.domain")
	vaultPath := flag.String("vault-path", "aws/creds/"+*clusterName, "The vault path, example: aws/creds/clustername.")
	githubTokenPath := flag.String("github-token-path", homeDir+"/.github-token", "Path to get the github credential")

	appVersion := flag.Bool("version", false, "Shows application version")
	flag.Parse()

	if *appVersion == true {
		out := fmt.Sprintf("%s %s", "eks-login", version)
		fmt.Println(out)
		os.Exit(0)
	}
	if *vaultAddr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	loadEnv(*clusterName)
	if !leaseIsValid() {
		if err := fetchAwsCredsFromVault(*clusterName, *vaultAddr, *vaultPath, *githubTokenPath); err != nil {
			log.Fatalln(err)
		}
	}
	out, _ := getEKSToken(*clusterName)
	fmt.Println(string(out))
}
