{{tfDoNotEditStamp}}

module "app-dev-anotherapp" {
  source        = "git@github.com:Automatic/taskhawk-terraform//lambda_app?ref={{VERSION}}"
  function_arn  = "arn:aws:lambda:us-west-2:12345:function:myFunction"
  function_name = "myFunction"
  name          = "dev-anotherapp"
}
