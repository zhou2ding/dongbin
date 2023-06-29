package minio

import (
	"blog/pkg/v"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/pkg/bucket/policy"
	"github.com/minio/pkg/bucket/policy/condition"
)

func ConnectClient() (*minio.Client, error) {
	MinioClient, err := minio.New(fmt.Sprintf("%s:%d", v.GetViper().GetString("minio.url"), v.GetViper().GetInt("minio.port")), &minio.Options{
		Creds:  credentials.NewStaticV4(v.GetViper().GetString("minio.secretKey"), v.GetViper().GetString("minio.accessKey"), ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	} else {
		return MinioClient, nil
	}
}

func AddBucket(ctx context.Context, client *minio.Client, name string) error {
	exist, err := client.BucketExists(ctx, name)
	if err != nil {
		return err
	}
	if !exist {
		opts := minio.MakeBucketOptions{}
		err = client.MakeBucket(ctx, name, opts)
		if err != nil {
			return err
		}
		po, _ := policy.Policy{
			Version: policy.DefaultVersion,
			Statements: []policy.Statement{
				policy.NewStatement("",
					policy.Allow, policy.NewPrincipal("*"),
					policy.NewActionSet(policy.GetObjectAction, policy.ListMultipartUploadPartsAction, policy.PutObjectAction, policy.AbortMultipartUploadAction),
					policy.NewResourceSet(policy.NewResource(name, "*")),
					condition.NewFunctions(),
				),
			},
		}.MarshalJSON()
		err = client.SetBucketPolicy(ctx, name, string(po))
		if err != nil {
			return err
		}
	}
	return nil
}
