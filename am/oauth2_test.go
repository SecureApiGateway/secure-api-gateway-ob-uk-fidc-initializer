package am

import (
	"github.com/secureBankingAccessToolkit/securebanking-openbanking-uk-fidc-initialiszer/common"
	"io/ioutil"
	"testing"

	mocks "github.com/secureBankingAccessToolkit/securebanking-openbanking-uk-fidc-initialiszer/mocks/am"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindSoftwarePublisherAgent(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	common.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("oauth2-test.json")
	mockRestReaderWriter.On("Get", mock.Anything, mock.Anything).
		Return(buffer)

	b := SoftwarePublisherAgentExists("OBRI")
	assert.True(t, b)
	mockRestReaderWriter.AssertCalled(t, "Get", mock.Anything, mock.Anything)

	b = SoftwarePublisherAgentExists("test-publisher")
	assert.True(t, b)
}

func TestFindRemoteConsent(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	common.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("remote-consent-test.json")
	mockRestReaderWriter.On("Get", mock.Anything, mock.Anything).
		Return(buffer)

	b := RemoteConsentExists("forgerock-rcs")

	assert.True(t, b)
}

func TestReturnScriptId(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	common.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("script-test.json")
	mockRestReaderWriter.On("Get", mock.Anything, mock.Anything).
		Return(buffer)

	s := GetScriptIdByName("test script")

	assert.Equal(t, "123", s)
	mockRestReaderWriter.AssertCalled(t, "Get", mock.Anything, mock.Anything)

	s = GetScriptIdByName("Doesnt existy")
	assert.Equal(t, "", s)
}

func TestFindOAuth2Provider(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	common.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("oauth2provider-test.json")
	mockRestReaderWriter.On("Get", mock.Anything, mock.Anything).
		Return(buffer)

	b := Oauth2ProviderExists("oauth-oidc")

	assert.True(t, b)
}
