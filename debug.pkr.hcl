source "googlecompute" "ex" {
  image_name              = "test-packer-example2"
  machine_type            = "e2-small"
  source_image            = "debian-10-buster-v20210316"
  ssh_username            = "packer"
  temporary_key_pair_type = "rsa"
  temporary_key_pair_bits = 2048
  zone                    = "us-central1-a"
  project_id              = "evaluation-blitz"
  running_instance_name   = "instance-6"
  use_os_login = true
  use_iap = true
}

build {
  sources = ["source.googlecompute.ex"]
  provisioner "shell" {
    inline = [
      "echo Hello From ${source.type} ${source.name}"
    ]
  }
}
