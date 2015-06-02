# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "ubuntu/trusty64"

  config.vm.provision "shell", inline: <<SCRIPT
set -x
set -u
set -e

apt-get -qq update
apt-get -qq install git-core  # build-essential curl libpcre3-dev pkg-config zip
git clone https://go.googlesource.com/go /opt/go
cd /opt/go
git checkout go1.4.2
cd src
./all.bash

mkdir -p /opt/gopath
cat <<EOF >/etc/profile.d/gopath.sh
export GOPATH="/opt/gopath"
export PATH="/opt/go/bin:/opt/gopath/bin:\$PATH"
EOF

mkdir -p /opt/gopath/src/github.com/rackspace
ln -s /vagrant /opt/gopath/src/github.com/rackspace/gophercloud

chown -R vagrant:vagrant /opt/go
chown -R vagrant:vagrant /opt/gopath

cd /opt/gopath/src/github.com/rackspace/gophercloud
go get
go build
go test
SCRIPT

  config.vm.provider "virtualbox" do |p|
    p.memory = 2048
    p.cpus = 2
  end
  
  ["vmware_fusion", "vmware_workstation"].each do |p|
    config.vm.provider "p" do |v|
      v.vmx["memsize"] = "2048"
      v.vmx["numvcpus"] = "2"
      v.vmx["cpuid.coresPerSocket"] = "1"
    end
  end
end
