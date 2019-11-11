package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

var homeDir = os.Getenv("HOME")
var defaultDir = homeDir + "/" + ".eks-login"

func createFile(file string, content string) {

	if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
		os.Mkdir(defaultDir, 0700)
	}

	f, err := os.Create(defaultDir + "/" + file)
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(content)
}

func getEKSToken(clusterName string) string {
	out, _ := exec.Command("aws", "eks", "get-token", "--cluster-name", clusterName).Output()
	return string(out)
}

func getVaultToken() string {
	token, _ := ioutil.ReadFile(homeDir + "/.vault-token")
	return string(token)
}

func main() {
	vaultAddress := flag.String("vault-address", "", "The vault address, example: https://your.vault.domain")
	clusterName := flag.String("cluster-name", "k8s-sandbox", "EKS cluster name, you can see this name in EKS console")
	vaultPath := flag.String("vault-key", "aws/creds/"+*clusterName, "The vault path, example: aws/creds/clustername.")
	flag.Parse()

	if *vaultAddress == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	_ = godotenv.Load(defaultDir + "/" + *clusterName)

	svc := sts.New(session.New())
	input := &sts.GetCallerIdentityInput{}
	if _, err := svc.GetCallerIdentity(input); err != nil {
		config := &api.Config{
			Address: *vaultAddress,
		}
		client, err := api.NewClient(config)
		if err != nil {
			log.Fatal(err)
		}
		client.SetToken(string(getVaultToken()))
		cl := client.Logical()
		secret, err := cl.Read(*vaultPath)
		if err != nil {
			log.Fatal(err)
		}
		content := fmt.Sprintf("AWS_ACCESS_KEY_ID = %s \n"+
			"AWS_SECRET_ACCESS_KEY = %s \n"+
			"AWS_SESSION_TOKEN = %s",
			secret.Data["access_key"],
			secret.Data["secret_key"],
			secret.Data["security_token"])

		createFile(*clusterName, content)
	}
	fmt.Println(getEKSToken(*clusterName))
}
