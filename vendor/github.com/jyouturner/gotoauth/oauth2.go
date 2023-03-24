package gotoauth

//
// reference https://github.com/golang/oauth2/issues/84
//
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// TokenNotifyFunc is a function that accepts an oauth2 Token upon refresh, and
// returns an error if it should not be used.
type TokenNotifyFunc func(*oauth2.Token) error

// NotifyRefreshTokenSource is essentially `oauth2.ResuseTokenSource` with `TokenNotifyFunc` added.
type NotifyRefreshTokenSource struct {
	new oauth2.TokenSource
	mu  sync.Mutex // guards t
	t   *oauth2.Token
	f   TokenNotifyFunc // called when token refreshed so new refresh token can be persisted
}

//TokenStorage defines the interface for the token storages for example AWS S3, AWS Secret Manager or K8S secret
type TokenStorage interface {
	SaveNewToken(t *oauth2.Token) error
	LoadToken() (*oauth2.Token, error)
}

// Token returns the current token if it's still valid, else will
// refresh the current token (using r.Context for HTTP client
// information) and return the new one.
func (s *NotifyRefreshTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		log.Info("returning existing token")
		return s.t, nil
	}
	t, err := s.new.Token()
	if err != nil {
		log.Errorf("failed to get token %v", err)
		return nil, err
	}
	s.t = t
	return t, s.f(t)
}

// NewClient will return a http client with the tokens from the token file. The client can refresh the token automatically, since it is aware of
// the refresh_token in the Token Source. You should only need to call this function once in the flow and reuse the client in furture requests.
func NewClient(tokenStorage TokenStorage, config *oauth2.Config) (*http.Client, error) {

	ctx := context.Background()
	tok, err := tokenStorage.LoadToken()

	if err != nil {
		log.Error("token file does not exist, please genereate it first")
		return nil, err
	}

	if !tok.Valid() {
		log.Info("Access token is either not existing or expired")
	}

	nrts := &NotifyRefreshTokenSource{
		new: config.TokenSource(ctx, tok),
		t:   tok,
		f:   tokenStorage.SaveNewToken,
	}

	return oauth2.NewClient(ctx, nrts), nil
}

/*
func refreshToken() {
	tokenSource := conf.TokenSource(oauth2.NoContext, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalln(err)
	}

	if newToken.AccessToken != token.AccessToken {
		SaveToken(newToken)
		log.Println("Saved new token:", newToken.AccessToken)
	}

	client := oauth2.NewClient(oauth2.NoContext, tokenSource)
	resp, err := client.Get(url)
}
*/

//ConfigFromJSON will create the oauth2 config from JSON, with specified scope
func ConfigFromJSON(jsonKey []byte, scope []string) (*oauth2.Config, error) {
	type cred struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURI  string `json:"redirect_uri"`
		AuthURI      string `json:"auth_uri"`
		TokenURI     string `json:"token_uri"`
	}
	var j struct {
		Web       *cred `json:"web"`
		Installed *cred `json:"installed"`
	}
	if err := json.Unmarshal(jsonKey, &j); err != nil {
		return nil, err
	}
	var c *cred
	switch {
	case j.Web != nil:
		c = j.Web
	case j.Installed != nil:
		c = j.Installed
	default:
		return nil, fmt.Errorf("no credentials found")
	}

	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURI,
		Scopes:       scope,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AuthURI,
			TokenURL: c.TokenURI,
		},
	}, nil
}
