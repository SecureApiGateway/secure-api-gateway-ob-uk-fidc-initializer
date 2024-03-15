package main

import (
	"fmt"
	"os"
	"reflect"
	"secure-banking-uk-initializer/pkg/common"
	"secure-banking-uk-initializer/pkg/httprest"
	platform "secure-banking-uk-initializer/pkg/identity-platform"
	"secure-banking-uk-initializer/pkg/securebanking"
	"secure-banking-uk-initializer/pkg/types"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// init function is execute before main to initialize the program,
// this function is called when the package is initialized
func init() {
	fmt.Println("initializing the program with defaults.....")
	viper.AutomaticEnv()
	viper.SetDefault("ENVIRONMENT.VERBOSE", false)
	viper.SetDefault("ENVIRONMENT.STRICT", true)
	viper.SetDefault("ENVIRONMENT.VIPER_CONFIG", "default")
	viper.SetDefault("IDENTITY.AM_REALM", "alpha")
	// load default logger
	fmt.Println("initializing the default logger.....")
	loadLogger()
	loadConfiguration()
	// load logger again to update the level set in the configuration file
	loadLogger()
	checks()
	// after call 'loadConfiguration' we have an object with all configuration mapped
	if common.Config.Environment.Verbose {
		verboseProgramInfo()
	}

	if viper.GetBool("ENVIRONMENT.ONLY_CONFIG") {
		os.Exit(0)
	}
}

// verboseProgramInfo is a method to add all additional information about the program to output in the console in verbose/debug mode
func verboseProgramInfo() {
	fmt.Println("IdentityPlatformFQDN:", common.Config.Hosts.IdentityPlatformFQDN)
	zap.S().Infow("Configuration", "config", config)
}

// config to get configuration values
var config types.Configuration

func main() {
	// operations
	checkValidPlatformCert()
	session := getIdentityPlatformSession()

	// operation not supported on FIDC (identity cloud platform)
	createIdentityPlatformOAuth2AdminClient(session)

	//Make CDK looks like FIDC
	createRealm(session)
	//configure Identity platform (both CDK and FIDC) with non speicifc OB config
	createServerConfig(session)

	fmt.Println("Resty initialization....")

	//get IDM auth code
	session.Authenticate()
	//to obtain cookies values
	httprest.InitRestReaderWriter(session.Cookie, session.AuthToken.AccessToken)

	fmt.Println("Attempting to configure AM CORS Service")
	securebanking.ConfigureAmCorsService(session.Cookie)

	fmt.Println("Attempting to create AM Validation Service")
	securebanking.CreateAmValidationService(session.Cookie)

	if common.Config.Environment.Type != "FIDC" {
		fmt.Println("Attempting to configure AM Global Services Platform")
		securebanking.ConfigureAmPlatformService(session.Cookie)
	}

	fmt.Println("Attempting to configure Realm Default User Authentication Service")
	securebanking.ConfigureRealmDefaultUserAuthenticationService()

	fmt.Println("Attempting to configure Google Secret Store")
	securebanking.ConfigureGoogleSecretStore(session.Cookie)

	fmt.Println("Attempt PSD2 authentication trees initialization...")
	securebanking.CreateSecureBankingPSD2AuthenticationTrees()
	
	if common.Config.Environment.SapigType == "ob" {
    	fmt.Println("Attempt to create secure banking remote consent...")
		securebanking.CreateSecureBankingRemoteConsentService()
		
		fmt.Println("Attempt to create OBRI software publisher agent...")
		securebanking.CreateSoftwarePublisherAgentOBRI()
		
		fmt.Println("Attempt to create OB Test Directory software publisher agent...")
		securebanking.CreateSoftwarePublisherAgentOBTestDirectory()
	}
	
	fmt.Println("Attempt to create Test software publisher agent...")
	securebanking.CreateSoftwarePublisherAgentTestPublisher()
	
	fmt.Println("Attempt to create OIDC claims script..")
	id := securebanking.CreateOIDCClaimsScript(session.Cookie)
	fmt.Println("Attempt to Create OAUTH2 Providers...")
	securebanking.UpdateOAuth2Provider(id)
	securebanking.CreateBaseURLSourceService(session.Cookie)

	time.Sleep(5 * time.Second)

	platform.CreateIGServiceUser()
	platform.CreateIGOAuth2Client()
	platform.CreateIGPolicyAgent()

	platform.ApplySystemClients(session.Cookie)

	time.Sleep(5 * time.Second)

	if common.Config.Environment.SapigType == "ob" {
		fmt.Println("Attempt to Add OB Managed Objects...")
		securebanking.AddOBManagedObjects()
	}
	
	securebanking.CreateApiJwksEndpoint()

}

func loadLogger() {
	logger, e := common.ConfigureLogger()
	if e != nil {
		panic(e)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger)

	zap.ReplaceGlobals(logger)
}

func loadConfiguration() {
	fmt.Println("Load the [", viper.GetString("ENVIRONMENT.VIPER_CONFIG"), "] configuration.....")
	err := common.LoadConfigurationByEnv(viper.GetString("ENVIRONMENT.VIPER_CONFIG"))
	if err != nil {
		zap.S().Fatalw("Cannot load config:", "error", err)
	}
	config = common.Config
	zap.S().Infof("Config is %s", types.ToStr(config))
}

func checks() {
	fmt.Println("Making some checks.....")
	checkPaths()
}
func checkPaths() {
	zap.L().Debug("Checking trailing slash in paths")
	suffix := "/"
	value := reflect.ValueOf(config.Environment.Paths)
	for i := 0; i < value.NumField(); i++ {
		if !strings.HasSuffix(value.Field(i).String(), suffix) {
			zap.S().Fatalw(value.Type().Field(i).Name + " must have a trailing slash /")
		}
		zap.S().Debugw("index["+strconv.Itoa(i)+"]", "Field", value.Type().Field(i).Name, "value", value.Field(i).String())
	}
}

// Operations
func checkValidPlatformCert() {
	zap.L().Info("Check valid cert")
	if !platform.IsValidX509() {
		zap.L().Fatal("No Valid SSL certificate present in the cdk")
	}
}

func getIdentityPlatformSession() *common.Session {
	zap.L().Info("Get CookieName")
	cookieName := platform.GetCookieNameFromAm()
	zap.L().Info("Get user session")
	return platform.FromUserSession(cookieName)
}

//This creates an admin user in CDK deployment that can be used to create new config. THis does not run when initializer is run against
// FIDC
func createIdentityPlatformOAuth2AdminClient(session *common.Session) {
	// operation not supported on CDM (identity cloud platform)
	if config.Environment.Type == types.Platform.Instance().CDK {
		platform.CreateIdentityPlatformOAuth2AdminClient(session.Cookie)
	} else {
		zap.S().Infow("SKIP: Creating OAuth2Client Identity Platform admin client, the platform instance is not a CDK",
			"platform type", config.Environment.Type)
	}
}

func createRealm(session *common.Session) {
	// the alpha realm exist in identity cloud by default
	if !platform.RealmExist(session.Cookie) {
		platform.CreateRealm(session.Cookie)
	}
}

//sets org.forgerock.http.TrustTransactionHeader to true in AM's
//https://<domain>/am/ui-admin/#configure/serverDefaults/advanced config
//This makes AM trust the x-forgerock-transactionod header provided by IG and allows us to trace a trquest through their system
func createServerConfig(session *common.Session) {
	platform.CreateServerConfig(session.Cookie)
}
