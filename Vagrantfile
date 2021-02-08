# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "hashicorp/bionic64"
  config.vm.hostname = "beehive"

  # expose rabbitmq to host
  config.vm.network :forwarded_port, guest: 30000, host: 5671
  config.vm.network :forwarded_port, guest: 30001, host: 15671

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "2048"
  end

  config.vm.provision "shell", inline: <<-SCRIPT
    curl -sfL https://get.k3s.io | sh -
  SCRIPT
   
end
