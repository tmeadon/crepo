#!/bin/bash

# Download binary to /usr/local/bin
curl -L https://github.com/tmeadon/crepo/releases/download/v0.1.6/crepo.tar.gz -o /tmp/crepo.tar.gz
tar -xzvf /tmp/crepo.tar.gz -C /usr/local/bin

# Print the wrapper function to .bashrc
/usr/local/bin/crepo --print >> ~/.bashrc

# Source the updated .bashrc
. ~/.bashrc