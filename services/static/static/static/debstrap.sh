#!/bin/sh

set -eu

sudo sed -i 's/bullseye/testing/g' /etc/apt/sources.list
sudo apt-get -y update
sudo apt-get -y dist-upgrade

(
    curl https://dl.google.com/linux/linux_signing_key.pub | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/google-chrome.gpg
    echo "deb [arch=$(dpkg --print-architecture)] https://dl.google.com/linux/chrome/deb/ stable main" | sudo tee /etc/apt/sources.list.d/google-chrome.list
)
(
    curl https://dl.google.com/linux/linux_signing_key.pub | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/chrome-remote-desktop.gpg
    echo "deb [arch=$(dpkg --print-architecture)] https://dl.google.com/linux/chrome-remote-desktop/deb/ stable main" | sudo tee /etc/apt/sources.list.d/chrome-remote-desktop.list
)
(
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo tee /etc/apt/trusted.gpg.d/google-cloud-sdk.gpg
    echo "deb https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list 
)
(
    curl https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo tee /etc/apt/trusted.gpg.d/github-cli.gpg
    echo "deb [arch=$(dpkg --print-architecture)] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list
)
(
    curl https://repo.charm.sh/apt/gpg.key | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/charm.gpg
    echo "deb [arch=$(dpkg --print-architecture)] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
)

sudo apt-get -y update
sudo apt-get -y install 

sudo apt-get -y install \
    apt-transport-https \
    bspwm \
    ca-certificates \
    fonts-jetbrains-mono \
    git \
    gnome-dust-icon-theme \
    gnome-keyring \
    gnupg \
    golang \
    htop \
    jq \
    libnspr4 \
    libnss3 \
    lsb-release \
    polybar \
    ranger \
    rofi \
    sox \
    stow \
    sxhkd \
    tree \
    vim \
    thunar \
    kitty \
    xserver-xephyr \
    chrome-remote-desktop \
    google-chrome-stable \
    google-cloud-sdk \
    gh \
    gum \
    vhs

mkdir -p "${HOME}/.local/bin"
arch="$(uname -m)"
if [ "${arch}" = "x86_64" ]
then
    curl https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/hashicorp.gpg
    echo "deb  https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee -a /etc/apt/sources.list.d/hashicorp.list
    curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor  | sudo tee /etc/apt/trusted.gpg.d/vscode.gpg
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

mkdir -p "${HOME}/.dotfiles"
(
    cd "${HOME}/.dotfiles"
    git clone https://github.com/helloworlddan/dotfiles
    stow dotfiles 
)
echo ''  >> "${HOME}/.bashrc"
echo 'source "${HOME}/.bash_profile"' >> "${HOME}/.bashrc"
echo ''  >> "${HOME}/.bashrc"

mkdir -p "${HOME}/.go/"
mkdir -p "${HOME}/Code/"

bash

go install -v golang.org/x/tools/gopls@latest
go install -v github.com/go-delve/delve/cmd/dlv@latest
go install -v github.com/mdempsky/gocode@latest
go install -v github.com/stamblerre/gocode@latest
go install -v github.com/ramya-rao-a/go-outline@latest
go install -v github.com/acroca/go-symbols@latest
go install -v golang.org/x/tools/cmd/guru@latest
go install -v golang.org/x/tools/cmd/gorename@latest
go install -v github.com/fatih/gomodifytags@latest
go install -v github.com/haya14busa/goplay/cmd/goplay@latest
go install -v github.com/josharian/impl@latest
go install -v github.com/tylerb/gotype-live@latest
go install -v github.com/rogpeppe/godef@latest
go install -v github.com/zmb3/gogetdoc@latest
go install -v mvdan.cc/gofumpt/gofumports@latest
go install -v mvdan.cc/gofumpt@latest
go install -v golang.org/x/tools/cmd/goimports@latest
go install -v github.com/sqs/goreturns@latest
go install -v winterdrache.de/goformat/goformat@latest
go install -v github.com/cweill/gotests/gotests@latest
go install -v golang.org/x/lint/golint@latest
go install -v honnef.co/go/tools/cmd/staticcheck@latest
go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install -v github.com/mgechev/revive@latest
go install -v github.com/mgechev/revive@latest
go install -v github.com/godoctor/godoctor@latest

sudo apt-get autoremove -y

gcloud auth login
gcloud auth application-default login
gcloud config set project hwsh-api
gcloud config set run/region europe-west4
gcloud config set deploy/region europe-west4
gcloud config set compute/region europe-west4
gcloud config set artifacts/location europe-west4

exit
