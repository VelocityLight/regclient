package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"log"
	"regexp"
)

type ecrAuth struct {
	RoleARN   string
	Reg       string
	Region    string
	Token     *string
	UserName  string
	Password  *string
	ecrClient *ecr.Client
}

func NewEcrAuth(roleARN string, reg string) (Authenticator, error) {
	compileRegex := regexp.MustCompile(`^(\d+)\.dkr\.ecr\.([\w\-]+)\.amazonaws\.com$`)
	matchArr := compileRegex.FindStringSubmatch(reg)
	if len(matchArr) != 3 {
		return nil, fmt.Errorf("the format of ecr url is illegal: %s", reg)
	}

	res := &ecrAuth{
		RoleARN: roleARN,
		Reg:     reg,
		Region:  matchArr[2],
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(res.Region),
	)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	res.ecrClient = ecr.NewFromConfig(cfg)
	//client := sts.NewFromConfig(cfg)

	//if res.RoleARN != "" {
	//	creds := stscreds.NewAssumeRoleProvider(client, res.RoleARN)
	//	res.ecrClient = ecr.NewFromConfig(aws.Config{Credentials: creds, Region: res.Region})
	//}

	token, err := res.ecrClient.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, fmt.Errorf("error get ECR TOKEN: %s", err)
	}

	res.Token = token.AuthorizationData[0].AuthorizationToken

	res.UserName = "AWS"

	pwb := make([]byte, base64.StdEncoding.DecodedLen(len(*res.Token)))
	_, err = base64.StdEncoding.Decode(pwb, []byte(*res.Token))
	if err != nil {
		return nil, fmt.Errorf("error decoding ECR TOKEN: %s", err)
	}

	pw := string(pwb)
	if pw[0:4] == "AWS:" {
		pw = pw[4:]
	} else {
		return nil, fmt.Errorf("error decoding ECR TOKEN: should start with 'AWS:'")
	}
	res.Password = &pw
	return res, nil
}

func (e ecrAuth) CheckRepo(repo string) (string, error) {
	var repoRes *types.Repository
	describeRes, err := e.ecrClient.DescribeRepositories(context.TODO(), &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repo},
	})
	if err != nil {
		createRes, err := e.ecrClient.CreateRepository(context.TODO(), &ecr.CreateRepositoryInput{RepositoryName: &repo})
		if err != nil {
			return "", fmt.Errorf("error create ECR repo when describe failed: %s", err)
		}
		repoRes = createRes.Repository
	} else {
		repoRes = &describeRes.Repositories[0]
	}

	return *repoRes.RepositoryUri, nil
}

func (e ecrAuth) GetToken() string {
	if e.Token == nil {
		return ""
	}
	return *e.Token
}

func (e ecrAuth) GetLoginPassword() (string, string) {
	if e.Password == nil {
		return "", ""
	}
	return e.UserName, *e.Password
}
