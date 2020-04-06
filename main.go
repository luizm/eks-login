package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/luizm/eks-login/internal/eks"
	"github.com/luizm/eks-login/internal/vault"
)

const version string = "v1.0.1"

var homeDir = os.Getenv("HOME")
var eksLoginDir = filepath.Join(os.Getenv("HOME"), ".eks-login")

func main() {
	clusterName := flag.String("cluster-name", "k8s-sandbox", "EKS cluster name, you can see this name in EKS console")
	region := flag.String("region", "us-east-1", "AWS region where EKS cluster is running")
	vaultAddr := flag.String("vault-addr", "", "The vault address, example: https://your.vault.domain")
	vaultPath := flag.String("vault-path", "aws/creds/"+*clusterName, "The vault endpoint path, example: aws/creds/clustername")
	githubTokenPath := flag.String("github-token-path", homeDir+"/.github-token", "Path to file with github credential")
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

	godotenv.Load(filepath.Join(eksLoginDir, *clusterName))
	if !vault.LeaseIsValid() {
		if content, err := vault.FetchAwsCredsFromVault(*clusterName, *vaultAddr, *vaultPath, *githubTokenPath); err != nil {
			log.Fatalln(err)
		} else {
			if err := os.MkdirAll(eksLoginDir, 0700); err != nil {
				log.Fatalln(err)
			}
			if err := ioutil.WriteFile(filepath.Join(eksLoginDir, *clusterName), []byte(content), 0644); err != nil {
				log.Fatalln(err)
			}
			godotenv.Load(filepath.Join(eksLoginDir, *clusterName))
		}
	}
	out, err := eks.GetEKSToken(*clusterName, *region)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(out)
}
