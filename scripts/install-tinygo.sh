## Installing tinygo on the server (debian based)

echo "Install TinyGo (AMD64)"
wget https://github.com/tinygo-org/tinygo/releases/download/v0.26.0/tinygo_0.26.0_amd64.deb
dpkg -i tinygo_0.26.0_amd64.deb
export PATH=$PATH:/usr/local/tinygo/bin
echo "Finished installing TinyGo: $(tinygo version)"
