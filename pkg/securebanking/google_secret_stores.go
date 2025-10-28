package securebanking

import (
	"fmt"
	"net/http"
	"net/url"
	"secure-banking-uk-initializer/pkg/common"
	"secure-banking-uk-initializer/pkg/types"

	"go.uber.org/zap"
)

// ConfigureGoogleSecretStore Configures Google Secret Store in AM if "GOOGLE_SECRET_STORE_NAME" is defined in config
//
// For CDK environments: the store will be created if it does not exist
// For all environments: secrets will be mapped to the store if they are defined in config
func ConfigureGoogleSecretStore(cookie *http.Cookie) {
	storeName := common.Config.Identity.GoogleSecretStoreName
	if storeName == "" {
		zap.S().Infow("No Google Secret Stores found in config, nothing to do.")
		return
	}
	if common.Config.Environment.CloudType == "CDK" {
		configureCdkGoogleSecretStore(cookie, storeName)
	}

	oAuth2CaCertsSecretName := common.Config.Identity.GoogleSecretStoreOAuth2CaCertsSecretName
	if oAuth2CaCertsSecretName != "" {
		oAuth2Secret := types.SecretMapping{
			SecretId: "am.services.oauth2.tls.client.cert.authentication",
			Alias:    oAuth2CaCertsSecretName,
		}
		configSecretMappings(storeName, []types.SecretMapping{oAuth2Secret}, cookie)
	}

	devigobsigningkey := common.Config.Identity.GoogleSecretStoreOBSigningKeySecretName
	if devigobsigningkey != "" {
		// Add in mapping for OB Signing Key
		oAuth2Signing := types.SecretMapping{
			SecretId: "am.services.oauth2.oidc.signing.RSA",
			Alias:    devigobsigningkey,
		}
		configSecretMappings(storeName, []types.SecretMapping{oAuth2Signing}, cookie)
	}
}

func configureCdkGoogleSecretStore(cookie *http.Cookie, storeName string) {
	project := common.Config.Identity.GoogleSecretStoreProject
	if project == "" {
		zap.S().Fatal("GOOGLE_SECRET_STORE_PROJECT must be configured")
	}
	store := types.GoogleSecretStore{
		Name: storeName, ServiceAccount: "default",
		Project:               project,
		SecretFormat:          "PEM",
		ExpiryDurationSeconds: 3600,
	}
	createStoreUrl, storeRequest := buildCreateStoreRequest(store)
	zap.S().Infow("Attempting to configure Google Secret Store", "store", store,
		"requestUrl", createStoreUrl, "requestJson", storeRequest)
	CreateOrUpdateCrestResource("PUT", createStoreUrl, storeRequest, cookie)
}

func buildCreateStoreRequest(store types.GoogleSecretStore) (string, map[string]interface{}) {
	requestBody := make(map[string]interface{})
	requestBody["_id"] = store.Name
	requestBody["serviceAccount"] = store.ServiceAccount
	requestBody["project"] = store.Project
	requestBody["expiryDurationSeconds"] = store.ExpiryDurationSeconds
	requestBody["secretFormat"] = store.SecretFormat

	createStoreUrl := fmt.Sprintf("https://%s/am/json/realms/root/realms/%s/realm-config/secrets/stores/GoogleSecretManagerSecretStoreProvider/%s",
		common.Config.Hosts.IdentityPlatformFQDN, common.Config.Identity.AmRealm, url.PathEscape(store.Name))
	return createStoreUrl, requestBody
}

func configSecretMappings(storeName string, secrets []types.SecretMapping, cookie *http.Cookie) {
	zap.S().Infow("Attempting to map secrets to store", "store", storeName)
	for _, secret := range secrets {
		createMappingUrl, mappingRequest := buildSecretMappingRequest(storeName, secret)
		CreateOrUpdateCrestResource("PUT", createMappingUrl, mappingRequest, cookie)
	}
}

func buildSecretMappingRequest(storeName string, secretMapping types.SecretMapping) (string, map[string]interface{}) {
	createMappingUrl := fmt.Sprintf("https://%s/am/json/realms/root/realms/%s/realm-config/secrets/stores/GoogleSecretManagerSecretStoreProvider/%s/mappings/%s",
		common.Config.Hosts.IdentityPlatformFQDN, common.Config.Identity.AmRealm, url.PathEscape(storeName), secretMapping.SecretId)

	requestBody := make(map[string]interface{})
	requestBody["secretId"] = secretMapping.SecretId
	requestBody["aliases"] = []string{secretMapping.Alias}

	return createMappingUrl, requestBody
}
