package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

var homeDir = os.Getenv("HOME")
var eksLoginDir = filepath.Join(homeDir, ".eks-login")

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
	return exec.Command("aws", "eks", "get-token", "--cluster-name", clusterName).Output()
}

func getVaultToken() ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(homeDir, ".vault-token"))
}

func fetchAwsCredsFromVault(clusterName, vaultAddr, vaultPath string) error {
	config := &api.Config{
		Address: vaultAddr,
	}
	vaultToken, err := getVaultToken()
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
		"AWS_SESSION_TOKEN = %s",
		secret.Data["access_key"],
		secret.Data["secret_key"],
		secret.Data["security_token"])
	return createFile(clusterName, content)
}

func canAuthenticateToAws(clusterName string) bool {
	godotenv.Load(filepath.Join(eksLoginDir, clusterName))
	svc := sts.New(session.New())
	input := &sts.GetCallerIdentityInput{}
	if _, err := svc.GetCallerIdentity(input); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func main() {
	clusterName := flag.String("cluster-name", "k8s-sandbox", "EKS cluster name, you can see this name in EKS console")
	vaultAddr := flag.String("vault-addr", "", "The vault address, example: https://your.vault.domain")
	vaultPath := flag.String("vault-path", "aws/creds/"+*clusterName, "The vault path, example: aws/creds/clustername.")
	flag.Parse()

	if *vaultAddr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if !canAuthenticateToAws(*clusterName) {
		fetchAwsCredsFromVault(*clusterName, *vaultAddr, *vaultPath)
	}
	out, _ := getEKSToken(*clusterName)
	fmt.Println(string(out))
}
