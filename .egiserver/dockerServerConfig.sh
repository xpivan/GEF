#!/bin/bash
sudo chmod 0400 $HOME/.docker/server-key.pem
sudo service docker stop
sudo dockerd -D --tls=true --tlscert=$HOME/.docker/server.pem --tlskey=$HOME/.docker/server-key.pem -H tcp://$1:2376
sudo sudo groupadd docker
sudo gpasswd -a $USER docker 
