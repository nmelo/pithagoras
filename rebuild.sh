#!/usr/bin/env bash

echo "Pulling git..." &&
cd ~/go/src/github.com/nmelo/pithagoras/ &&
git pull &&
cd ~ &&
echo "Done. Building server..." &&
go build github.com/nmelo/pithagoras/cmd/server &&
echo "Done. Restarting service..." &&
sudo service pithagoras stop &&
sudo cp ./server /usr/local/bin/pithagoras &&
sudo service pithagoras start &&
echo "Service started!"