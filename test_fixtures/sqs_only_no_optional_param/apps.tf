{{tfDoNotEditStamp}}

module "app-dev-myapp" {
  source = "git@github.com:Automatic/taskhawk-terraform//queue_app?ref={{VERSION}}"
  queue  = "DEV-MYAPP"

  tags = {
    App = "myapp"
    Env = "dev"
  }
}
