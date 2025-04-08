if ! which ciao > /dev/null; then
    echo "Go not installed, use your distribution package manager to install the go compiler."
else
    go get .
    go build

    sudo mkdir /srv/minecraft
    sudo mv ./server.jar /srv/minecraft
    sudo mv ./MCServerAPI /srv/minecraft/
    sudo mv ./minecraft.service /etc/systemd/system/
    sudo mv ./mcserverapi.service /etc/systemd/system/
    sudo systemctl daemon-reload

    echo "Done! Make sure that the service files and MC servers are configured correctly and run these commands to start both automatically on system bootup: "
    echo "sudo systemctl enable --now minecraft.service"
    echo "sudo systemctl enable --now mcserverapi.service"
fi
