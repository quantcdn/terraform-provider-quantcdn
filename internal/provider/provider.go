package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	quant "github.com/quantcdn/quant-go"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("QUANT_CLIENT_ID", nil),
					Description: "A registered Quant client name",
				},

				"project": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("QUANT_PROJECT", nil),
					Description: "A registered Quant project name",
				},

				"api_token": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("QUANT_TOKEN", nil),
					Description: "The API token for operations",
				},

				"api_hostname": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("QUANT_HOSTNAME", "https://api.quantcdn.io"),
				},

				"api_basepath": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("QUANT_BASEPATH", "/v1"),
				},
			},

			DataSourcesMap: map[string]*schema.Resource{},

			ResourcesMap: map[string]*schema.Resource{
				"quantcdn_form":     resourceQuantForm(),
				"quantcdn_revision": resourceQuantRevision(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		client := quant.NewClient(
			d.Get("api_token").(string),
			d.Get("client_id").(string),
			d.Get("project").(string),
		)

		if apiHost, ok := d.GetOk("api_hostname"); ok {
			client.Host = apiHost.(string)
		}

		if apiBase, ok := d.GetOk("api_base"); ok {
			client.Base = apiBase.(string)
		}

		return client, nil
	}
}
