package backend

import (
	"github.com/pearsontechnology/terraform-provider-secret/backend/gpg"
	"github.com/pearsontechnology/terraform-provider-secret/backend/kms"
)

// SecretBackend is an interface that need to be implemented for different decryption models.
type SecretBackend interface {
	Configure(map[string]interface{}) error
	Decrypt(string) (string, error)
	Validate() error
}

func Plugins() map[string]SecretBackend {
	return map[string]SecretBackend{
		"kms": &kms.KMSsecret{},
		"gpg": &gpg.GPGsecret{},
	}
}
