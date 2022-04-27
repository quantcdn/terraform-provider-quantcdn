package provider

import (
	"context"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	quant "github.com/quantcdn/quant-go"
)

func resourceQuantRevision() *schema.Resource {
	return &schema.Resource{
		Description: "Manage revisions for URLs for the project",

		CreateContext: resourceQuantRevisionCreate,
		ReadContext:   resourceQuantRevisionRead,
		UpdateContext: resourceQuantRevisionCreate,
		DeleteContext: resourceQuantRevisionDelete,

		Schema: map[string]*schema.Schema{
			"url": {
				Description:  "The URL path to the revision",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^/"), "must start with '/'"),
			},
			"published": {
				Type:        schema.TypeBool,
				Description: "The published status of the revision",
				Optional:    true,
				Default:     true,
			},
			"find_attachments": {
				Type:        schema.TypeBool,
				Description: "If the Quant API should crawl external assets",
				Optional:    true,
				Default:     true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "File path to markup to send to Quant",
				Required:    true,
			},
		},
	}
}

func resourceQuantRevisionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*quant.Client)

	file, err := ioutil.ReadFile(d.Get("content").(string))

	if err != nil {
		return diag.Errorf("error loading content %s", err)
	}

	r, err := client.AddMarkupRevision(quant.MarkupRevision{
		Url:             d.Get("url").(string),
		FindAttachments: d.Get("find_attachments").(bool),
		Published:       d.Get("published").(bool),
		Content:         file,
	}, false)

	if err != nil && !strings.Contains(err.Error(), "Published version already has md5") {
		return diag.Errorf("unable to create revisions %s", err)
	}

	d.SetId(r.Url)
	return resourceQuantRevisionRead(ctx, d, meta)
}

func resourceQuantRevisionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*quant.Client)
	r, err := client.GetRevision(quant.RevisionQuery{
		Url: d.Get("url").(string),
	})

	if err != nil {
		return diag.Errorf("error retrieving revision %s", err)
	}

	d.SetId(r.Url)
	d.Set("url", r.Url)
	d.Set("published", r.Published)

	return nil
}

func resourceQuantRevisionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
