package awshelper

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/project-miko/miko/conf"
	"github.com/project-miko/miko/tools/log"
)

var (
	ExtSvg = ".svg"
	ExtPng = ".png"

	TypeJpg    = "image/jpeg"
	TypePng    = "image/png"
	TypeSvg    = "image/svg+xml"
	TypeMp3    = "audio/mpeg"
	TypeStream = "application/octet-stream"
	cfg        aws.Config
)

type S3CopyObjectAPI interface {
	CopyObject(ctx context.Context,
		params *s3.CopyObjectInput,
		optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error)
}

func InitAwsSDK() {
	var err error
	accessKey := conf.GetConfigString("aws", "aws_access_key_id")
	secretKey := conf.GetConfigString("aws", "aws_secret_access_key")
	region := conf.GetConfigString("aws", "region")

	if len(accessKey) == 0 || len(secretKey) == 0 || len(region) == 0 {
		panic("aws config can not be null")
	}

	cfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			},
		}),
	)

	if err != nil {
		panic(err)
	}
}

func CopyItem(c context.Context, api S3CopyObjectAPI, input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	return api.CopyObject(c, input)
}

func CopyObject(kFrom string, kTo string, bucket string) error {
	client := s3.NewFromConfig(cfg)
	input := &s3.CopyObjectInput{
		Bucket: aws.String(bucket),

		CopySource: aws.String(bucket + "/" + kFrom),
		Key:        aws.String(kTo),
	}

	_, err := CopyItem(context.TODO(), client, input)
	if err != nil {

		return err
	}

	return nil
}

func UploadPublicAsset(bucket string, typePtr *string, res io.Reader, subpath string, fileKey string) (location string, err error) {

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)

	// ex:
	// subpath := fmt.Sprintf("%s/%s%s", uploadDir, fileKey, extName)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(subpath),
		Body:        res,
		ContentType: typePtr,
	})
	if err != nil {
		return "", err
	}

	log.Debug("", "upload done %s %s %s", result.Location, result.UploadID, *(result.VersionID))
	return result.Location, nil
}

func Upload(bucket string, key string, typePtr *string, res io.Reader) (location string, err error) {
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        res,
		ContentType: typePtr,
	})
	if err != nil {
		return "", err
	}

	log.Debug("", "upload done %s %s", result.Location, result.UploadID)
	return result.Location, nil
}
