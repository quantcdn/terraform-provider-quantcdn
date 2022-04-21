package provider

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	quant "github.com/quantcdn/quant-go"
)

func resourceQuantForm() *schema.Resource {
	return &schema.Resource{
		Description: "Manage forms configuration for a project",

		CreateContext: resourceQuantFormCreate,
		ReadContext:   resourceQuantFormRead,
		UpdateContext: resourceQuantFormUpdate,
		DeleteContext: resourceQuantFormDelete,

		Schema: map[string]*schema.Schema{
			"url": {
				Description:  "The URL path to accept post values for this form",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^/"), "must start with '/'"),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"success_message": {
				Description: "Text to display when form submission is successful",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Thank you for your submission.",
			},
			"failure_message": {
				Description: "Text to display when the form fails to submit correctly",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "An error occurred. Please reload the page and try again.",
			},
			"mandatory_fields_message": {
				Description: "Text to display when mandatory fields are missing from the submission",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Some required values were missing, please try again.",
			},
			"mandatory_fields": {
				Description: "List of field names that are required",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"honeypot_fields": {
				Description: "List of field names that are treated as honeypot fields",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remove_fields": {
				Description: "List of field names to remove from submissions",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"notification_email": {
				Description: "Email notification configuration",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"to": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cc": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"from": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subject": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"text_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"include_results": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"notification_slack": {
				Description: "Slack notification configuration",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"webhook": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
		},
	}
}

func getQuantFormFromResource(d *schema.ResourceData) quant.Form {
	var email quant.FormNotificationEmail
	var slack quant.FormNotificationSlack
	var mandatoryFields []string
	var honeypotFields []string
	var removeFields []string

	if mandatoryFieldsRaw, ok := d.GetOk("mandatory_fields"); ok {
		for _, mf := range mandatoryFieldsRaw.([]interface{}) {
			mandatoryFields = append(mandatoryFields, mf.(string))
		}
	}

	if honeypotFieldsRaw, ok := d.GetOk("honeypot_fields"); ok {
		for _, hf := range honeypotFieldsRaw.([]interface{}) {
			honeypotFields = append(honeypotFields, hf.(string))
		}
	}

	if removeFieldsRaw, ok := d.GetOk("remove_fields"); ok {
		for _, rf := range removeFieldsRaw.([]interface{}) {
			removeFields = append(removeFields, rf.(string))
		}
	}

	if v, ok := d.GetOk("notification_email"); ok {
		for _, e := range v.(*schema.Set).List() {
			emailConfig := e.(map[string]interface{})
			email.To = emailConfig["to"].(string)
			email.Cc = emailConfig["cc"].(string)
			email.Subject = emailConfig["subject"].(string)
			email.Enabled = emailConfig["enabled"].(bool)
			email.Options.TextOnly = emailConfig["text_only"].(bool)
			email.Options.IncludeResults = emailConfig["include_results"].(bool)
			break
		}
	}

	if s, ok := d.GetOk("notification_slack"); ok {
		for _, e := range s.(*schema.Set).List() {
			slackConfig := e.(map[string]interface{})
			slack.Webhook = slackConfig["webhook"].(string)
			slack.Enabled = slackConfig["enabled"].(bool)
			break
		}
	}

	return quant.Form{
		Url:     d.Get("url").(string),
		Enabled: d.Get("enabled").(bool),
		Config: quant.FormConfig{
			Target:                d.Get("url").(string),
			HoneypotFields:        honeypotFields,
			MandatoryFields:       mandatoryFields,
			RemoveFields:          removeFields,
			SuccessMessage:        d.Get("success_message").(string),
			ErrorMessageGeneric:   d.Get("failure_message").(string),
			ErrorMessageMandatory: d.Get("mandatory_fields_message").(string),
			Notifications: quant.FormNotification{
				Email: email,
				Slack: slack,
			},
		},
	}
}

func resourceQuantFormCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	form := getQuantFormFromResource(d)

	client := meta.(*quant.Client)
	_, err := client.AddForm(form)

	if err != nil {
		return diag.Errorf("error retrieving form %s", err)
	}

	d.SetId(form.Url)
	return resourceQuantFormRead(ctx, d, meta)
}

func resourceQuantFormRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*quant.Client)
	f, err := client.GetForm(quant.RevisionQuery{
		Url: d.Get("url").(string),
	})

	if err != nil {
		return diag.Errorf("error retrieving form %s", err)
	}

	if f.Url == "" {
		d.SetId("")
		return diag.Errorf("empty form %s", client.ApiClient)
	}

	d.Set("url", f.Url)
	d.Set("success_message", f.Config.SuccessMessage)
	d.Set("failure_message", f.Config.ErrorMessageGeneric)
	d.Set("mandatory_fields_message", f.Config.ErrorMessageMandatory)
	d.Set("mandatory_fields", f.Config.MandatoryFields)
	d.Set("honeypot_fields", f.Config.HoneypotFields)
	d.Set("remove_fields", f.Config.RemoveFields)

	notificationEmail := make([]map[string]interface{}, 0)
	emailConfig := make(map[string]interface{})
	emailConfig["to"] = f.Config.Notifications.Email.To
	emailConfig["from"] = f.Config.Notifications.Email.From
	emailConfig["subject"] = f.Config.Notifications.Email.Subject
	emailConfig["cc"] = f.Config.Notifications.Email.Cc
	emailConfig["enabled"] = f.Config.Notifications.Email.Enabled
	emailConfig["text_only"] = f.Config.Notifications.Email.Options.TextOnly
	emailConfig["include_results"] = f.Config.Notifications.Email.Options.IncludeResults
	notificationEmail = append(notificationEmail, emailConfig)
	d.Set("notification_email", notificationEmail)

	notificationSlack := make([]map[string]interface{}, 0)
	slackConfig := make(map[string]interface{})
	slackConfig["webhook"] = f.Config.Notifications.Slack.Webhook
	slackConfig["enabled"] = f.Config.Notifications.Slack.Enabled
	notificationSlack = append(notificationSlack, slackConfig)
	d.Set("notification_slack", notificationSlack)

	return nil
}

func resourceQuantFormUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	form := getQuantFormFromResource(d)
	client := meta.(*quant.Client)
	_, err := client.UpdateForm(form)

	if err != nil {
		if strings.Contains(err.Error(), "Published version already has md5") {
			return nil
		}

		return diag.Errorf("error updating form %s", err)
	}

	return nil
}

func resourceQuantFormDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*quant.Client)
	f, err := client.GetForm(quant.RevisionQuery{
		Url: d.Get("url").(string),
	})

	f.Enabled = false
	client.UpdateForm(f)

	if err != nil {
		return diag.Errorf("error deleting form %s", err)
	}

	return nil
}
