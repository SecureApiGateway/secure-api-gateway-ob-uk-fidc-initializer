package securebanking

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

func CreateSecureBankingRemoteConsentService() {
	remoteConsentId := common.Config.Identity.RemoteConsentId
	if remoteConsentExists(remoteConsentId) {
		zap.L().Info("Remote consent exists. skipping")
		return
	}
	zap.L().Info("Creating remote consent service")
	signingPublicKey := common.Config.Identity.RemoteConsentSigningPublicKey
	if signingPublicKey == "" {
		zap.S().Fatal("RemoteConsentSigningPublicKey must be configured")
	}
	signingKeyId := common.Config.Identity.RemoteConsentSigningKeyId
	if signingKeyId == "" {
		zap.S().Fatal("RemoteConsentSigningKeyId must be configured")
	}
	rcsJwks := CreateRcsJwks(signingPublicKey, signingKeyId)

	rc := &types.RemoteConsent{
		RemoteConsentRequestEncryptionAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "RSA-OAEP-256",
		},
		PublicKeyLocation: types.InheritedValueString{
			Inherited: false,
			Value:     "jwks",
		},
		JwksCacheTimeout: types.InheritedValueInt{
			Inherited: false,
			Value:     0,
		},
		RemoteConsentResponseSigningAlg: types.InheritedValueString{
			Inherited: false,
			Value:     "PS256",
		},
		RemoteConsentRequestSigningAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "RS256",
		},
		JwkSet: types.JwkSet{
			Inherited: false,
			Value:     rcsJwks,
		},
		JwkStoreCacheMissCacheTime: types.InheritedValueInt{
			Inherited: false,
			Value:     0,
		},
		RemoteConsentResponseEncryptionMethod: types.InheritedValueString{
			Inherited: false,
			Value:     "A128GCM",
		},
		RemoteConsentRedirectURL: types.InheritedValueString{
			Inherited: false,
			Value:     fmt.Sprintf("https://%s/rcs/ui/consent", common.Config.Hosts.IgFQDN),
		},
		RemoteConsentRequestEncryptionEnabled: types.InheritedValueBool{
			Inherited: false,
			Value:     false,
		},
		RemoteConsentRequestEncryptionMethod: types.InheritedValueString{
			Inherited: false,
			Value:     "A128GCM",
		},
		RemoteConsentResponseEncryptionAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "RSA-OAEP-256",
		},
		RequestTimeLimit: types.InheritedValueInt{
			Inherited: false,
			Value:     common.Config.Identity.RemoteConsentTimeLimitSeconds,
		},
		JwksURI: types.InheritedValueString{
			Inherited: false,
			//Value:     "http://securebanking-openbanking-uk-rcs:8080/api/rcs/consent/jwk_pub",
		},
		Type: types.Type{
			ID:         "RemoteConsentAgent",
			Name:       "OAuth2 Remote Consent Service",
			Collection: true,
		},
		Userpassword: common.Config.Ig.IgRcsSecret,
	}
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/RemoteConsentAgent/" + remoteConsentId

	s := httprest.Client.Put(path, rc, map[string]string{
		"Accept":             "*/*",
		"Connection":         "keep-alive",
		"Accept-API-Version": "protocol=2.0,resource=1.0",
	})

	zap.S().Infow("Remote Consent Service", "statusCode", s)
}

func remoteConsentExists(name string) bool {
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/RemoteConsentAgent?_queryFilter=true&_pageSize=10&_fields=agentgroup"
	consent := &types.AmResult{}
	b, _ := httprest.Client.Get(path, map[string]string{
		"Accept":             "application/json",
		"X-Requested-With":   "ForgeRock Identity Cloud Postman Collection",
		"Accept-Api-Version": "protocol=2.0,resource=1.0",
	})

	err := json.Unmarshal(b, consent)
	if err != nil {
		panic(err)
	}

	return common.Find(name, consent, func(r *types.Result) string {
		return r.ID
	})
}

func CreateSoftwarePublisherAgentOBRI() {
	if softwarePublisherAgentExists(common.Config.Identity.ObriSoftwarePublisherAgent) {
		zap.L().Info("Skipping creation of Software publisher agent")
		return
	}

	zap.L().Info("Creating software publisher agent")
	pa := types.PublisherAgent{
		PublicKeyLocation: types.InheritedValueString{
			Inherited: false,
			Value:     "jwks_uri",
		},
		JwksCacheTimeout: types.InheritedValueInt{
			Inherited: false,
			Value:     3600000,
		},
		SoftwareStatementSigningAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "PS256",
		},
		JwkSet: types.JwkSet{
			Inherited: false,
		},
		Issuer: types.InheritedValueString{
			Inherited: false,
			Value:     "ForgeRock",
		},
		JwkStoreCacheMissCacheTime: types.InheritedValueInt{
			Inherited: false,
			Value:     60000,
		},
		JwksURI: types.InheritedValueString{
			Inherited: false,
			Value:     "https://service.directory.ob.forgerock.financial/api/directory/keys/jwk_uri",
		},
	}
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/SoftwarePublisher/" + common.Config.Identity.ObriSoftwarePublisherAgent
	s := httprest.Client.Put(path, pa, map[string]string{
		"Accept":             "*/*",
		"Connection":         "keep-alive",
		"Accept-API-Version": "protocol=2.0,resource=1.0",
	})

	zap.S().Infow("Software Publisher Agent", "statusCode", s)
}

func CreateSoftwarePublisherAgentOBTestDirectory() {
	if softwarePublisherAgentExists(common.Config.Identity.ObTestDirectorySoftwarePublisherAgent) {
		zap.L().Info("Skipping creation of Software publisher agent")
		return
	}

	zap.S().Infof("Creating OB Test Directory software publisher agent '%s'", common.Config.Identity.ObTestDirectorySoftwarePublisherAgent)
	pa := types.PublisherAgent{
		PublicKeyLocation: types.InheritedValueString{
			Inherited: false,
			Value:     "jwks_uri",
		},
		JwksCacheTimeout: types.InheritedValueInt{
			Inherited: false,
			Value:     3600000,
		},
		SoftwareStatementSigningAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "PS256",
		},
		JwkSet: types.JwkSet{
			Inherited: false,
		},
		Issuer: types.InheritedValueString{
			Inherited: false,
			Value:     "OpenBanking Ltd",
		},
		JwkStoreCacheMissCacheTime: types.InheritedValueInt{
			Inherited: false,
			Value:     60000,
		},
		JwksURI: types.InheritedValueString{
			Inherited: false,
			Value:     "https://" + common.Config.Hosts.IgFQDN + "/jwkms/jwksproxy/keystore.openbankingtest.org.uk/keystore/openbanking.jwks",
		},
	}
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/SoftwarePublisher/" + common.Config.Identity.ObTestDirectorySoftwarePublisherAgent
	s := httprest.Client.Put(path, pa, map[string]string{
		"Accept":             "*/*",
		"Connection":         "keep-alive",
		"Accept-API-Version": "protocol=2.0,resource=1.0",
	})

	zap.S().Infow("Software Publisher Agent", "statusCode", s)
}

func CreateSoftwarePublisherAgentTestPublisher() {
	if softwarePublisherAgentExists(common.Config.Identity.SecureApiGatewayDevTrustedDirectory) {
		zap.L().Info("Skipping creation of Software publisher agent")
		return
	}

	zap.S().Infof("Creating software publisher agent '%s'", common.Config.Identity.SecureApiGatewayDevTrustedDirectory)
	pa := types.PublisherAgent{
		Userpassword: common.Config.Ig.IgSsaSecret,
		PublicKeyLocation: types.InheritedValueString{
			Inherited: false,
			Value:     "jwks_uri",
		},
		JwksCacheTimeout: types.InheritedValueInt{
			Inherited: false,
			Value:     3600000,
		},
		SoftwareStatementSigningAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "PS256",
		},
		JwkSet: types.JwkSet{
			Inherited: false,
		},
		Issuer: types.InheritedValueString{
			Inherited: false,
			Value:     "test-publisher",
		},
		JwkStoreCacheMissCacheTime: types.InheritedValueInt{
			Inherited: false,
			Value:     60000,
		},
		JwksURI: types.InheritedValueString{
			Inherited: false,
			Value:     "https://" + common.Config.Hosts.IgFQDN + "/jwkms/testdirectory/jwks",
		},
	}
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/SoftwarePublisher/" + common.Config.Identity.SecureApiGatewayDevTrustedDirectory
	s := httprest.Client.Put(path, pa, map[string]string{
		"Accept":             "*/*",
		"Connection":         "keep-alive",
		"Accept-API-Version": "protocol=2.0,resource=1.0",
	})

	zap.S().Infow("Software Publisher Agent", "statusCode", s)
}

func softwarePublisherAgentExists(name string) bool {
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/SoftwarePublisher?_queryFilter=true&_pageSize=10&_fields=agentgroup"
	agent := &types.AmResult{}
	b, _ := httprest.Client.Get(path, map[string]string{
		"Accept":             "application/json",
		"X-Requested-With":   "ForgeRock Identity Cloud Postman Collection",
		"Accept-Api-Version": "protocol=2.0,resource=1.0",
	})

	err := json.Unmarshal(b, agent)
	if err != nil {
		panic(err)
	}

	return common.Find(name, agent, func(r *types.Result) string {
		return r.ID
	})
}

// CreateOIDCClaimsScript -
func CreateOIDCClaimsScript(cookie *http.Cookie) string {

	zap.L().Info("Creating OIDC claims script")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking + "oidc-script.json")
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("https://%s/am/json/"+common.Config.Identity.AmRealm+"/scripts/?_action=create", common.Config.Hosts.IdentityPlatformFQDN)

	claimsScript := &types.RequestScript{}

	err = json.Unmarshal(b, claimsScript)
	if err != nil {
		panic(err)
	}

	id := httprest.GetScriptIdByName(claimsScript.Name)
	if id != "" {
		zap.L().Info("Script exists")
		return id
	}

	resp, err := restClient.R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/json").
		SetHeader("Connection", "keep-alive").
		SetHeader("Accept-API-Version", "protocol=2.0,resource=1.0").
		SetContentLength(true).
		SetCookie(cookie).
		SetResult(claimsScript).
		SetBody(b).
		Post(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	zap.S().Infow("OIDC claims script", "statusCode", resp.StatusCode(), "claimsScriptID", claimsScript.ID, "createdBy", claimsScript.CreatedBy)
	return claimsScript.ID
}

// UpdateOAuth2Provider - update the oauth 2 provider, must supply the claimScript ID
func UpdateOAuth2Provider(claimsScriptID string) {
	zap.S().Info("UpdateOAuth2Provider() Creating OAuth2Provider service in the " + common.Config.Identity.AmRealm + " realm")

	oauth2Provider := &types.OAuth2Provider{}
	fileName := ""
	if common.Config.Environment.SapigType == "core" {
		zap.S().Info("UpdateOAuth2Provider() Using Core Config")
		fileName = "oauth2provider-core-update.json"
	} else if common.Config.Environment.SapigType == "ob" {
		zap.S().Info("UpdateOAuth2Provider() Using OB Config")
		fileName = "oauth2provider-update.json"
	} else {
		panic(fmt.Sprintf("Unrecognised SapigType %v", common.Config.Environment.SapigType))
	}
    zap.S().Infof("Shaun Test %+v", oauth2Provider)
	err := common.Unmarshal(common.Config.Environment.Paths.ConfigSecureBanking+fileName, &common.Config, oauth2Provider)
	if err != nil {
		panic(err)
	}

	if oauth2ProviderExists(oauth2Provider.Type.ID) {
		zap.L().Info("UpdateOAuth2Provider() OAuth2 provider exists")
		return
	}

	zap.S().Infof("Pushing the following config %+v", oauth2Provider)

	oauth2Provider.PluginsConfig.OidcClaimsScript = claimsScriptID
	zap.S().Infow("UpdateOAuth2Provider() Updating OAuth2 provider", "claimScriptId", claimsScriptID)
	path := "/am/json/" + common.Config.Identity.AmRealm + "/realm-config/services/oauth-oidc"
	zap.S().Info("UpdateOAuth2Provider() Updating OAuth2Provider via the following path {}", path)
	s := httprest.Client.Put(path, oauth2Provider, map[string]string{
		"Accept":           "*/*",
		"Content-Type":     "application/json",
		"Connection":       "keep-alive",
		"X-Requested-With": "ForgeRock Identity Cloud Postman Collection",
	})

	zap.S().Infow("UpdateOAuth2Provider() OAuth2 provider", "statusCode", s)
}

func oauth2ProviderExists(id string) bool {
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/services?_queryFilter=true"
	r := &types.AmResult{}
	b, _ := httprest.Client.Get(path, map[string]string{
		"Accept":             "application/json",
		"X-Requested-With":   "ForgeRock Identity Cloud Postman Collection",
		"Accept-Api-Version": "protocol=1.0,resource=1.0",
	})

	err := json.Unmarshal(b, r)
	if err != nil {
		panic(err)
	}

	return common.Find(id, r, func(r *types.Result) string {
		return r.ID
	})
}

func CreateBaseURLSourceService(cookie *http.Cookie) {
	zap.S().Info("Creating BaseURLSource service in the " + common.Config.Identity.AmRealm + " realm")

	s := &types.Source{}
	err := common.Unmarshal(common.Config.Environment.Paths.ConfigSecureBanking+"create-base-url-source.json", &common.Config, s)
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("https://%s/am/json/realms/root/realms/"+common.Config.Identity.AmRealm+"/realm-config/services/baseurl?_action=create",
		common.Config.Hosts.IdentityPlatformFQDN)

	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-API-Version", "protocol=1.0,resource=1.0").
		SetHeader("Content-Type", "application/json").
		SetContentLength(true).
		SetCookie(cookie).
		SetBody(s).
		Post(path)

	zap.S().Info("resp is " + resp.String())
	if resp != nil && resp.StatusCode() == 409 {
		zap.S().Info("Did not create BaseURLSource service in " + common.Config.Identity.AmRealm + " realm. It already exists.")
	} else {
		common.RaiseForStatus(err, resp.Error(), resp.StatusCode())
		zap.S().Info("Created Base URL Service in AM's " + common.Config.Identity.AmRealm + " realm")
	}
}
