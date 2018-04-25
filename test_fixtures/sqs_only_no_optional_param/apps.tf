{{tfDoNotEditStamp}}

module "app-dev-myapp" {
  source  = "Automatic/taskhawk-queue/aws"
  version = "~> {{TFQueueModuleVersion}}"

  queue = "DEV-MYAPP"

  tags = {
    App = "myapp"
    Env = "dev"
  }
}
