package eks

import (
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

//GetEKSToken return the temporary token in json format
func GetEKSToken(clusterName, region string) (string, error) {
	gen, err := token.NewGenerator(false, false)
	if err != nil {
		return "", err
	}

	tok, err := gen.GetWithOptions(&token.GetTokenOptions{
		ClusterID:            clusterName,
		Region:               region,
		AssumeRoleARN:        "",
		AssumeRoleExternalID: "",
		SessionName:          "",
	})
	if err != nil {
		return "", err
	}

	return gen.FormatJSON(tok), nil
}
