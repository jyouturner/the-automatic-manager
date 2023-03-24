//This program has the client code to handle AWS services including saving to and reading from S3, reading from Sercret Manager
package awsserverless

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

//AWSClient wraps the aws config and receiver of multiple AWS client functions
type AWSClient struct {
	Config aws.Config
}

//S3Get can fetch S3 object
func (p AWSClient) S3Get(bucket string, key string) ([]byte, error) {

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(p.Config)
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	output, err := client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3, %v", err)
	}
	defer output.Body.Close()
	return ioutil.ReadAll(output.Body)

}

//S3Save uploads data to S3
func (p AWSClient) S3Save(data []byte, bucket string, key string) error {

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(p.Config)
	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	}
	_, err := client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload to S3, %v", err)
	}

	return nil
}

//S3Ls list S3 buckets
func (p AWSClient) S3Ls() (*s3.ListBucketsOutput, error) {

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(p.Config)

	output, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3, %v", err)
	}

	for _, bucket := range output.Buckets {
		log.Info(*bucket.Name)
	}
	return output, nil
}

//GetSecret read the secret from AWS secret manager
func (p AWSClient) GetSecret(name string) (*map[string]string, error) {

	var secretData = make(map[string]string)
	client := secretsmanager.NewFromConfig(p.Config)

	result, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: &name,
	})

	if err != nil {
		log.Error("failed to read secret ", err)
		return nil, err
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	if result.SecretString != nil {
		secretString := *result.SecretString
		err := json.Unmarshal([]byte(secretString), &secretData)
		if err != nil {
			return nil, fmt.Errorf("could not unmarsh the secrete data to JSON %v", err)
		}
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			return nil, fmt.Errorf("could not base64 decode the secret data %v", err)
		}
		json.Unmarshal(decodedBinarySecretBytes[:len], &secretData)

	}

	return &secretData, nil
}

//SecretToEnvVariables load the AWS secret, then set to environment variables.
func (p AWSClient) SecretToEnvVariables(secretName string) error {
	secretData, err := p.GetSecret(secretName)
	if err != nil {
		return fmt.Errorf("failed to get secret %s %v", secretName, err)
	}

	for key, value := range *secretData {
		os.Setenv(key, fmt.Sprintf("%v", value))
	}
	return nil
}

//GetSecretValueFromKey will fetch the aws secret, then try to find the value of the matching key, return the []byte of value
func (p AWSClient) GetSecretValueFromKey(secretName string, key string) ([]byte, error) {
	secretData, err := p.GetSecret(secretName)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s %v", secretName, err)
	}
	for k, v := range *secretData {
		if k == key {
			return []byte(fmt.Sprintf("%v", v)), nil
		}

	}
	return nil, nil
}
