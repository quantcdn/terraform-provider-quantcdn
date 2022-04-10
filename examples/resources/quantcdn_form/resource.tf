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
