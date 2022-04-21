package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestForm(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testForm(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quantcdn_form.foo", "url", "/content/duis"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testForm() string {
	return fmt.Sprintf(`
resource "quantcdn_form" "foo" {
	url = "/content/duis"
	enabled = true
	success_message = "Great success %s"
	failure_message = "Such errors"
	mandatory_fields_message = "Not all errors"
	mandatory_fields = ["test", "update"]
	remove_fields = ["email"]
	notification_email {
		to = "test@test.com"
		from = "test@noreply.com"
		subject = "You've got mail"
		cc = "another@cc.com"
		enabled = true
		text_only = false
		include_results = false
	}
	notification_slack {
		webhook = "https://test.com.au/asdfasds"
		enabled = true
	}
}`, time.Now().String())
}
