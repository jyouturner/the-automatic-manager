package integrationtesting

import (
	"testing"

	"github.com/jyouturner/automaticmanager/pkg/atlanssian"
	automaticmanager "github.com/jyouturner/automaticmanager/tam"
	log "github.com/sirupsen/logrus"
)

func TestJiraSoftwareApiClient_GetBoard(t *testing.T) {
	userCfg, err := automaticmanager.GetUserConfigFromLocalFile("testdata/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	p := atlanssian.NewJiraSoftwareApiClient(userCfg.Atlanssian.JiraSoftware.JiraSoftwareUrl, userCfg.Atlanssian.JiraSoftware.BasicAuthUser, userCfg.Atlanssian.JiraSoftware.BasicAuthToken)
	got, err := p.GetBoard("132")
	if err != nil {
		t.Errorf("JiraSoftwareApiClient.GetBoard() error = %v", err)
		return
	}
	log.Info(got)

}
