package awsserverless

import (
	"encoding/json"
	"log"
	"strings"
)

//StateToken implements the OauthState
type StateToken struct {
	//User has the user's organiation and identifier
	User OrgUser `json:"user"`
	//provider is Google, Atlanssian
	Provider string `json:"provider"`
	//scope is the space delimited string like "scope1 scope2 scope3"
	Scope string `json:"scope"`
	//
	SuccessRedirectUrl string `json:"successRedirectTo"`
}

func StateTokenFromBytes(data []byte) StateToken {
	um := StateToken{}
	err := json.Unmarshal(data, &um)
	if err != nil {
		log.Fatalf("failed to unmarshal  Auth state token, %v", err)
	}
	return um
}

//GetStateData implement the AuthState function
func (p StateToken) GetStateData() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("failed to marshal json of Auth state token, %v", err)
	}
	return b
}

//GetProvider implement the AuthState function
func (p StateToken) GetProvider() string {
	return p.Provider
}

func (p StateToken) GetScope() []string {
	return strings.Split(p.Scope, " ")
}
