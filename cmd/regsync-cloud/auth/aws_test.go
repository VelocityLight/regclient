package auth

import (
	"context"
	"log"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func TestDescribeRepo(t *testing.T) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	ecrClient := ecr.NewFromConfig(cfg)

	describeRes, err := ecrClient.DescribeRepositories(context.TODO(), &ecr.DescribeRepositoriesInput{
		RegistryId:      &[]string{"TARGET_REG"}[0],
		RepositoryNames: []string{"TARGET_REPO"},
	})

	if err != nil && strings.Contains(err.Error(), "does not exist in the registry") {
		log.Fatalf("failed to describe repository, %v", err)
	}

	log.Println(describeRes)
}
