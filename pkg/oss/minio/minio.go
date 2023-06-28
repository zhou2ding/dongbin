package minio

import (
	"blog/pkg/v"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectClient() (*minio.Client, error) {
	if MinioClient, e0 := minio.New(fmt.Sprintf("%s:%d", v.GetViper().GetString("minio.url"), v.GetViper().GetInt("minio.port")), &minio.Options{
		Creds:  credentials.NewStaticV4(v.GetViper().GetString("minio.secretKey"), v.GetViper().GetString("minio.accessKey"), ""),
		Secure: false,
	}); e0 != nil {
		return nil, e0
	} else {
		return MinioClient, nil
	}
}
