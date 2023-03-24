package awsserverless

import (
	"encoding/json"
	"fmt"

	"github.com/jyouturner/gotoauth"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

//S3TokenStorage provides the methods to persist the token in AWS S3 bucket.
type s3TokenStorage struct {
	Bucket string
	Key    string
	client AWSClient
}

// NewS3TokenStorage return a S3TokenStorage to store the oauth2 token with the provided bucket and key.
func NewS3TokenStorage(client AWSClient, bucket string, key string) (gotoauth.TokenStorage, error) {

	return s3TokenStorage{
		Bucket: bucket,
		Key:    key,
		client: client,
	}, nil
}

// LoadToken retrieves a token from a S3 file, mplement the LoadToken of TokenStorage
func (p s3TokenStorage) LoadToken() (*oauth2.Token, error) {

	b, err := p.client.S3Get(p.Bucket, p.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get JSON from S3, %v", err)
	}
	log.Debugf("fetch token %s", string(b))
	tok := &oauth2.Token{}

	err = json.Unmarshal(b, &tok)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from S3, %v", err)
	}
	return tok, nil
}

//SaveNewToken save the token to AWS S3, implement the SaveNewToken of TokenStorage
func (p s3TokenStorage) SaveNewToken(token *oauth2.Token) error {
	body, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data, %v", err)
	}

	return p.client.S3Save(body, p.Bucket, p.Key)
}
