{{tfDoNotEditStamp}}

module "app-dev-anotherapp" {
  source  = "Automatic/taskhawk-lambda/aws"
  version = "~> {{TFLambdaModuleVersion}}"

  name          = "dev-anotherapp"
  function_arn  = "arn:aws:lambda:us-west-2:12345:function:myFunction"
  function_name = "myFunction"
}
