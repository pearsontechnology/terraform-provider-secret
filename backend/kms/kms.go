package kms

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	kmsService "github.com/aws/aws-sdk-go/service/kms"
	homedir "github.com/mitchellh/go-homedir"
)

// KMSsecret implements SecretBackend to decrypt using AWS KMS.
type KMSsecret struct {
	AccessKey           string
	SecretKey           string
	Profile             string
	Region              string
	CredentialsFilename string
	Client              *kmsService.KMS
}

func (k *KMSsecret) Configure(m map[string]interface{}) error {
	if key, ok := m["access_key"]; ok {
		k.AccessKey = key.(string)
	}

	if key, ok := m["secret_key"]; ok {
		k.SecretKey = key.(string)
	}

	if key, ok := m["profile"]; ok {
		k.Profile = key.(string)
	}

	if key, ok := m["region"]; ok {
		k.Region = key.(string)
	}

	if key, ok := m["shared_credentials_file"]; ok {
		credsPath, err := homedir.Expand(key.(string))
		if err != nil {
			return err
		}
		k.CredentialsFilename = credsPath

	}
	// finally, configure kms.Client
	return k.configureClient()

}

func (k *KMSsecret) Decrypt(encrypted string) (string, error) {
	blob, _ := base64.StdEncoding.DecodeString(encrypted)
	result, err := k.Client.Decrypt(&kmsService.DecryptInput{CiphertextBlob: blob})
	if err != nil {
		return "", err
	}
	value := string(result.Plaintext)
	return value, nil
}

func (k *KMSsecret) configureClient() error {
	opt := session.Options{}

	cfg := aws.Config{
		Region: aws.String(k.Region),
	}

	cfg.Credentials = k.getCredentials()
	opt.Config = cfg

	sess, err := session.NewSessionWithOptions(opt)
	if err != nil {
		return err
	}

	_, err = sess.Config.Credentials.Get()
	if err != nil {
		return err
	}

	k.Client = kmsService.New(sess)
	return nil
}

func (k *KMSsecret) Validate() error {
	return nil
}

func (k *KMSsecret) getCredentials() *awsCredentials.Credentials {
	var metadataClient = ec2metadata.New(session.New(aws.NewConfig()))
	providers := []awsCredentials.Provider{
		&awsCredentials.StaticProvider{Value: awsCredentials.Value{
			AccessKeyID:     k.AccessKey,
			SecretAccessKey: k.SecretKey,
		}},
		&awsCredentials.EnvProvider{},
		&awsCredentials.SharedCredentialsProvider{
			Filename: k.CredentialsFilename,
			Profile:  k.Profile,
		},
		&ec2rolecreds.EC2RoleProvider{
			Client: metadataClient,
		},
	}
	return awsCredentials.NewChainCredentials(providers)
}
