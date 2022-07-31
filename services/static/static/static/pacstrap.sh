#!/bin/sh

set -eu

sudo pacman -Syyuu
sudo pacman -S --noconfirm git vim stow go gnupg jq ranger lsb-release htop gnome-keyring libsecret libgnome-keyring terraform

git clone https://aur.archlinux.org/yay.git
(
    cd yay
    makepkg -s
    sudo pacman -U yay*.tar.xz
)
rm -rf yay

yes | yay -S --noconfirm google-cloud-sdk visual-studio-code-bin 

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
