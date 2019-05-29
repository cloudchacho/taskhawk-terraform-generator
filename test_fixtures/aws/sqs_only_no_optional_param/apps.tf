{{tfDoNotEditStamp}}

module "app-dev-myapp" {
  source  = "Automatic/taskhawk-queue/aws"
  version = "~> {{TFAWSQueueModuleVersion}}"

  queue = "DEV-MYAPP"

  tags = {
    App = "myapp"
    Env = "dev"
  }
}
