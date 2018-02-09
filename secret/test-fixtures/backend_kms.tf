provider "secret" {
  backend = "kms"
  config {
    aws_access_key = "123"
  }
}