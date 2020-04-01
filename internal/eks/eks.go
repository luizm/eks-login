package eks

import (
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

//GetEKSToken get the temporary token in json format
func GetEKSToken(clusterName, region string) (string, error) {
	var tok token.Token
	var out string
	var err error

	gen, err := token.NewGenerator(false, false)
	if err != nil {
		return "", err
	}

	tok, err = gen.GetWithOptions(&token.GetTokenOptions{
		ClusterID:            clusterName,
		Region:               region,
		AssumeRoleARN:        "",
		AssumeRoleExternalID: "",
		SessionName:          "",
	})

	if err != nil {
		return "", err
	}

	out = gen.FormatJSON(tok)

	return out, nil
}
