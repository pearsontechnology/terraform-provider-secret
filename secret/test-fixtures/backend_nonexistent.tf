provider "secret" {
  backend = "nonexistent"
  config {
    aws_access_key = "123"
  }
}