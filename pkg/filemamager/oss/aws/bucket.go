package aws

import (
	"blog/pkg/l"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func AddBucket(ctx context.Context, cli *s3.S3, userId string) error {
	exist := true
	bucketName := aws.String(userId)
	_, err := cli.HeadBucket(&s3.HeadBucketInput{Bucket: bucketName})
	if err != nil {
		if aerr, ok := err.(awserr.RequestFailure); ok {
			if aerr.StatusCode()/100 == 4 {
				exist = false
			} else {
				l.Logger().Errorf("HeadBucket error, err: %v", err)
				return err
			}
		}
	}

	if !exist {
		out, err := cli.CreateBucket(&s3.CreateBucketInput{Bucket: bucketName})
		if err != nil {
			l.Logger().Errorf("CreateBucket error, err: %v", err)
			return err
		}
		l.Logger().Errorf("AddBucket success, bucket: %v, location: %v", userId, out.Location)
	}
	return nil
}
