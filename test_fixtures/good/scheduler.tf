{{tfDoNotEditStamp}}

module "scheduler-dev-myapp-nightly-job" {
  source  = "Automatic/taskhawk-scheduler/aws"
  version = "~> {{TFSchedulerModuleVersion}}"

  queue       = "${module.app-dev-myapp.default_queue_arn}"
  name        = "dev-myapp-nightly-job"
  description = "nightly job for sqs app"

  format_version = "v1.0"

  headers = {
    request_id = "<id>"
  }

  task = "tasks.send_email"

  args = [
    "hello@automatic.com",
    "Hello!",
    10,
  ]

  kwargs = {
    from_email = "spam@example.com"
    with_delay = 100
  }

  schedule_expression = "cron(0 10 * * ? *)"
}

module "scheduler-dev-anotherapp-nightly-job" {
  source  = "Automatic/taskhawk-scheduler/aws"
  version = "~> {{TFSchedulerModuleVersion}}"

  topic = "${module.app-dev-anotherapp.sns_topic_default_arn}"
  name  = "dev-anotherapp-nightly-job"

  function_name      = "myFunction"
  function_qualifier = "deployed"

  task = "tasks.cleanup_task"

  schedule_expression = "cron(0 5 * * ? *)"
}
