package securebanking

import (
	"io/ioutil"
	"secure-banking-uk-initializer/pkg/httprest"
	"testing"

	mocks "secure-banking-uk-initializer/pkg/mocks/am"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindSoftwarePublisherAgent(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	httprest.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("oauth2-test.json")
	mockRestReaderWriter.On("Get", mock.Anything, mock.Anything).
		Return(buffer, 200)

	b := softwarePublisherAgentExists("OBRI")
	assert.True(t, b)
	mockRestReaderWriter.AssertCalled(t, "Get", mock.Anything, mock.Anything)

	b = softwarePublisherAgentExists("test-publisher")
	assert.True(t, b)
}

func TestFindRemoteConsent(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	httprest.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("remote-consent-test.json")
	mockRestReaderWriter.On("Get", mock.Anything, mock.Anything).
		Return(buffer, 200)

	b := remoteConsentExists("secure-open-banking-rcs")

	assert.True(t, b)
}

func TestReturnScriptId(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	httprest.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("script-test.json")
	mockRestReaderWriter.On("Get", mock.Anything, mock.Anything).
		Return(buffer, 200)

	s := httprest.GetScriptIdByName("test script")

	assert.Equal(t, "123", s)
	mockRestReaderWriter.AssertCalled(t, "Get", mock.Anything, mock.Anything)

	s = httprest.GetScriptIdByName("Doesnt existy")
	assert.Equal(t, "", s)
}

func TestWireProviderScripts(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	httprest.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("oauth2provider-live-test.json")
	mockRestReaderWriter.On("GetRS", mock.Anything, mock.Anything).
		Return(buffer, 200)
	mockRestReaderWriter.On("Put", mock.Anything, mock.Anything, mock.Anything).
		Return(200)

	wireProviderScripts(readOAuth2Provider("https://iam/path/services/oauth-oidc"), "claims-id", "dcr-id", "atm-id")

	mockRestReaderWriter.AssertCalled(t, "GetRS", mock.Anything, mock.Anything)
	mockRestReaderWriter.AssertCalled(t, "Put", mock.Anything, mock.Anything, mock.Anything)

	putBody := mockRestReaderWriter.Calls[1].Arguments.Get(1).(map[string]interface{})
	assert.NotContains(t, putBody, "_id")
	assert.NotContains(t, putBody, "_rev")

	plugins := putBody["pluginsConfig"].(map[string]interface{})
	assert.Equal(t, "claims-id", plugins["oidcClaimsScript"])
	assert.Equal(t, "atm-id", plugins["accessTokenModificationScript"])

	registration := putBody["clientDynamicRegistrationConfig"].(map[string]interface{})
	assert.Equal(t, "dcr-id", registration["dynamicClientRegistrationScript"])

	// untouched settings survive the round-trip
	core := putBody["coreOAuth2Config"].(map[string]interface{})
	assert.Equal(t, true, core["statelessTokensEnabled"])
	assert.InEpsilon(t, 360000.0, core["accessTokenLifetime"], 0.0001)
	advanced := putBody["advancedOAuth2Config"].(map[string]interface{})
	assert.Equal(t, "ssl-client-cert", advanced["tlsClientCertificateTrustedHeader"])
}

//func TestUpdateOBOAuth2Provider(t *testing.T) {
//	mockRestReaderWriter := &mocks.RestReaderWriter{}
//	httprest.Client = mockRestReaderWriter
//	buffer, _ := ioutil.ReadFile("")
//}

func TestReadOAuth2ProviderAbsent(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	httprest.Client = mockRestReaderWriter
	mockRestReaderWriter.On("GetRS", mock.Anything, mock.Anything).
		Return([]byte("{\"code\":404,\"reason\":\"Not Found\"}"), 404)

	assert.Nil(t, readOAuth2Provider("https://iam/path/services/oauth-oidc"))
}
