## Installing tinygo on the server (debian based)

echo "Installing TinyGo (AMD64)"
wget https://github.com/tinygo-org/tinygo/releases/download/v0.26.0/tinygo_0.26.0_amd64.deb
dpkg -i tinygo_0.26.0_amd64.deb

export PATH=$PATH:/usr/local/tinygo/bin

rm tinygo_0.26.0_amd64.deb

echo "Finished installing TinyGo: $(tinygo version)"
