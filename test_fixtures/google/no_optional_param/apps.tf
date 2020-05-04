{{tfDoNotEditStamp}}

module "app-dev-myapp" {
  source  = "standard-ai/taskhawk-queue/google"
  version = "~> {{TFGoogleQueueModuleVersion}}"

  queue = "dev-myapp"

  labels = {
    app = "myapp"
    env = "dev"
  }

  enable_firehose_all_messages = var.enable_firehose_all_messages
}

module "app-dev-secondapp" {
  source  = "standard-ai/taskhawk-queue/google"
  version = "~> {{TFGoogleQueueModuleVersion}}"

  queue = "dev-secondapp"

  labels = {
    app = "secondapp"
    env = "dev"
  }

  enable_firehose_all_messages = var.enable_firehose_all_messages
}
