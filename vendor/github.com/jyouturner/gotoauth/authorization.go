//This program can handle the 2-leg or 3-leg oauth authorization process, for example, get the oauth2 auth code url and then exchange it with the auth provider
//for access token (and often refresh token too)
package gotoauth

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

//OauthState defines the functions to manage the data saved with the auth state(nounce), it is expected to have the provide data and user data in it.
type OauthState interface {
	GetStateData() []byte
	GetProvider() string
	GetScope() []string
}

type OauthNounceStateWriter interface {
	Save(nounce string, state OauthState) error
}

type OauthNounceStateReader interface {
	Read(nounce string) (OauthState, error)
}

//AppOauthConfig wraps the oatuh secret (for example the oauth2 config JSON) and the token storage specfic to the user
type AppOauthConfig struct {
	Secret            []byte
	OauthTokenStorage TokenStorage
}

//OAuthConfigSource defines the functions that the config provider (for example AWS, Local or database) implement to provide the needed configuration data
type OAuthConfigSource interface {
	//GetAppOathConfig function returns the oauth config data of the given user with the given oauth provider (Google, Atlanssian etc)
	GetAppOathConfig(oauthProvider string) (*AppOauthConfig, error)
}

//GetAuthUrl returns the oauth2 url to get autocode from the provider. It will create a random nounce and use it as key to store the "state" (which wraps the user's identifier, auth provider, whatnot).
func GetAuthUrl(state OauthState, configProvider OAuthConfigSource, nounceWriter OauthNounceStateWriter) (*string, error) {
	log.Debugf("trying to find auth provider %s", state.GetProvider())
	authConfig, err := configProvider.GetAppOathConfig(state.GetProvider())
	if err != nil {
		return nil, fmt.Errorf("failed to get the auth config of provider %s %v", state.GetProvider(), err)
	}

	cfg, err := ConfigFromJSON(authConfig.Secret, state.GetScope())

	if err != nil {
		return nil, fmt.Errorf("could not get the oauth config of the provider %v", err)
	}

	if cfg == nil {
		return nil, fmt.Errorf("no secret found matching the name")
	}

	nounce := String(16)

	//save the nounce for later verification
	err = nounceWriter.Save(nounce, state)
	if err != nil {
		return nil, fmt.Errorf("failed to save nounce %v", err)
	}

	//get the auth URL
	//force approval to support gmail api, and not a bad idea anyway.
	authUrl := cfg.AuthCodeURL(nounce, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	//redirect to the oauth provider authrization URL
	return &authUrl, nil

}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

//Exchange gets the oauth2 access token from the auth code and save it, based on the oauth config.
func Exchange(authcode string, state string, configProvider OAuthConfigSource, nounceReader OauthNounceStateReader) error {

	//find the user data from the nounce(state)
	stateTokenData, err := nounceReader.Read(state)
	if err != nil {
		return fmt.Errorf("failed to locate auth state data %v", err)
	}
	if stateTokenData == nil {
		return fmt.Errorf("no matching user found with the given nounce")
	}
	//state data has user identifier and the auth provider
	//find the oauth config data by the auth provider and user
	authEnv, err := configProvider.GetAppOathConfig(stateTokenData.GetProvider())
	if err != nil {
		return fmt.Errorf("failed to get auth config of provider %v", err)
	}
	//create
	cfg, err := ConfigFromJSON(authEnv.Secret, stateTokenData.GetScope())

	if err != nil {
		return fmt.Errorf("could not get the oauth config of the provider %v", err)
	}

	if cfg == nil {
		return fmt.Errorf("no secret found matching the name")
	}

	token, err := cfg.Exchange(context.TODO(), authcode)
	if err != nil {
		return fmt.Errorf("failed to exchange with oauth provider to get access token from auth code %v", err)
	}

	ts := authEnv.OauthTokenStorage
	if err != nil {
		return fmt.Errorf("error creating token storage %v", err)
	}
	err = ts.SaveNewToken(token)
	if err != nil {
		return fmt.Errorf("failed to save auth token %v", err)
	}
	return nil
}
