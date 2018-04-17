{{tfDoNotEditStamp}}

module "alerts-dev-myapp" {
  source = "git@github.com:Automatic/taskhawk-terraform//alert?ref={{VERSION}}"
  queue  = "DEV-MYAPP"

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
}

module "alerts-dev-secondapp" {
  source = "git@github.com:Automatic/taskhawk-terraform//alert?ref={{VERSION}}"
  queue  = "DEV-SECONDAPP"

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
}
