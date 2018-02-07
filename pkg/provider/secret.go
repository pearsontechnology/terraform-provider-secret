package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pearsontechnology/terraform-provider-secret/pkg/backend"
	"github.com/pearsontechnology/terraform-provider-secret/pkg/resource"
)

var descriptions map[string]string
var backends map[string]backend.SecretBackend

func init() {
	descriptions = map[string]string{
		"backend": "The decryption backend to use. Currently supported backends\n" +
			"are kms and gpg",

		"config": "Backend specific config.",
	}

	backends = backend.Plugins()
}

func Secret() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"backend": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  descriptions["backend"],
				ValidateFunc: validateBackendType,
			},
			"config": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: descriptions["config"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"secret": resource.DataSourceSecret(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	backendType := d.Get("backend").(string)
	b := backends[backendType]
	cfg := d.Get("config").(map[string]interface{})
	err := b.Configure(cfg)
	return b, err
}

func validateBackendType(v interface{}, k string) (ws []string, es []error) {
	val := v.(string)
	if b, ok := backends[val]; ok {
		if err := b.Validate(); err != nil {
			es = append(es, err)
		}
	} else {
		es = append(es, fmt.Errorf("unsupported secret backend type: %s", val))
	}
	return
}
