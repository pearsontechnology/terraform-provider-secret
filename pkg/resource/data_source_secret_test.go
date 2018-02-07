package resource

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	terraformResource "github.com/hashicorp/terraform/helper/resource"
)

func TestAccSecret_basic(t *testing.T) {
	terraformResource.Test(t, terraformResource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testProviders,
		Steps: []terraformResource.TestStep{
			{
				Config: testAccCheckAwsAutoscalingGroupsConfig(acctest.RandInt(), acctest.RandInt(), acctest.RandInt()),
			},
			{
				Config: testAccCheckAwsAutoscalingGroupsConfigWithDataSource(acctest.RandInt(), acctest.RandInt(), acctest.RandInt()),
				Check: terraformResource.ComposeTestCheckFunc(
					testAccCheckAwsAutoscalingGroups("data.aws_autoscaling_groups.group_list"),
					terraformResource.TestCheckResourceAttr("data.aws_autoscaling_groups.group_list", "names.#", "3"),
				),
			},
		},
	})
}
