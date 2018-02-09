package secret

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSecret_kms(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecretConfig(acctest.RandInt(), acctest.RandInt(), acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRandomEncrypt("data.secret.mysecret"),
					resource.TestCheckResourceAttr("data.secret.mysecret", "value", "myvalue"),
				),
			},
		},
	})
}

func testDataSourceSecretConfig(rInt1, rInt2, rInt3 int) string {
	enc := doKmsEncrypt("myvalue")
	return fmt.Sprintf(`
		provider "secret" {
			backend = "kms"
			config = {
				region = "eu-west-1"
				shared_credentials_file = "~/.aws/credentials"
			}
		}
		data "secret" "mysecret" {
			encrypted_value = "%s"
		}
	`, enc)
}

func testAccCheckRandomEncrypt(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find secret datasource: %s", n)
		}
		return nil
	}
}

func doKmsEncrypt(v string) string {
	sess := session.Must(session.NewSession())
	kmsClient := kms.New(sess, aws.NewConfig().WithRegion("eu-west-1"))
	k, _ := kmsClient.ListKeys(&kms.ListKeysInput{})
	if len(k.Keys) == 0 {
		return ""
	}
	keyID := k.Keys[0].KeyId
	retval, err := kmsClient.Encrypt(&kms.EncryptInput{
		KeyId:     keyID,
		Plaintext: []byte(v),
	})
	if err != nil {
		log.Println(err)
	}
	encoded := base64.StdEncoding.EncodeToString(retval.CiphertextBlob)
	return encoded
}
