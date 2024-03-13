package database

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

var MC *minio.Client

func InitMinioClient() *minio.Client {
	minioEndpoint := viper.GetString("minio.endpoint")
	accessKey := viper.GetString("minio.accessID")
	secretKey := viper.GetString("minio.accessKey")
	minioOpt := &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
	}
	mc, err := minio.New(minioEndpoint, minioOpt)
	if err != nil {
		fmt.Println("minio error ", err)
		return nil
	}
	MC = mc
	return mc
}

func GetMinioClient() *minio.Client {
	return MC
}
