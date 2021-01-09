{{tfDoNotEditStamp}}

module "app-dev-myapp" {
  source  = "standard-ai/taskhawk-queue/google"
  version = "~> {{TFGoogleQueueModuleVersion}}"

  queue = "dev-myapp"

  enable_alerts    = var.enable_alerts
  alerting_project = var.alerting_project

  iam_service_accounts = [
    "myapp@project.iam.gserviceaccount.com"
  ]

  labels = {
    app = "myapp"
    env = "dev"
  }

  queue_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  dlq_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  dataflow_freshness_alert_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]

  queue_alarm_high_message_count_threshold               = 5000
  queue_alarm_high_priority_high_message_count_threshold = 5000
  queue_alarm_low_priority_high_message_count_threshold  = 5000
  queue_alarm_bulk_high_message_count_threshold          = 5000

  enable_firehose_all_messages = var.enable_firehose_all_messages
  dataflow_tmp_gcs_location    = var.dataflow_tmp_gcs_location
  dataflow_template_gcs_path   = var.dataflow_template_pubsub_to_storage_gcs_path
  dataflow_zone                = var.dataflow_zone
  dataflow_region              = var.dataflow_region
  dataflow_output_directory    = var.dataflow_output_directory

  scheduler_jobs = [
    {
      name        = "nightly-job"
      description = "nightly job"

      format_version = "v1.0"

      schedule = "0 10 * * ? *"
      timezone = "America/Los_Angeles"

      headers = {
        request_id = "79c7104b-1364-323e-ad09-89ed72089f98"
      }

      task = "tasks.send_email"
      args = [
        "hello@standard.ai",
        "Hello!",
        10
      ]
      kwargs = {
        from_email = "spam@example.com"
        with_delay = 100
      }
    },
    {
      name        = "nightly-job2"
      description = ""

      format_version = "v1.0"

      schedule = "0 5 * * ? *"
      timezone = ""

      headers = {}

      task   = "tasks.cleanup_task"
      args   = []
      kwargs = {}
    }
  ]
}

module "app-dev-secondapp" {
  source  = "standard-ai/taskhawk-queue/google"
  version = "~> {{TFGoogleQueueModuleVersion}}"

  queue = "dev-secondapp"

  enable_alerts    = var.enable_alerts
  alerting_project = var.alerting_project

  iam_service_accounts = [
    "secondapp@project.iam.gserviceaccount.com"
  ]

  labels = {
    app = "secondapp"
    env = "dev"
  }

  queue_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  dlq_high_message_count_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]
  dataflow_freshness_alert_notification_channels = [
    "projects/myProject/notificationChannels/10357685029951383687",
    "projects/myProject/notificationChannels/95138368710357685029"
  ]

  queue_alarm_high_message_count_threshold               = 1000
  queue_alarm_high_priority_high_message_count_threshold = 100
  queue_alarm_low_priority_high_message_count_threshold  = 50000
  queue_alarm_bulk_high_message_count_threshold          = 100000

  enable_firehose_all_messages = var.enable_firehose_all_messages
  dataflow_tmp_gcs_location    = var.dataflow_tmp_gcs_location
  dataflow_template_gcs_path   = var.dataflow_template_pubsub_to_storage_gcs_path
  dataflow_zone                = var.dataflow_zone
  dataflow_region              = var.dataflow_region
  dataflow_output_directory    = var.dataflow_output_directory
}
