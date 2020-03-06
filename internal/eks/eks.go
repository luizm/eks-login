package eks

import (
	"os/exec"
)

//GetEKSToken get the EKS token using aws cli
func GetEKSToken(clusterName string) ([]byte, error) {
	return exec.Command("aws", "eks", "get-token", "--cluster-name", clusterName).Output()
}
