#!/bin/sh

echo "btw"

set -eu

sudo pacman -Syyuu
sudo pacman -S --noconfirm \
    bspwm \
    git \
    github-cli \
    gnome-keyring \
    gnupg \
    go \
    htop \
    jq \
    kitty \
    libgnome-keyring \
    libsecret \
    lsb-release \
    nitrogen \
    polybar \
    ranger \
    rofi  \
    stow \
    sxhkd \
    terraform \
    ttf-jetbrains-mono \
    vim 

git clone https://aur.archlinux.org/yay.git
(
    cd yay
    makepkg -s
    sudo pacman -U yay*.tar.xz
)
rm -rf yay

yes | sudo yay -S --noconfirm google-chrome google-cloud-sdk visual-studio-code-bin chrome-remote-desktop  vhs-bin  gum

mkdir -p "${HOME}/.local/bin"
(
    cd "${HOME}/.local/bin"
    if [ $(uname -m) == "x86_64" ]
    then
        curl -Lo tresor https://github.com/helloworlddan/tresor/releases/download/v1.1.4/tresor_linux_amd64
    else
        curl -Lo tresor https://github.com/helloworlddan/tresor/releases/download/v1.1.4/tresor_linux_arm64
    fi
    chmod +x tresor
)

mkdir -p "${HOME}/.go/"
mkdir -p "${HOME}/Code/"

mkdir -p "${HOME}/.dotfiles"
(
    cd "${HOME}/.dotfiles"
    git clone https://github.com/helloworlddan/dotfiles
    stow dotfiles 
)
echo '\nsource "${HOME}/.bash_profile"\n' >> "${HOME}/.bashrc"

yes | sudo yay -Scc --noconfirm

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

gcloud auth login
gcloud auth application-default login
gcloud config set project hwsh-api
gcloud config set run/region europe-west4
gcloud config set deploy/region europe-west4
gcloud config set compute/region europe-west4
gcloud config set artifacts/location europe-west4

exit
