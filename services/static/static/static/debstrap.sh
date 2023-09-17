#!/bin/sh

set -eu

sudo sed -i 's/bullseye/testing/g' /etc/apt/sources.list # Somehow whacks sudo/root on crostini
sudo apt-get -y update
sudo apt-get -y dist-upgrade

curl https://dl.google.com/linux/linux_signing_key.pub             | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/google-chrome.gpg         >/dev/null
curl https://dl.google.com/linux/linux_signing_key.pub             | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/chrome-remote-desktop.gpg >/dev/null
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg         | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/google-cloud-sdk.gpg      >/dev/null
curl https://cli.github.com/packages/githubcli-archive-keyring.gpg                 | sudo tee /etc/apt/trusted.gpg.d/github-cli.gpg            >/dev/null
curl https://repo.charm.sh/apt/gpg.key                             | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/charm.gpg                 >/dev/null
curl https://apt.releases.hashicorp.com/gpg                        | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/hashicorp.gpg             >/dev/null
curl https://packages.microsoft.com/keys/microsoft.asc             | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/vscode.gpg                >/dev/null

echo "deb [arch=$(dpkg --print-architecture)] https://dl.google.com/linux/chrome/deb/ stable main"                | sudo tee /etc/apt/sources.list.d/google-chrome.list         >/dev/null
echo "deb [arch=$(dpkg --print-architecture)] https://dl.google.com/linux/chrome-remote-desktop/deb/ stable main" | sudo tee /etc/apt/sources.list.d/chrome-remote-desktop.list >/dev/null
echo "deb https://packages.cloud.google.com/apt cloud-sdk main"                                                   | sudo tee /etc/apt/sources.list.d/google-cloud-sdk.list      >/dev/null
echo "deb [arch=$(dpkg --print-architecture)] https://cli.github.com/packages stable main"                        | sudo tee /etc/apt/sources.list.d/github-cli.list            >/dev/null
echo "deb [arch=$(dpkg --print-architecture)] https://repo.charm.sh/apt/ * *"                                     | sudo tee /etc/apt/sources.list.d/charm.list                 >/dev/null
echo "deb  https://apt.releases.hashicorp.com $(lsb_release -cs) main"                                            | sudo tee /etc/apt/sources.list.d/hashicorp.list             >/dev/null
echo "deb https://packages.microsoft.com/repos/vscode stable main"                                                | sudo tee /etc/apt/sources.list.d/vscode.list                >/dev/null

sudo apt-get -y update

sudo apt-get -y install \
    cros-adapta \
    cros-apt-config \
    cros-garcon \
    cros-host-fonts \
    cros-logging \
    cros-notificationd \
    cros-pulse-config \
    cros-sftp \
    cros-sommelier \
    cros-sommelier-config \
    cros-sudo-config \
    cros-systemd-overrides \
    cros-ui-config \
    cros-wayland

sudo apt-get -y install \
    build-essential \
    apt-transport-https \
    ca-certificates \
    fonts-jetbrains-mono \
    git \
    gnome-keyring \
    gnupg \
    golang \
    htop \
    jq \
    libnspr4 \
    libnss3 \
    lsb-release \
    ranger \
    sox \
    stow \
    tree \
    vim \
    google-cloud-sdk \
    gh \
    gum \
    ruby-rubygems \
    ruby-bundler \
    ruby-dev \
    terraform \
    npm \
    libfuse2 \
    vhs

sudo apt-get -y install \
    gnome-dust-icon-theme \
    gnome-noble-icon-theme \
    rofi \
    xterm \
    polybar \
    bspwm \
    sxhkd \
    xserver-xephyr \
    thunar \
    code \
    chrome-remote-desktop \
    google-chrome-stable 

mkdir -p "${HOME}/.local/bin"
curl -L tresor https://github.com/helloworlddan/tresor/releases/download/v1.1.4/tresor_linux_amd64 > "${HOME}/.local/bin/tresor"
chmod +x "${HOME}/.local/bin/tresor"

mkdir -p "${HOME}/.dotfiles"
(
    cd "${HOME}/.dotfiles"
    git clone https://github.com/helloworlddan/dotfiles
    stow dotfiles 
)

echo 'source "${HOME}/.bash_profile"' >> "${HOME}/.bashrc"

mkdir -p "${HOME}/Code/"

mkdir -p "${HOME}/.go/"
export GOPATH="${HOME}/.go/"

go install -v github.com/helloworlddan/tortune@latest
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
