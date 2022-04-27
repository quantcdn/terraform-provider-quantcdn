package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestRevision(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testRevisionDefinition(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quantcdn_revision.foo", "url", "/test/content"),
				),
			},
		},
	})
}

func testRevisionDefinition() string {
	return `
resource "quantcdn_revision" "foo" {
	url = "/test/content"
	published = true
	find_attachments = false
	content = "./fixtures/test.html"
}`
}
