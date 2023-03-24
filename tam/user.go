package automaticmanager

import (
	"fmt"
	"io/ioutil"

	"github.com/jyouturner/gotoauth/example/awsserverless"
	"gopkg.in/yaml.v3"
)

const S3_BUCKET_PREFIX = "tam-org"

type User struct {
	//we use the "OrgUser" from the gotoauth package. It is a simple structure with OrgId and UserId, used in the oauth authorization flow.
	Id awsserverless.OrgUser
	//the S3 config of the user, where to find the user's configuration yml file.
	BucketName string
	ConfigFile string
}

//NewUser create a user object with the given user id and orgnization id, these two will decdie the location of the config.yml file in S3
func NewUser(userId string, orgId string) User {
	return NewUserWithOrg(awsserverless.OrgUser{
		OrgId:  orgId,
		UserId: userId,
	})
}

func NewUserWithOrg(orgUser awsserverless.OrgUser) User {
	return User{
		Id:         orgUser,
		BucketName: getBucketOfUser(orgUser.UserId, orgUser.OrgId),
		ConfigFile: fmt.Sprintf("%s/config.yml", orgUser.UserId),
	}
}

func getBucketOfUser(userId string, orgId string) string {
	return fmt.Sprintf("%s-%s", S3_BUCKET_PREFIX, orgId)
}

//The config.yml file
type UserConfig struct {
	Calendar struct {
		NumberOfEvents    int    `yaml:"number_of_events"`
		ExcludeEmailsList string `yaml:"exclude_emails_list"`
	} `yaml:"calendar"`
	Atlanssian struct {
		JiraClound struct {
			JiraCloudId   string `yaml:"jira_cloud_id"`
			JiraCloundUrl string `yaml:"jira_cloud_url"`
			JiraMonitor   struct {
				JQL string `yaml:"jql"`
			} `yaml:"jira_monitor"`
		} `yaml:"jira_cloud"`
		JiraSoftware struct {
			BasicAuthUser   string `yaml:"basic_auth_user"`
			BasicAuthToken  string `yaml:"basic_auth_token"`
			JiraSoftwareUrl string `yaml:"jira_software_url"`
		} `yaml:"jira_software"`
		Confluence struct {
			BasicAuthUser  string `yaml:"basic_auth_user"`
			BasicAuthToken string `yaml:"basic_auth_token"`
			ConfluenceUrl  string `yaml:"confluence_software_url"`
		} `yaml:"confluence"`
	} `yaml:"atlanssian"`
	Notion struct {
		ApiKey         string `yaml:"api_key"`
		TaskDatabaseId string `yaml:"task_database_id"`
	} `yaml:"notion"`
	Github struct {
		AccessToken string `yaml:"access_token"`
	} `yaml:"github"`
}

//GetUserConfig will read the YML file and return the UserEnv
func GetUserConfigFromS3(awsClient awsserverless.AWSClient, user User) (*UserConfig, error) {
	//read the s3 file
	b, err := awsClient.S3Get(user.BucketName, user.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read yml file %v", err)
	}
	var cfg = UserConfig{}
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall the yml file %v", err)
	}
	return &cfg, nil
}

func GetUserConfigFromLocalFile(filepath string) (*UserConfig, error) {
	//read the s3 file
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read yml file %v", err)
	}
	var cfg = UserConfig{}
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall the yml file %v", err)
	}
	return &cfg, nil
}
