package gpg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/openpgp"
)

type GPGsecret struct {
	PrivateKey []byte
	Passphrase string
}

func (g *GPGsecret) Configure(m map[string]interface{}) error {
	if val, ok := m["private_key"]; ok {
		filename := val.(string)
		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("can't read gpg private key: %s", err.Error())
		}
		g.PrivateKey = buf
	}

	if val, ok := m["passphrase"]; ok {
		g.Passphrase = val.(string)
	}
	return nil
}

func (g *GPGsecret) Decrypt(encrypted string) (string, error) {
	base64decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode secret: %s", err.Error())
	}

	entitylist, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(g.PrivateKey))
	if err != nil {
		log.Fatal(err)
	}
	entity := entitylist[0]

	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		err := entity.PrivateKey.Decrypt([]byte(g.Passphrase))
		if err != nil {
			return "", err
		}
	}
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
			err := subkey.PrivateKey.Decrypt([]byte(g.Passphrase))
			if err != nil {
				return "", fmt.Errorf("failed to decrypt subkey: %s", err)
			}
		}
	}

	md, err := openpgp.ReadMessage(bytes.NewBuffer(base64decoded), entitylist, nil, nil)
	if err != nil {
		return "", fmt.Errorf("could not decrypt secret: %s", err.Error())
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (g *GPGsecret) Validate() error {
	return nil
}
