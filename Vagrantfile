# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "hashicorp/bionic64"

  config.vm.hostname = "beehive"
  config.vm.network "private_network", ip: "10.31.81.200"
  config.vm.network "forwarded_port", guest: 8086, host: 8086

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "2048"
  end

  config.vm.provision "shell", inline: <<-SCRIPT
    curl -sfL https://get.k3s.io | sh -
  SCRIPT
   
end
