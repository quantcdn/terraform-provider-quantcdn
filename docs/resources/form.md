---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "quantcdn_form Resource - terraform-provider-quantcdn"
subcategory: ""
description: |-
  Manage forms configuration for a project
---

# quantcdn_form (Resource)

Manage forms configuration for a project

## Example Usage

```terraform
resource "quantcdn_form" "example" {
  url                      = "/content/duis"
  enabled                  = true
  success_message          = "Great success"
  failure_message          = "Such errors"
  mandatory_fields_message = "Not all errors"
  mandatory_fields         = ["test"]
  honeypot_fields          = ["honeypot"]
  remove_fields            = ["email"]
  notification_email {
    to              = "test@test.com"
    from            = "test@noreply.com"
    subject         = "You've got mail"
    cc              = "another@cc.com"
    enabled         = true
    text_only       = false
    include_results = false
  }
  notification_slack {
    webhook = "https://test.com.au/asdfasds"
    enabled = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `url` (String) The URL path to accept post values for this form

### Optional

- `enabled` (Boolean)
- `failure_message` (String) Text to display when the form fails to submit correctly
- `honeypot_fields` (Set of String) List of field names that are treated as honeypot fields
- `id` (String) The ID of this resource.
- `mandatory_fields` (Set of String) List of field names that are required
- `mandatory_fields_message` (String) Text to display when mandatory fields are missing from the submission
- `notification_email` (Block Set, Max: 1) Email notification configuration (see [below for nested schema](#nestedblock--notification_email))
- `notification_slack` (Block Set, Max: 1) Slack notification configuration (see [below for nested schema](#nestedblock--notification_slack))
- `remove_fields` (Set of String) List of field names to remove from submissions
- `success_message` (String) Text to display when form submission is successful

<a id="nestedblock--notification_email"></a>
### Nested Schema for `notification_email`

Required:

- `from` (String)
- `subject` (String)
- `to` (String)

Optional:

- `cc` (String)
- `enabled` (Boolean)
- `include_results` (Boolean)
- `text_only` (Boolean)


<a id="nestedblock--notification_slack"></a>
### Nested Schema for `notification_slack`

Required:

- `webhook` (String)

Optional:

- `enabled` (Boolean)


