package platform

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"secure-banking-uk-initializer/pkg/common"
	"secure-banking-uk-initializer/pkg/types"
	"strings"
	_ "strings"
)

func GetCookieNameFromAm() string {
	zap.L().Info("Getting Cookie name from Identity Platform")
	path := fmt.Sprintf("%s://%s/am/json/serverinfo/*", common.Config.Hosts.Scheme, common.Config.Hosts.IdentityPlatformFQDN)
	zap.S().Infow("Getting Cookie name from Identity Platform", "path", path)
	result := &types.ServerInfo{}
	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetResult(result).
		Get(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	cookieName := result.CookieName

	zap.S().Infow("Got cookie from am",
		zap.String("cookieName", cookieName))
	return cookieName
}

// FromUserSession - get a session token from AM for authentication
//    returns the Session object with embedded session cookie
func FromUserSession(cookieName string) *common.Session {
	zap.L().Info("Getting an admin session from Identity Platform")

	path := ""
	platformType := common.Config.Environment.Type
	if platformType == "FIDC" {
		path = fmt.Sprintf("https://%s/am/json/realms/root/authenticate", common.Config.Hosts.IdentityPlatformFQDN)
	} else {
		path = fmt.Sprintf("https://%s/am/json/realms/root/authenticate?authIndexType=service&authIndexValue=ldapService", common.Config.Hosts.IdentityPlatformFQDN)
	}

	zap.S().Infow("Path to authenticate the user", "path", path)

	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-API-Version", "resource=2.0, protocol=1.0").
		SetHeader("X-OpenAM-Username", common.Config.Users.FrPlatformAdminUsername).
		SetHeader("X-OpenAM-Password", common.Config.Users.FrPlatformAdminPassword).
		Post(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())
	zap.S().Infof("Got response code %v from %v", resp.StatusCode(), path)

	if platformType == "FIDC" {
		dismiss2faDialog(path, resp)
	}

	var cookieValue = ""
	for _, cookie := range resp.Cookies() {
		zap.S().Infow("Cookies found", "cookie", cookie)
		if cookie.Name == cookieName {
			cookieValue = cookie.Value
		}
	}

	if cookieValue == "" {
		zap.S().Fatalw("Cannot find cookie",
			"statusCode", resp.StatusCode(),
			"cookieName", cookieName,
			"advice", `Calling this method twice might result in this error,
				 AM will not react well to successive session requests`,
			"error", resp.Error())
	}

	c := &http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Domain:   common.Config.Hosts.IdentityPlatformFQDN,
	}

	s := &common.Session{}
	s.Cookie = c
	zap.S().Infow("New Identity Platform session created", "cookie", s.Cookie)
	return s
}

type AuthIdFromJson struct {
	AuthId string `json:"authId,omitempty"`
}

func dismiss2faDialog(requestUri string, resp *resty.Response) {
	bodyString := string(resp.Body()[:])
	zap.S().Infof("Dismissing 2FA dialog as authing to FIDC. Auth response body is %v", bodyString)

	var jsonMap map[string]interface{}
	json.Unmarshal(resp.Body(), &jsonMap)
	authIdString := fmt.Sprintf("%q", jsonMap["authId"])
	zap.S().Infof("authId is %v", authIdString)

	cookies := resp.Cookies()

	content, erro := ioutil.ReadFile(common.Config.Environment.Paths.ConfigAuthHelper + "FidcDismiss2FA.json")
	if erro != nil {
		// Handle error here
		zap.S().Warnf("Failed to read file %v", common.Config.Environment.Paths.ConfigAuthHelper+"FidcDismiss2FA.json")
	} else {
		zap.S().Infof("File content is %v", string(content))
	}

	contentString := string(content)
	substitutedContentString := strings.ReplaceAll(contentString, "{{authId}}", authIdString)
	zap.S().Infof("Substituted body of request: %v", substitutedContentString)

	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-API-Version", "resource=2.1, protocol=1.0").
		SetHeader("Content-Type", "application/json").
		SetCookies(cookies).
		SetBody(substitutedContentString).
		Post(requestUri)

	if err != nil {
		zap.S().Warnf("Failed to dismiss 2FA. ErrorCode %v", resp.StatusCode())
	} else {
		var jsonMap map[string]interface{}
		json.Unmarshal(resp.Body(), &jsonMap)
		zap.S().Infof("Dismissed 2FA - statusCode: %v,  %v", resp.StatusCode(), jsonMap)
	}
}
