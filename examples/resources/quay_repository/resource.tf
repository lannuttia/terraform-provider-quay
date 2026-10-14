resource "quay_repository" "example" {
  namespace   = "my-org"
  name        = "my-app"
  visibility  = "private"
  description = "Container images for my-app"
}
