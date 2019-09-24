# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/bionic64"
  
	#Port de connexion sur SSH : 55555
	config.vm.network "forwarded_port", guest: 22, host: 55555
    config.vm.network "forwarded_port", guest: 5672, host: 5672
    config.vm.network "forwarded_port", guest: 5673, host: 5673
    config.vm.network "forwarded_port", guest: 15672, host: 15672
	config.vm.network "forwarded_port", guest: 10000, host: 10000
	config.vm.network "forwarded_port", guest: 10001, host: 10001
  
   config.vm.provider "virtualbox" do |vb|
     # Display the VirtualBox GUI when booting the machine
     vb.gui = false
  
     # Customize the amount of memory on the VM:
     vb.memory = "4096"
	 vb.cpus = 2
   end

  #config.vm.synced_folder "./", "/home/vagrant/sources", disabled: false

	config.vm.provision "shell", inline: <<-SHELL
		#
		# Configure hostfiles to use local squid proxy.
		#
	export DEBIAN_FRONTEND=noninteractive
    sudo apt-get update
    sudo apt-get -o DPkg::Options::=--force-confdef install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository \
    "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) \
    stable"
    sudo apt-get update
    sudo apt-get install -y docker-ce
    sudo usermod -aG docker vagrant
    sudo service docker restart
    sudo curl -L https://github.com/docker/compose/releases/download/1.18.0/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
    sudo chmod a+x /usr/local/bin/docker-compose
    sudo apt-get install -y putty-tools
    puttygen /vagrant/.vagrant/machines/default/virtualbox/private_key -O private -o /vagrant/.vagrant/machines/default/virtualbox/private_key.ppk
    echo "START putty.exe -P 55555 -i ".vagrant/machines/default/virtualbox/private_key.ppk" vagrant@127.0.0.1" > /vagrant/vm_login.bat
    cd /vagrant/registry
    sh setup.sh
    cd -
	SHELL
end