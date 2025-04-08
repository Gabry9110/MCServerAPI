# MCServerAPI
REST API to start/stop a Java Minecraft Server from a systemd service. Written in Go using Gin, it also checks periodically for the current number of players to shut down the server automatically (similar to how Aternos works, but completely self-hosted). Useful if you already have an always on computer (like a home lab) and want to self host a MC server, without wasting resources if no one is playing.
Requires a Linux system with systemd.

## How to use
After running both the API server and the MC server (on the same system), there will be various endpoints usable from a web client (such as curl or a web browser). Default port is 8080 and the available endpoints are shown by Gin in the API server logs (also commented in `endpoints.go`).

## How to install
There is a bash script (`install.sh`) for convenient installation, but it still requires manual work:
- Downloading a `server.jar` file from Mojang
- Editing the `.service` files with the user(s)/group(s) that should execute both servers (remember that only `root` can open well-known ports)
- Changing ownership (recursively)
- Accept eula and/or edit server.proprieties after the `server.jar` file gets moved in `/srv/minecraft`

## TODOs and possible future features
- [ ] Possible authentication that allows only trusted parties to start/stop the server
- [ ] `/status` endpoint that returns more information
- [ ] Don't allow for the server to stop if someone is connected (online players check in `/stop`)
- [ ] Change the API port with a CLI parameter
- [ ] Config file to manually specify the MC server port/API server port/timeout/systemd service name
- [ ] More?
