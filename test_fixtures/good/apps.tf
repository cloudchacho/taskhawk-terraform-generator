{{tfDoNotEditStamp}}

module "app-dev-myapp" {
  source = "git@github.com:Automatic/taskhawk-terraform//queue_app?ref={{VERSION}}"
  queue  = "DEV-MYAPP"
  iam    = "true"

  tags = {
    App = "myapp"
    Env = "dev"
  }
}

module "app-dev-secondapp" {
  source = "git@github.com:Automatic/taskhawk-terraform//queue_app?ref={{VERSION}}"
  queue  = "DEV-SECONDAPP"
  iam    = "true"

  tags = {
    App = "secondapp"
    Env = "dev"
  }
}

module "app-dev-anotherapp" {
  source             = "git@github.com:Automatic/taskhawk-terraform//lambda_app?ref={{VERSION}}"
  function_arn       = "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed"
  function_name      = "myFunction"
  function_qualifier = "deployed"
  name               = "dev-anotherapp"
}
