package platform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"secure-banking-uk-initializer/pkg/common"
	"secure-banking-uk-initializer/pkg/httprest"
	"secure-banking-uk-initializer/pkg/types"

	"go.uber.org/zap"
)

func CreateIGServiceUser() {
	if httprest.ServiceIdentityExists(common.Config.Ig.IgIdmUser) {
		zap.L().Info("Skipping creation of IG service user")
		return
	}

	zap.L().Info("Creating IG service user")

	user := &types.ServiceUser{
		UserName:  common.Config.Ig.IgIdmUser,
		SN:        "Service Account",
		GivenName: "IG",
		Mail:      "ig@acme.com",
		Password:  common.Config.Ig.IgIdmPassword,
		AuthzRole: []types.AuthzRole{
			{
				Ref: "internal/role/openidm-admin",
			},
		},
	}
	path := "/openidm/managed/user/?_action=create"
	//FIDC IDM default user managed objects use a different naming pattern <realm>_user Eg:alpha_user
	if common.Config.Environment.CloudType == "FIDC" {
		path = "/openidm/managed/" + common.Config.Identity.AmRealm + "_user/?_action=create"
	}
	_, s := httprest.Client.Post(path, user, map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/json",
		"Connection":   "keep-alive",
	})

	zap.S().Infow("IG Service User", "statusCode", s)
}

// CreateIGOAuth2Client -
func CreateIGOAuth2Client() {
	if httprest.OAuth2AgentClientsExist(common.Config.Ig.IgClientId) {
		zap.S().Infof("Skipping creation of IG Oauth2 client. OAuth2 client %s already exists", common.Config.Ig.IgClientId)
		return
	}

	zap.S().Infof("Creating IG OAuth2 client with id %s", common.Config.Ig.IgClientId)
	oauth2Client := &types.OAuth2Client{}

	err := common.Unmarshal(common.Config.Environment.Paths.ConfigIdentityPlatform+"ig-oauth2-client.json", &common.Config, oauth2Client)
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("/am/json/"+common.Config.Identity.AmRealm+"/realm-config/agents/OAuth2Client/%s", common.Config.Ig.IgClientId)
	s := httprest.Client.Put(path, oauth2Client, map[string]string{
		"Accept":           "application/json",
		"Content-Type":     "application/json",
		"Connection":       "keep-alive",
		"X-Requested-With": "securebanking-openbanking-uk-fidc-initializer",
	})

	zap.S().Infow("IG OAuth2 Client", "statusCode", s)
}

// CreateIGPolicyAgent -
func CreateIGPolicyAgent() {
	zap.L().Info("Creating IG Policy agent")
	policyAgent := &types.PolicyAgent{
		Userpassword: common.Config.Ig.IgAgentPassword,
		IgTokenIntrospection: types.IgTokenIntrospection{
			Value:     "Realm",
			Inherited: false,
		},
	}
	path := fmt.Sprintf("/am/json/realms/root/realms/"+common.Config.Identity.AmRealm+"/realm-config/agents/IdentityGatewayAgent/%s", common.Config.Ig.IgAgentId)
	s := httprest.Client.Put(path, policyAgent, map[string]string{
		"Accept":           "application/json",
		"Content-Type":     "application/json",
		"Connection":       "keep-alive",
		"X-Requested-With": "ForgeRock Identity Cloud Postman Collection",
	})

	zap.S().Infow("IG Policy Agent", "statusCode", s)
}

func CreateIdentityPlatformOAuth2AdminClient(cookie *http.Cookie) {
	zap.L().Info("Creating Identity Platform admin oauth2 client")

	oauth2Client := &types.OAuth2Client{}
	e := common.Unmarshal(common.Config.Environment.Paths.ConfigIdentityPlatform+"oauth2-admin-client.json", &common.Config, oauth2Client)
	if e != nil {
		panic(e)
	}
	zap.S().Debugw("Admin client request", "body", oauth2Client)
	path := "https://" + common.Config.Hosts.IdentityPlatformFQDN + "/am/json/realm-config/agents/OAuth2Client/idmAdminClient"
	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("Connection", "keep-alive").
		SetHeader("X-Requested-With", "ForgeRock Identity Cloud Postman Collection").
		SetContentLength(true).
		SetCookie(cookie).
		SetBody(oauth2Client).
		Put(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	zap.S().Infow("IDM Admin Client", "statusCode", resp.StatusCode(), "redirect", oauth2Client.CoreOAuth2ClientConfig.RedirectionUris.Value)
}

// CreateRealm creates the realm for a new deployment of CDK
func CreateRealm(cookie *http.Cookie) {
	zap.S().Infof("Creating '%s' Realm", common.Config.Identity.AmRealm)
	b, err := common.Template(common.Config.Environment.Paths.ConfigIdentityPlatform+"realm-template.json", &common.Config)
	if err != nil {
		panic(err)
	}
	path := "https://" + common.Config.Hosts.IdentityPlatformFQDN + "/am/json/global-config/realms?_action=create"
	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Requested-With", "ForgeRock Identity Cloud Postman Collection").
		SetContentLength(true).
		SetCookie(cookie).
		SetBody(b).
		Post(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	zap.S().Infow(common.Config.Identity.AmRealm+" Realm Created", "statusCode", resp.StatusCode())
}

// RealmExist check if the realm exist
func RealmExist(cookie *http.Cookie) bool {
	zap.S().Infof("Checking if '%s' Realm already exist", common.Config.Identity.AmRealm)
	var realmExist = false
	path := fmt.Sprintf("https://%s/am/json/global-config/realms?_queryFilter=true", common.Config.Hosts.IdentityPlatformFQDN)
	serviceIdentityFilter := &types.ResultFilter{}
	resp, errResp := restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-With", "ForgeRock Identity Cloud Postman Collection").
		SetHeader("Accept-Api-Version", "protocol=2.0,resource=1.0").
		SetCookie(cookie).
		Get(path)
	if errResp != nil {
		panic(errResp)
	}
	err := json.Unmarshal(resp.Body(), serviceIdentityFilter)
	if err != nil {
		panic(err)
	}

	for _, s := range serviceIdentityFilter.Result {
		if s.Name == common.Config.Identity.AmRealm {
			realmExist = true
		}
	}
	zap.S().Infow("Check realm exist", "realm", common.Config.Identity.AmRealm, "exist", realmExist)
	return realmExist
}

func CreateServerConfig(cookie *http.Cookie) {
	zap.L().Info("Pushing Creating ServerDefault - Advanced Settings")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigIdentityPlatform + "server-defaults.json")
	if err != nil {
		panic(err)
	}
	path := "https://" + common.Config.Hosts.IdentityPlatformFQDN + "/am/json/global-config/servers/server-default/properties/advanced"
	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Requested-With", "ForgeRock Identity Cloud Postman Collection").
		SetHeader("accept-api-version", "protocol=1.0,resource=1.0").
		SetContentLength(true).
		SetCookie(cookie).
		SetBody(b).
		Put(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	zap.S().Infow("Pushed server default - Advanced Settings")
}
