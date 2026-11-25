package securebanking

import (
	"encoding/base64"
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
		// SH-ToDO - Check if RS FQDN?
		RemoteConsentRedirectURL: types.InheritedValueString{
			Inherited: false,
			Value:     fmt.Sprintf("https://%s/rcs/ui/consent", common.Config.Hosts.RSFQDN),
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

	if err := json.Unmarshal(b, consent); err != nil {
		panic(err)
	}

	return common.Find(name, consent, func(r *types.Result) string {
		return r.ID
	})
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
			Value:     "https://keystore.openbankingtest.org.uk/keystore/openbanking.jwks",
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
			Value:     "https://" + common.Config.Hosts.TrustedDirFQDN + "/api/directory/jwks",
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

	if err := json.Unmarshal(b, agent); err != nil {
		panic(err)
	}

	return common.Find(name, agent, func(r *types.Result) string {
		return r.ID
	})
}

// CreateOIDCClaimsScript -
func CreateOIDCClaimsScript(cookie *http.Cookie) string {
    return CreateScript(cookie, "oidc-script", "groovy")
}
// CreateAccessTokenModificationScript -
func CreateAccessTokenModificationScript(cookie *http.Cookie) string {
    return CreateScript(cookie, "oauth2-access-token-modification-script", "js")
}
// CreateDCRScript -
func CreateDCRScript(cookie *http.Cookie) string {
    return CreateScript(cookie, "oauth2-dcr-script", "js")
}

func CreateScript(cookie *http.Cookie, scriptIdentifier string, scriptExt string) string {

	zap.L().Info("Creating script " + scriptIdentifier)
	jsonBytes, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking + scriptIdentifier + ".json")
	if err != nil {
		panic(err)
	}

	requestPayload := &types.RequestScript{}

	if err = json.Unmarshal(jsonBytes, requestPayload); err != nil {
		panic(err)
	}

	if id := httprest.GetScriptIdByName(requestPayload.Name); id != "" {
		zap.L().Info("Script exists")
		return id
	}
	scriptBytes, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking + scriptIdentifier + "." +
	scriptExt)
	if err != nil {
		panic(err)
	}
	requestPayload.Script = base64.StdEncoding.EncodeToString(scriptBytes)

	path := fmt.Sprintf("https://%s/am/json/"+common.Config.Identity.AmRealm+"/scripts/?_action=create", common.Config.Hosts.IdentityPlatformFQDN)
	resultPayload := &types.RequestScript{}
	resp, err := restClient.R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/json").
		SetHeader("Connection", "keep-alive").
		SetHeader("Accept-API-Version", "protocol=2.0,resource=1.0").
		SetContentLength(true).
		SetCookie(cookie).
		SetResult(resultPayload).
		SetBody(requestPayload).
		Post(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	zap.S().Infow("Script", "statusCode", resp.StatusCode(), "scriptID", resultPayload.ID, "createdBy", resultPayload.CreatedBy)
	return resultPayload.ID
}

// UpdateOAuth2Provider - update the oauth 2 provider, must supply the claimScript ID
func ConfigureOBOAuth2Provider(claimsScriptID string, dcrScriptId string, accessTokenModScriptId string) {
	zap.S().Info("ConfigureOAuth2Provider() Configuring the OAuth2Provider service in the " + common.Config.Identity.AmRealm + " realm")
	oauth2Provider := &types.OBOAuth2Provider{}
	if err := common.Unmarshal(common.Config.Environment.Paths.ConfigSecureBanking+"oauth2provider-update.json", &common.Config, oauth2Provider); err != nil {
		panic(err)
	}

	path := "/am/json/" + common.Config.Identity.AmRealm + "/realm-config/services/oauth-oidc"
	if currentConfig := readOAuth2Provider(path); currentConfig != nil {
		zap.L().Info("ConfigureOAuth2Provider() OAuth2 provider exists, wiring plugin scripts")
		wireProviderScripts(currentConfig, claimsScriptID, dcrScriptId, accessTokenModScriptId)
		return
	}

	zap.S().Infof("Pushing the following config %+v", oauth2Provider)

	oauth2Provider.PluginsConfig.OidcClaimsScript = claimsScriptID
	zap.S().Infow("ConfigureOAuth2Provider() Configuring OAuth2 provider", "claimScriptId", claimsScriptID)
	oauth2Provider.ClientDynamicRegistrationConfig.DynamicClientRegistrationScript = dcrScriptId
	zap.S().Infow("ConfigureOAuth2Provider() Configuring OAuth2 provider", "dcrScriptId", dcrScriptId)
	oauth2Provider.PluginsConfig.AccessTokenModificationScript = accessTokenModScriptId
	zap.S().Infow("ConfigureOAuth2Provider() Configuring OAuth2 provider", "accessTokenModScriptId", accessTokenModScriptId)
	zap.S().Info("ConfigureOAuth2Provider() Configuring the OAuth2Provider via the following path {}", path)
	s := httprest.Client.Put(path, oauth2Provider, map[string]string{
		"Accept":           "*/*",
		"Content-Type":     "application/json",
		"Connection":       "keep-alive",
		"X-Requested-With": "ForgeRock Identity Cloud Postman Collection",
	})

	zap.S().Infow("ConfigureOAuth2Provider() OAuth2 provider", "statusCode", s)
}

func ConfigureCoreOAuth2Provider(claimsScriptID string, dcrScriptId string, accessTokenModScriptId string) {
	zap.S().Info("ConfigureOAuth2Provider() Configuring the OAuth2Provider service in the " + common.Config.Identity.AmRealm + " realm")
	oauth2Provider := &types.CoreOAuth2Provider{}
	if err := common.Unmarshal(common.Config.Environment.Paths.ConfigSecureBanking+"oauth2provider-core-update.json", &common.Config, oauth2Provider); err != nil {
		panic(err)
	}

	path := "/am/json/" + common.Config.Identity.AmRealm + "/realm-config/services/oauth-oidc"
	if currentConfig := readOAuth2Provider(path); currentConfig != nil {
		zap.L().Info("ConfigureOAuth2Provider() OAuth2 provider exists, wiring plugin scripts")
		wireProviderScripts(currentConfig, claimsScriptID, dcrScriptId, accessTokenModScriptId)
		return
	}

	zap.S().Infof("Pushing the following config %+v", oauth2Provider)

	oauth2Provider.PluginsConfig.OidcClaimsScript = claimsScriptID
	zap.S().Infow("ConfigureOAuth2Provider() Configuring OAuth2 provider", "claimScriptId", claimsScriptID)
	oauth2Provider.ClientDynamicRegistrationConfig.DynamicClientRegistrationScript = dcrScriptId
	zap.S().Infow("ConfigureOAuth2Provider() Configuring OAuth2 provider", "dcrScriptId", dcrScriptId)
	oauth2Provider.PluginsConfig.AccessTokenModificationScript = accessTokenModScriptId
	zap.S().Infow("ConfigureOAuth2Provider() Configuring OAuth2 provider", "accessTokenModScriptId", accessTokenModScriptId)
	zap.S().Info("ConfigureOAuth2Provider() Configuring the OAuth2Provider via the following path {}", path)
	s := httprest.Client.Put(path, oauth2Provider, map[string]string{
		"Accept":           "*/*",
		"Content-Type":     "application/json",
		"Connection":       "keep-alive",
		"X-Requested-With": "ForgeRock Identity Cloud Postman Collection",
	})

	zap.S().Infow("ConfigureOAuth2Provider() OAuth2 provider", "statusCode", s)
}

// readOAuth2Provider reads the OAuth2 provider configuration, returning nil when the
// service does not exist. GetRS is used instead of Get because the latter treats a
// 404 as fatal in strict mode.
func readOAuth2Provider(path string) map[string]interface{} {
	b, status := httprest.Client.GetRS("https://"+common.Config.Hosts.IdentityPlatformFQDN+path, map[string]string{
		"Accept":           "application/json",
		"X-Requested-With": "ForgeRock Identity Cloud Postman Collection",
	})
	if status == 404 {
		return nil
	}
	if status != 200 {
		zap.S().Fatalw("readOAuth2Provider() Unexpected status reading the OAuth2 provider", "httpStatus", status)
	}

	provider := map[string]interface{}{}
	if err := json.Unmarshal(b, &provider); err != nil {
		panic(err)
	}
	return provider
}

// wireProviderScripts wires the OIDC claims, dynamic client registration and access
// token modification plugin scripts into the existing OAuth2 provider. The provider
// config is updated in place — only the three script references change — and PUT
// back, leaving every other provider setting untouched. The services endpoint does
// not support PATCH.
func wireProviderScripts(provider map[string]interface{}, claimsScriptID string, dcrScriptId string, accessTokenModScriptId string) {
	delete(provider, "_id")
	delete(provider, "_rev")

	plugins, ok := provider["pluginsConfig"].(map[string]interface{})
	if !ok {
		plugins = map[string]interface{}{}
		provider["pluginsConfig"] = plugins
	}
	plugins["oidcClaimsScript"] = claimsScriptID
	plugins["accessTokenModificationScript"] = accessTokenModScriptId

	registration, ok := provider["clientDynamicRegistrationConfig"].(map[string]interface{})
	if !ok {
		registration = map[string]interface{}{}
		provider["clientDynamicRegistrationConfig"] = registration
	}
	registration["dynamicClientRegistrationScript"] = dcrScriptId

	zap.S().Infow("wireProviderScripts() Wiring OAuth2 provider plugin scripts",
		"claimScriptId", claimsScriptID,
		"dcrScriptId", dcrScriptId,
		"accessTokenModScriptId", accessTokenModScriptId)
	path := "/am/json/" + common.Config.Identity.AmRealm + "/realm-config/services/oauth-oidc"
	status := httprest.Client.Put(path, provider, map[string]string{
		"Accept":           "*/*",
		"Content-Type":     "application/json",
		"Connection":       "keep-alive",
		"X-Requested-With": "ForgeRock Identity Cloud Postman Collection",
	})

	zap.S().Infow("wireProviderScripts() OAuth2 provider", "statusCode", status)
}

func CreateBaseURLSourceService(cookie *http.Cookie) {
	zap.S().Info("Creating BaseURLSource service in the " + common.Config.Identity.AmRealm + " realm")

	s := &types.Source{}
	if err := common.Unmarshal(common.Config.Environment.Paths.ConfigSecureBanking+"create-base-url-source.json", &common.Config, s); err != nil {
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
