package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

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
	fmt.Println(secret.LeaseDuration)
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

func timeNow() int64 {
	t := time.Now()
	return t.Unix()
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

	godotenv.Load(filepath.Join(eksLoginDir, *clusterName))
	if !leaseIsValid() {
		fetchAwsCredsFromVault(*clusterName, *vaultAddr, *vaultPath)
	}
	out, _ := getEKSToken(*clusterName)
	fmt.Println(string(out))
}
