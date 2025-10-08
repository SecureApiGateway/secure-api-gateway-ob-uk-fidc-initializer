package securebanking

import (
	"fmt"
	"net/http"
	"secure-banking-uk-initializer/pkg/common"
	"secure-banking-uk-initializer/pkg/types"

	"go.uber.org/zap"
)

// configureGoogleServiceAccount Configures Google Service Account Service in AM
//
// For CDK environments: the service account will be created if it does not exist
// For all environments: secrets will be mapped to the store if they are defined in config
func ConfigureGoogleServiceAccount(cookie *http.Cookie) {
	storeName := common.Config.Identity.GoogleSecretStoreName
	if storeName == "" {
		zap.S().Infow("No Google Secret Stores found in config, nothing to do.")
		return
	}
	if common.Config.Environment.CloudType == "CDK" {
		configureCdkGoogleServiceAccount(cookie)
	}
}

func configureCdkGoogleServiceAccount(cookie *http.Cookie) {
	serviceAccount := types.GoogleServiceAccount{
		Id:                           "default",
		AllowedRealms:                []string{"*"},
		DisallowedSecretNamePatterns: []string{},
	}
	createServiceAccountUrl, createRequest := buildCreateServiceAccountRequest(serviceAccount)
	zap.S().Infow("Attempting to create Google Service Account Servce", "service_account", serviceAccount,
		"requestUrl", createServiceAccountUrl, "requestJson", createRequest)
	CreateOrUpdateCrestResource("POST", createServiceAccountUrl, createRequest, cookie)
}

func buildCreateServiceAccountRequest(serviceAccount types.GoogleServiceAccount) (string, map[string]interface{}) {
	requestBody := make(map[string]interface{})
	requestBody["_id"] = serviceAccount.Id
	requestBody["allowedRealms"] = serviceAccount.AllowedRealms
	requestBody["disallowedSecretNamePatterns"] = serviceAccount.DisallowedSecretNamePatterns

	createServiceAccountUrl := fmt.Sprintf("https://%s/am/json/global-config/services/GoogleCloudServiceAccountService/serviceAccounts?_action=create",
		common.Config.Hosts.IdentityPlatformFQDN)
	return createServiceAccountUrl, requestBody
}
