package secret

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pearsontechnology/terraform-provider-secret/backend"
)

func DataSourceSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encrypted_value": {
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSecretRead(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(backend.SecretBackend)
	encryptedValue := d.Get("encrypted_value").(string)
	decryptedValue, err := svc.Decrypt(encryptedValue)
	if err != nil {
		return err
	}

	d.SetId(encryptedValue)
	d.Set("value", decryptedValue)
	return nil
}
