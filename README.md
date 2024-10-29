# mctui-server

This is the backend for mctui. Check its repo for more details.

**Disclaimer**:
_Some information may be incomplete or incorrect._
_I will improve the documentation over time._

## Description

This is a webserver that helps managing a minecraft server. It uses RCON to send commands, and has it's own solution for managing backups.

This program should be used with [mctui](). You also can implement your frontend if you like.


## API

- `POST /login`: authenticate users

```json
{
    username="admin"
    password="1234"
}
```

**Response**: jwt token.

- `POST /command`: send comands to RCON

```json
{
    command="your-rcon-command"
}
```

**Response**: 200 OK or error.

- `POST /backup`: make a backup

**Response**: 200 OK or error.

- `GET /backups`: list all backups

**Response**: 200 OK or error.

```json
{
    [
        "backup-xxxyyyzzz.zip",
        "backup-xxxyyyaaa.zip",
        "backup-xxxyyybbb.zip",
    ]
}
```

- `POST /restore`: restore the selected backup

```json
{
    filename="backup-xxxyyyzzz.zip"
}
```

**Response**: 200 OK or error.

### CLI

- `add-user`: inserts an user in the database
- `list`: list all usernames

> Currently you can't remove users using the cli, but it's easilly done using the sqlite3 utility

## Setup

### Dependencies
- For compile yourself, install [go1.23.2](https://github.com/golang/go) or later
- Install [mcrcon](https://github.com/Tiiffi/mcrcon)
- Install [minecraft server](https://www.minecraft.net/en-us/download/server)
- Install [java](https://www.java.com/en/download/)
(must match you minecraft server version)
- Install [sqlite](https://www.sqlite.org/)
- You need systemd
- Go will download any other dependency
- If you use the backups script you also need [zip](https://infozip.sourceforge.net/Zip.html)

### Systemd

- Create the `minecraft.service` file in `~/.local/config/systemd/user/minecraft.service`. You may wanna check the example in `systemd/minecraft.service`.
- Reload the daemon `systemctl --user daemon-reload`
- Start the service `systemctl --user start minecraft`
#### Optional
There are 2 more services and 1 bash script inside the systemd directory: 

- `minecraft-backup.service`: service to create backups automatically
- `minecraft-backup.timer`: timer to trigger the service
- `minecraft-backup.sh`: example of a backup script

It is important that the backup files match the pattern
`"backup-2006-01-02-15-04-05.zip"`.

Keep that in mind if you create your own script.


### Minecraft server

- Create an empty direcory to use as you **backups** directory
- Create an empty direcory to use as you `minecraft-server` directory
- Run the server at least once. This will generate the `eula.txt` file
- Accept the eula
- Edit the [server.properties]()
    - Set the RCON `password` and `port`
    - Set `server-ip=0.0.0.0` if you want hosts outside the network to join
    (you probaly want)
    - Set `enable-rcon=true`

### HTTPS (certificates)
- Create the certificates ([this](https://stackoverflow.com/questions/10175812/how-to-generate-a-self-signed-ssl-certificate-using-openssl) may help)
- Put them somewhere the executable can read
- Add the `path` as environment variable (see bellow)


### Environment variables

```bash
# Secret used by jwt
export JWT_SECRET="you-secret-key"
# Database to store users
export DB_FILE="database.db"
# Certificate files fot HTTPS
export TLS_CERT_FILE="cert/cert.pem"
export TLS_KEY_FILE="cert/key.pem"
# Where your save resides. Should be named "world"
export WORLD_DIR="/home/$USER/tmp/minecraft-server/world"
# Where you backups are made. Should contain only backups
export BACKUP_DIR="/home/$USER/tmp/minecraft-backups"
# Used to comunicate with RCON. Must be the same on server.properties
export RCON_PASSWORD="minecraft"
export RCON_ADDRESS="127.0.0.1:25575"
```

- You need all variables listed here defined
- Using [direnv]() really simplify this process

### Adding Users

`mctui-server add-user --username="admin" --password="1234"`

You will need this credentials to loggin in `mctui`.


## Troubleshooting
- The first thing it to look the logs
- You can save the logs to `debug.log` with the environment `DEBUG=1` 

A common issue is **confliting options**. Sometimes, you use different paths on your scripts, systemd service, env variable, etc. I plan to automate this in the future, but for now, always double check them.
