package local

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jyouturner/gotoauth"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

//LocalTokenStorage provides the methods to persist the token in local file. This can be used for simple oauth2 testing.
type LocalTokenStorage struct {
	TokenFile string
}

// LoadToken retrieves a token from a local file.
func (p LocalTokenStorage) LoadToken() (*oauth2.Token, error) {
	f, err := os.Open(p.TokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// SaveNewToken saves a token to a file path.
func (p LocalTokenStorage) SaveNewToken(token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", p.TokenFile)
	f, err := os.OpenFile(p.TokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Errorf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)

}

//ConfigFromLocalJSONFile loads oauth2 config from local json file. This is a helper function.
func ConfigFromLocalJsonFile(secretFile string, scope []string) (*oauth2.Config, error) {
	//read secret file into oauth config
	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Errorf("Unable to read client secret file: %v", err)
		return nil, err
	}
	return gotoauth.ConfigFromJSON(b, scope)
}
