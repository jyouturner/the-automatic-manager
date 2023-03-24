//when the app is deployed at AWS, we will use below to save the environment and configurations. This is the SaaS style.
//aws secret manager to store the google and atlanssian secret JSON. In this case, we need the secret name, which can have multiple key-value pairs.
//the oauth token will be saved in S3, under the user's folder, the bucket will be the organization bucket.
//the user config (config.yml) will be saved in the S3 as well, under user folder.
//the S3 bucket structure
// tam-auth-state
// tam-org-12345
//		/abcde
//			/google_token.json
//			/config.yml
package awsserverless

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jyouturner/gotoauth"
	log "github.com/sirupsen/logrus"
)

type awsEnv struct {
	awsClient AWSClient
	//the AWS Secret name. The secret has a bunch of key-value pair, each of which contains the OAuth Config json
	SecretName string
	//the S3 bucket to store the oauth access tokens
	TokenBucketName string
	//the S3 file path to store the oauth access tockets. Note this is the folder for example, if the path is "{user_id}/access_token", then the token is stored in {user_id}/access_token/google_token.json
	TokenFilePath string
	//OauthProviderConfigMap key is the capital case of the oauth provider for example GOOGLE, ATLANSSIAN, tbe value is the oauth config data including the credential secret, and token storage
	oauthProviders map[string]gotoauth.AppOauthConfig
	//the S3 bucket to store the auth nounce (the so called state when initializing the auth code)
	NounceBucket string
}

type UserMeta interface {
	Encode() ([]byte, error)
	GetAccessTokenFolderPath() string
}

//OrgUser implements the UserMeta with organization and user identifier. This is probably enough for most of use cases.
type OrgUser struct {
	OrgId  string `json:"org_id"`
	UserId string `json:"user_id"`
}

func FromJson(data []byte) (OrgUser, error) {
	um := OrgUser{}
	err := json.Unmarshal(data, &um)
	return um, err
}

func (p OrgUser) Encode() ([]byte, error) {
	return json.Marshal(p)
}

func (p OrgUser) GetAccessTokenFolderPath() string {
	return p.UserId
}

func NewAWSEnvByUser(client AWSClient, secretName string, tokenBucketName string, user UserMeta, nounceBucket string) (*awsEnv, error) {

	return NewAWSEnv(client, secretName, tokenBucketName, user.GetAccessTokenFolderPath(), nounceBucket)
}

func NewAWSEnv(client AWSClient, secretName string, tokenBucketName string, tokenFilePath string, nounceBucket string) (*awsEnv, error) {

	p := awsEnv{
		awsClient:       client,
		SecretName:      secretName,
		TokenBucketName: tokenBucketName,
		TokenFilePath:   tokenFilePath,
		oauthProviders:  make(map[string]gotoauth.AppOauthConfig),
		NounceBucket:    nounceBucket,
	}
	err := p.loadOAuthConfigOfUser()
	if err != nil {
		return nil, err
	}
	return &p, nil
}

//GetAppOathConfig implement the OAuthConfigSource
func (p awsEnv) GetAppOathConfig(oauthProvider string) (*gotoauth.AppOauthConfig, error) {
	config, exists := p.oauthProviders[strings.ToUpper(oauthProvider)]
	if !exists {
		return nil, fmt.Errorf("missing config for %s", oauthProvider)
	}
	return &config, nil
}

func (p *awsEnv) loadOAuthConfigOfUser() error {
	// read the aws secret from the secret manager and iterate through the oauth config JSON key-value pairs. This is app specific.
	secretData, err := p.awsClient.GetSecret(p.SecretName)
	if err != nil {
		return err
	}
	if len(*secretData) == 0 {
		return fmt.Errorf("the secret is empty")
	}
	for k, v := range *secretData {

		if strings.HasSuffix(k, "_OAUTH_CONFIG") {
			oauthProvider := substr(k, 0, strings.Index(k, "_"))

			// decide the user's the S3 token bucket and file name, this is user specific
			tokenStorage, err := NewS3TokenStorage(p.awsClient, p.TokenBucketName, fmt.Sprintf("%s/%s_token.json", p.TokenFilePath, strings.ToLower(oauthProvider)))
			if err != nil {
				return fmt.Errorf("failed to create S3 token storage %v", err)
			}
			p.oauthProviders[strings.ToUpper(oauthProvider)] = gotoauth.AppOauthConfig{
				Secret:            []byte(fmt.Sprintf("%v", v)),
				OauthTokenStorage: tokenStorage,
			}

		}
	}
	if len(p.oauthProviders) == 0 {
		return fmt.Errorf("no oauth config is found in the secret that match the expected pattern?")
	}
	return nil
}

func substr(s string, start, end int) string {
	counter, startIdx := 0, 0
	for i := range s {
		if counter == start {
			startIdx = i
		}
		if counter == end {
			return s[startIdx:i]
		}
		counter++
	}
	return s[startIdx:]
}

//SaveStateToken implement the OauthNounceStateWriter
func (p awsEnv) Save(nounce string, data gotoauth.OauthState) error {
	return p.saveNounceToS3(p.awsClient, p.NounceBucket, nounce, data.GetStateData())
}

//saveNounceToS3 stores the nounce data to S3
func (p awsEnv) saveNounceToS3(client AWSClient, bucketName string, nounce string, data []byte) error {
	//
	err := client.S3Save(data, bucketName, nounce)
	if err != nil {
		return fmt.Errorf("failed to save authcode nounce to s3 %s %v", bucketName, err)
	}
	return nil
}

//ReadStateToken implement the OauthNounceStateReader function to load the state data by the given nounce
func (p awsEnv) Read(nounce string) (gotoauth.OauthState, error) {
	b, err := p.awsClient.S3Get(p.NounceBucket, nounce)
	if err != nil {
		return nil, fmt.Errorf("failed to read from S3 err %v", err)
	}
	std := StateToken{}
	log.Println(string(b))
	err = json.Unmarshal(b, &std)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarsh the state token data %v", err)
	}
	return std, nil
}
