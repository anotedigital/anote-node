#!/bin/bash

# Get the files
wget https://github.com/anotedigital/anote-node/releases/download/v2.0-beta1/anote-node
wget https://github.com/anotedigital/anote-node/raw/main/conf/waves.conf
wget https://github.com/wavesplatform/Waves/releases/download/v1.4.11/waves_1.4.11_all.deb

# Stop waves
sudo service waves stop

# Remove old blockchain
sudo rm -rf /var/lib/anote/*

# Install new waves node
sudo dpkg -i waves_1.4.11_all.deb

# Configure new blockchain
chmod +x anote-node
mv seed secrets
source ./secrets
sed -i "s/D5u2FjJFcdit5di1fYy658ufnuzPULXRYG1YNVq68AH5/$ENCODED/g" waves.conf
sed -i "s/DTMZNMkjDzCwxNE1QLomcp5sXEQ9A3Mdb2RziN41BrYA/$KEYENCODED/g" waves.conf
sed -i "s/127.0.0.1:/$PUBLICIP:/g" waves.conf
mv waves.conf /etc/waves/waves.conf

# Stop waves
sudo service waves start