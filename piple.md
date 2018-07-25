# Install gvm
$ bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)

# Install the latest go version (as of this writing)
$ gvm install go1.4

# Set it as active, and the default choice
$ gvm use go1.4 --default
