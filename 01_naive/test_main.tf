provider "aws" {
  access_key                  = "anaccesskey"
  secret_key                  = "asceretkey"
  region                      = "us-east-1"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  s3_force_path_style         = true

  endpoints {
    s3         = "http://localhost:4572"
    apigateway = "http://localhost:4567"
    lambda     = "http://localhost:4574"
  }
}

locals {
  talk_phase    = "go18talk-phase2-test"
  function_name = "toy"
  filepath      = "${format("%s/%s.zip", path.module, local.function_name)}"
}

module "api" {
  source          = "../terraform/modules/api"
  api_name        = "${local.talk_phase}"
  api_description = "first POC implementation of our API"
  environment     = "sandbox"

  # S3 Config
  bucket_name       = "${format("%s-bucket", local.talk_phase)}"
  enable_expiration = true

  # Lambda Config
  lambda_function_name     = "${local.function_name}"
  lambda_executable_name   = "${local.function_name}"
  lambda_function_filepath = "${local.filepath}"
  lambda_timeout           = 5

  lambda_env_vars = {
    SALT = "tVfTXzAfsTtS87PXFzALb64sKzTWm1dNMggdAZdWonE="
  }

  tags = {
    Type  = "Demo projects"
    Event = "Gophercon 2018 Kickoff Party"
  }
}

output "api_url" {
  value = "${module.api.api_url}"
}
