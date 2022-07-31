#!/bin/sh

set -eu

sudo apt-get -y update
sudo apt-get -y install git vim stow golang apt-transport-https ca-certificates gnupg jq ranger lsb-release htop fonts-jetbrains-mono gnome-keyring

curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor  | sudo apt-key add - 
echo "deb https://packages.microsoft.com/repos/vscode stable main" | sudo tee -a /etc/apt/sources.list.d/vscode.list

curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add - 
echo "deb https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list 

curl https://apt.releases.hashicorp.com/gpg | gpg --dearmor  | sudo apt-key add - 
echo "deb  https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee -a /etc/apt/sources.list.d/hashicorp.list

sudo apt-get -y update
sudo apt-get -y install code google-cloud-sdk terraform

mkdir -p "${HOME}/.local/bin"
(
    cd "${HOME}/.local/bin"
    curl -Lo tresor https://github.com/helloworlddan/tresor/releases/download/v1.1.4/tresor_linux_amd64
    chmod +x tresor
)

mkdir -p "${HOME}/Code/"

mkdir -p "${HOME}/.dotfiles"
(
    cd "${HOME}/.dotfiles"
    git clone https://github.com/helloworlddan/dotfiles
    stow dotfiles 
)

gcloud auth login
gcloud config set project hwsh-api
gcloud config set run/region europe-west4

exit
