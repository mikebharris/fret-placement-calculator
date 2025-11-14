module "lambda" {
  source               = "./modules/lambda"
  region               = var.region
  contact              = var.contact
  project              = var.project
  orchestration        = var.orchestration
  distribution_bucket  = var.distribution_bucket
}