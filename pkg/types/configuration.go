package types

import "fmt"

func ToStr(config Configuration) string {
	return fmt.Sprintf("Config is %#v", config)
}

type Configuration struct {
	Environment environment `mapstructure:"ENVIRONMENT"`
	Hosts       hosts       `mapstructure:"HOSTS"`
	Identity    identity    `mapstructure:"IDENTITY"`
	Ig          ig          `mapstructure:"IG"`
	Users       users       `mapstructure:"USERS"`
	TLS         tls         `mapstructure:"TLS"`
	OB          ob          `mapstructure:"OB"`
}

type hosts struct {
	BaseFQDN             string   `mapstructure:"BASE_FQDN"`
	WildcardFQDN         string   `mapstructure:"WILDCARD_FQDN"`
	IgFQDN               string   `mapstructure:"IG_FQDN"`
	MtlsFQDN             string   `mapstructure:"MTLS_FQDN"`
	IdentityPlatformFQDN string   `mapstructure:"IDENTITY_PLATFORM_FQDN"`
	Scheme               string   `mapstructure:"SCHEME"`
	IgAudienceFQDNs      []string `mapstructure:"IG_AUDIENCE_FQDNS"`
}

type identity struct {
	AmRealm                                  string `mapstructure:"AM_REALM"`
	IdmClientId                              string `mapstructure:"IDM_CLIENT_ID"`
	IdmClientSecret                          string `mapstructure:"IDM_CLIENT_SECRET"`
	RemoteConsentId                          string `mapstructure:"REMOTE_CONSENT_ID"`
	RemoteConsentSigningPublicKey            string `mapstructure:"REMOTE_CONSENT_SIGNING_PUBLIC_KEY"`
	RemoteConsentSigningKeyId                string `mapstructure:"REMOTE_CONSENT_SIGNING_KEY_ID"`
	RemoteConsentTimeLimitSeconds            int    `mapstructure:"REMOTE_CONSENT_TIME_LIMIT_SECONDS"`
	ObriSoftwarePublisherAgent               string `mapstructure:"OBRI_SOFTWARE_PUBLISHER_AGENT_NAME"`
	ObTestDirectorySoftwarePublisherAgent    string `mapstructure:"OB_TEST_DIRECTORY_SOFTWARE_PUBLISHER_NAME"`
	SecureApiGatewayDevTrustedDirectory      string `mapstructure:"SECURE_API_GATEWAY_DEVELOPMENT_TRUSTED_DIRECTORY"`
	GoogleSecretStoreName                    string `mapstructure:"GOOGLE_SECRET_STORE_NAME"`
	GoogleSecretStoreProject                 string `mapstructure:"GOOGLE_SECRET_STORE_PROJECT"`
	GoogleSecretStoreOAuth2CaCertsSecretName string `mapstructure:"GOOGLE_SECRET_STORE_OAUTH2_CA_CERTS_SECRET_NAME"`
	DefaultUserAuthenticationService         string `mapstructure:"DEFAULT_USER_AUTHENTICATION_SERVICE"`
}

type ig struct {
	IgClientId      string `mapstructure:"IG_CLIENT_ID"`
	IgClientSecret  string `mapstructure:"IG_CLIENT_SECRET"`
	IgRcsSecret     string `mapstructure:"IG_RCS_SECRET"`
	IgSsaSecret     string `mapstructure:"IG_SSA_SECRET"`
	IgIdmUser       string `mapstructure:"IG_IDM_USER"`
	IgIdmPassword   string `mapstructure:"IG_IDM_PASSWORD"`
	IgAgentId       string `mapstructure:"IG_AGENT_ID"`
	IgAgentPassword string `mapstructure:"IG_AGENT_PASSWORD"`
}
type environment struct {
	Verbose   bool   `mapstructure:"VERBOSE"`
	Strict    bool   `mapstructure:"STRICT"`
	Type      string `mapstructure:"TYPE"`
	Paths     paths  `mapstructure:"PATHS"`
	SapigType string `mapstructure:"SAPIGTYPE"`
}

type paths struct {
	ConfigBaseDirectory    string `mapstructure:"CONFIG_BASE_DIRECTORY"`
	ConfigSecureBanking    string `mapstructure:"CONFIG_SECURE_BANKING"`
	ConfigIdentityPlatform string `mapstructure:"CONFIG_IDENTITY_PLATFORM"`
	ConfigAuthHelper       string `mapstructure:"CONFIG_AUTH_HELPER"`
}

type users struct {
	FrPlatformAdminUsername string `mapstructure:"FR_PLATFORM_ADMIN_USERNAME"`
	FrPlatformAdminPassword string `mapstructure:"FR_PLATFORM_ADMIN_PASSWORD"`
}

type ob struct {
	OrganisationId string `mapstructure:"ORGANISATION_ID"`
	SoftwareId     string `mapstructure:"SOFTWARE_ID"`
}

type tls struct {
	ClientCertHeaderName string `mapstructure:"CLIENT_CERT_HEADER_NAME"`
}
