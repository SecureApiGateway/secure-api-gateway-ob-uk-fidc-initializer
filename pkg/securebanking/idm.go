package securebanking

import (
	"io/ioutil"
	"secure-banking-uk-initializer/pkg/common"
	"secure-banking-uk-initializer/pkg/httprest"

	"go.uber.org/zap"
)

func CreateApiJwksEndpoint() {
	zap.L().Info("Creating API JWKS Endpoint")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking + "create-jwks-endpoint.json")
	if err != nil {
		panic(err)
	}

	path := "/openidm/config/endpoint/apiclientjwks"
	s := httprest.Client.Put(path, b, map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/json",
		"Connection":   "keep-alive",
	})

	zap.S().Infow("JWKS endpoint", "statusCode", s)
}
