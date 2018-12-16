package util

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func Auth() (*session.Session) {
	profile := os.Getenv("AWS_PROFILE")
	region := os.Getenv("AWS_DEFAULT_REGION")
	if profile == "" {
		profile = "default"
	}

	if region == "" {
		fmt.Println("Please set AWS_DEFAULT_REGION env variable")
		os.Exit(1)
	}

	fmt.Printf("loading credentials from the %s profile\n", profile)
	
	sess, err  := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("",""),
		Region: aws.String(region),
	})

	_, err = sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println("can not load credentials.. check your environment variables or .aws directory")
		os.Exit(1)
	}

	return sess
}
