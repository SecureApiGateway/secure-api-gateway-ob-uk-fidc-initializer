package types

// GoogleServiceAccount represents a service account service in AM which is backed by Google Service account
type GoogleServiceAccount struct {
	Id                           string
	AllowedRealms                []string
	DisallowedSecretNamePatterns []string
}
