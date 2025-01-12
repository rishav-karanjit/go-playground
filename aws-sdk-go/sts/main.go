package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func main() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithSharedConfigProfile("aws-crypto-tools-team+optools-ci-ToolsDevelopment"),
	)
	if err != nil {
		log.Fatal(err)
	}

	stsClient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsClient, "arn:aws:iam::370957321024:role/GitHub-CI-MPL-Dafny-Role-us-west-2", func(o *stscreds.AssumeRoleOptions) {
		o.RoleSessionName = "Go-ESDK-Client-Supplier-Example-Session"
	})
	cfg.Credentials = aws.NewCredentialsCache(provider)

	// without the following, I'm getting an error message: api error SignatureDoesNotMatch: The request signature we calculated does not match the signature you provided.
	_, err = cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
