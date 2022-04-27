---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "quantcdn_revision Resource - terraform-provider-quantcdn"
subcategory: ""
description: |-
  Manage revisions for URLs for the project
---

# quantcdn_revision (Resource)

Manage revisions for URLs for the project



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (String) File path to markup to send to Quant
- `url` (String) The URL path to the revision

### Optional

- `find_attachments` (Boolean) If the Quant API should crawl external assets
- `id` (String) The ID of this resource.
- `published` (Boolean) The published status of the revision

