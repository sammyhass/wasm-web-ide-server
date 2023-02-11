curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - &&
	apt-get install -y nodejs &&
	node -v &&
	npm -v &&
	npm install -g assemblyscript &&
	asc --version
