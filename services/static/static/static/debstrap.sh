#!/bin/sh

set -eu

sudo apt-get -y update
sudo apt-get -y install libnss3 libnspr4 git vim stow golang apt-transport-https ca-certificates gnupg jq ranger lsb-release htop fonts-jetbrains-mono gnome-keyring

curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add - 
echo "deb https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list 
sudo apt-get -y update
sudo apt-get -y install google-cloud-sdk

mkdir -p "${HOME}/.local/bin"
arch="$(uname -m)"
if [ "${arch}" = "x86_64" ]
then
    curl https://apt.releases.hashicorp.com/gpg | gpg --dearmor  | sudo apt-key add - 
    echo "deb  https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee -a /etc/apt/sources.list.d/hashicorp.list
    curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor  | sudo apt-key add - 
    echo "deb https://packages.microsoft.com/repos/vscode stable main" | sudo tee -a /etc/apt/sources.list.d/vscode.list
    sudo apt update
    sudo apt-get install -y code terraform
    (
        cd "${HOME}/.local/bin"
        curl -Lo tresor https://github.com/helloworlddan/tresor/releases/download/v1.1.4/tresor_linux_amd64
        chmod +x tresor
    )
elif [ "${arch}" = "aarch64" ]
then
    curl -L https://aka.ms/linux-arm64-deb > code_arm64.deb
    sudo apt install ./code_arm64.deb
    rm code_arm64.deb    
    (
        cd "${HOME}/.local/bin"
        curl -L https://releases.hashicorp.com/terraform/1.2.6/terraform_1.2.6_linux_arm64.zip | zcat > terraform
        chmod +x terraform
        curl -Lo tresor https://github.com/helloworlddan/tresor/releases/download/v1.1.4/tresor_linux_arm64
        chmod +x tresor
    )
else
    echo "unknown architecture\n"
fi

mkdir -p "${HOME}/Code/"

mkdir -p "${HOME}/.dotfiles"
(
    cd "${HOME}/.dotfiles"
    git clone https://github.com/helloworlddan/dotfiles
    stow dotfiles 
)
echo '\nsource "${HOME}/.bash_profile"\n' >> "${HOME}/.bashrc"

gcloud auth login
gcloud auth application-default login
gcloud config set project hwsh-api
gcloud config set run/region europe-west4

exit
