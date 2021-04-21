package platform

import (
	"time"
	"go.uber.org/zap"
	"github.com/spf13/viper"
)

//IsValidX509 - Pings the am health endpoint; an invalid certificate will return an error.
//  returns false if no valid cert present after 10 minutes. This is a naive implementation
//  and assumes the error thrown is related to an invalid SSL
func IsValidX509() bool {
	url := "https://" + viper.GetString("IAM_FQDN") + "/am/json/health/live"

	for i := 0; i < 60; i++ {
		zap.L().Info("Waiting for valid SSL certificate")
		_, err := client.R().
			Get(url)
		if err == nil {
			zap.L().Info("Got valid SSL cert")
			return true
		}
		zap.L().Info(err.Error())
		time.Sleep(10 * time.Second)
	}
	return false
}