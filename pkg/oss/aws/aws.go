package aws

import (
	"blog/pkg/l"
	"blog/pkg/v"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ConnectClient(ctx context.Context) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(v.GetViper().GetString("aws.region")),
		Endpoint: aws.String(v.GetViper().GetString("aws.endPoint")),
		Credentials: credentials.NewStaticCredentials(
			v.GetViper().GetString("aws.accessKeyId"),
			v.GetViper().GetString("aws.secretAccessKey"),
			"",
		),
	})
	if err != nil {
		l.Logger().Errorf("NewSession error, err: %v", err)
		return nil, err
	}
	return sess, err
}
