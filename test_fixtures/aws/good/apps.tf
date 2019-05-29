{{tfDoNotEditStamp}}

module "app-dev-myapp" {
  source  = "Automatic/taskhawk-queue/aws"
  version = "~> {{TFAWSQueueModuleVersion}}"

  queue            = "DEV-MYAPP"
  iam              = "true"
  alerting         = "true"
  enable_scheduler = "true"
  aws_region       = "${var.aws_region}"
  aws_account_id   = "${var.aws_account_id}"

  tags = {
    App = "myapp"
    Env = "dev"
  }

  dlq_alarm_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  dlq_ok_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_alarm_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_ok_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_alarm_high_message_count_threshold               = 5000
  queue_alarm_high_priority_high_message_count_threshold = 5000
  queue_alarm_low_priority_high_message_count_threshold  = 5000
  queue_alarm_bulk_high_message_count_threshold          = 5000
}

module "app-dev-secondapp" {
  source  = "Automatic/taskhawk-queue/aws"
  version = "~> {{TFAWSQueueModuleVersion}}"

  queue    = "DEV-SECONDAPP"
  iam      = "true"
  alerting = "true"

  tags = {
    App = "secondapp"
    Env = "dev"
  }

  dlq_alarm_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  dlq_ok_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_alarm_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_ok_high_message_count_actions = [
    "pager_action",
    "pager_action2",
  ]

  queue_alarm_high_message_count_threshold               = 1000
  queue_alarm_high_priority_high_message_count_threshold = 100
  queue_alarm_low_priority_high_message_count_threshold  = 50000
  queue_alarm_bulk_high_message_count_threshold          = 100000
}

module "app-dev-anotherapp" {
  source  = "Automatic/taskhawk-lambda/aws"
  version = "~> {{TFAWSLambdaModuleVersion}}"

  name               = "dev-anotherapp"
  function_arn       = "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed"
  function_name      = "myFunction"
  function_qualifier = "deployed"
}
