terraform {
  required_providers {
    discord = {
      source = "devcastops.com/stream/discord"
    }
  }
}

provider "discord" {}

data "discord_server" "this" {
  id = "1148301721953644624"
}