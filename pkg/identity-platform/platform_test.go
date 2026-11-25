package platform

import (
	"io/ioutil"
	"secure-banking-uk-initializer/pkg/httprest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mocks "secure-banking-uk-initializer/pkg/mocks/am"
)

func TestFindExistingOAuth2AgentClient(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	httprest.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("client-check-test.json")
	mockRestReaderWriter.On("Get", mock.Anything, mock.Anything).
		Return(buffer, 200)

	b := httprest.OAuth2AgentClientsExist("Doesnt existy")
	assert.False(t, b)
	mockRestReaderWriter.AssertCalled(t, "Get", mock.Anything, mock.Anything)

	b = httprest.OAuth2AgentClientsExist("ig-client")
	assert.True(t, b)
}
