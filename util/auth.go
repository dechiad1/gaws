package util

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
)

func Auth() *session.Session {
	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		profile = "default"
	}

	fmt.Printf("loading credentials from the %s profile\n", profile)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	}))

	_, err := sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println("can not load credentials.. expecting the ~/.aws directory")
		os.Exit(1)
	}

	return sess
}
