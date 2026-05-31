# lighthouse

**lighthouse** is a network monitor designed keep you aware of what's happening on your home network.

lighthouse periodically sends out [ARP](https://en.wikipedia.org/wiki/Address_Resolution_Protocol) requests to all hosts within a subnet and waits for responses. A responsive host is considered online, and a non-responsive host is considered offline after a (configurable) amount of time has passed. As some devices can't be relied on to reply to ARP requests in a timely fashion, if it looks like a device is about go offline an ICMP ping is sent to it to try to provoke a response, as a last-ditch effort.

lighthouse was born out of frustration with the [Fingbox](https://www.fing.com)'s increasing tendency to gate new features behind their premium offering while persistently upgrade-nagging whenever the app is opened.

## Features

- Track devices on the local network using ARP and ICMP
- Notify when devices go offline or come online
- Notify when devices stay online for longer than a specific amount of time
- Notify when an unknown device joins the network
- Service (port) scanning with notifications when new services are found
- Send ICMP packets to devices before they go offline to provoke a response as a last-ditch effort
- Assign icons to devices for easier identification
- Track DHCP-provided device hostname for easier identification
- User roles: admin, editor, viewer
- Support static mappings to work around dodgy devices like repeaters that rewrite ARP packets
- Passive mode which does not actively send out any packets
- Device vendor lookup using [gopacket's db](https://github.com/google/gopacket) with [macvendors](https://macvendors.com) fallback (requires API key)
- Support for devices with multiple adapters with different MAC addresses e.g. a laptop with ethernet & Wi-Fi
- Activity history, notification history, and audit log
- Minimal JavaScript (but not _none_)
- No remote dependencies
- Easy deployment with a single binary

## Deployment (binary)

The simplest way to run lighthouse is with the compiled binary, a `.env` file, and a MySQL or MariaDB database.

Permission is needed to send raw packets so it needs to run either as root (not recommended) or with capability `cap_net_raw`.

```bash
just make
sudo setcap cap_net_raw+ep dist/lighthouse
```

Set up the `.env` file (see [Configuration](#configuration)), then run with `--env`.

```bash
dist/lighthouse --env .env
```

## Deployment (docker)

Docker is also an option. In this case you'll need the built image, a `.env` file, and a MySQL or MariaDB database. Because it needs to run with `network_mode: host` (to be able to monitor the host network) the database server can be installed directly on the host or as an additional service placed in `compose.yml` (which should also be `network_mode: host`).

```yaml
  db:
    image: mariadb:latest
    restart: unless-stopped
    network_mode: host
    environment:
      MARIADB_USER: anon
      MARIADB_PASSWORD: anon
    volumes:
      - state:/var/lib/mysql
volumes:
  state:
```

Configure the `.env.prod` file to connect to it (see [Configuration](#configuration)), e.g.,

```
DB_USERNAME=anon
DB_PASSWORD=anon
DB_PROTOCOL=tcp
DB_ADDRESS=127.0.0.1:3306
```

Note: The `compose.yml` file references `.env.prod` instead of `.env`.

Build the image and run it.

```bash
just up
```

## Access

The default credentials if `DEFAULT_ROOT_PASS` or `DEFAULT_ANON_PASS` is not set is `root:root` and `anon:anon`, respectively.

## Configuration

lighthouse is configured with a `.env` file. Only the options necessary for startup are configured via the environment, and the rest are configured with the web interface.

Check out `.env.example` for a quick reference or continue reading for a more comprehensive rundown.

The first value (in **bold**) is the default.

| Variable             | Required      | Values                                            | Description                                                                                        |
|----------------------|---------------|---------------------------------------------------|----------------------------------------------------------------------------------------------------|
| `MODE`               | yes           | `development`, `production`                       | The execution mode. Should always be `production` unless developing lighthouse.                    |
| `HOST`               | yes           | `127.0.0.1`, `0.0.0.0`                           | The interface to bind to. `0.0.0.0` binds to all interfaces, `127.0.0.1` restricts to localhost.  |
| `PORT`               | in production | `1337`                                            | The port to listen on. A random port is used in development mode if not set.                       |
| `LOG_TYPE`           | yes           | `all`, `disk`, `stderr`, `none`                   | Where to send log output.                                                                          |
| `LOG_PATH`           | conditional   | `logs/lighthouse.log`                             | Log file path. Required if `LOG_TYPE` is `all` or `disk`.                                         |
| `DB_NAME`            | yes           | `lighthouse`                                      | The database name.                                                                                 |
| `DB_USERNAME`        | yes           | `anon`                                            | The database username.                                                                             |
| `DB_PASSWORD`        | no            | `anon`                                            | The database password. Defaults to empty.                                                          |
| `DB_PROTOCOL`        | no            | **`tcp`**, `unix`                                 | The database connection protocol.                                                                  |
| `DB_ADDRESS`         | no            | **`127.0.0.1:3306`**, `/run/mysqld/mysqld.sock`  | Host and port, or a unix socket path.                                                              |
| `DB_CHARSET`         | no            | **`utf8mb4`**                                     | The database character set.                                                                        |
| `DB_TIMEZONE`        | no            | **`UTC`**, `Europe/London`                        | The database timezone.                                                                             |
| `LIVE_RELOAD`        | no            | **`false`**, `true`                               | Enable live reload. Accepts `true`, `1`, or `yes`.                                                 |
| `DEFAULT_ROOT_PASS`  | no            | **`root`**                                        | Default password for the root user.                                                                |
| `DEFAULT_ANON_PASS`  | no            | **`anon`**                                        | Default password for the anon user.                                                                |
| `TRUSTED_PROXIES`    | no            | `127.0.0.1`                                       | Comma-separated IP addresses of trusted reverse proxies. Ensures correct IPs in the audit log.     |

## Mail

Mail is optional, but without it the only way to see notifications is to manually check using the web interface. Once up and running go to settings and configure the SMTP settings.

## Tests

Run tests with `just test`.

## Contributions

Open an [issue](https://github.com/crdx/lighthouse/issues) or send a [pull request](https://github.com/crdx/lighthouse/pulls).

## Licence

[GPLv3](LICENCE).
