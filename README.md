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

Docker is also an option. In this case you'll need the built image, a `.env` file, and a MySQL or MariaDB database. Because it needs to run with `network_mode: host` (to be able to monitor the host network) the database server can be installed directly on the host or as an additional service placed in `docker-compose.yml` (which should also be `network_mode: host`).

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

Note: The `docker-compose.yml` file references `.env.prod` instead of `.env`.

Build the image and run it.

```bash
just up
```

## Access

The default credentials if `DEFAULT_ROOT_PASS` or `DEFAULT_ANON_PASS` is not set is `root:root` and `anon:anon`, respectively.

## Configuration

lighthouse is configured with a `.env` file. Only the options necessary for startup are configured via the environment, and the rest are configured with the web interface.

Check out `.env.example` for a quick reference or continue reading for a more comprehensive rundown.

### MODE

- Required: yes
- Value: `development` or `production`

The execution mode. Unless you're working on lighthouse, this should always be set to `production`.

### HOST

- Required: yes
- Value: an IP address e.g., `127.0.0.1` or `0.0.0.0`

The interface to bind to.

`0.0.0.0` binds to all interfaces allowing access for all hosts, while `127.0.0.1` binds to the loopback interface and restricts access to the local system. Other values will bind to that specific interface. The exact value to set here will depend on your network configuration.

### PORT

- Required: only if `MODE` is `production`
- Value: port, e.g., `1337`

The port to listen on. If not specified in development mode then a random port will be used.

### LOG_TYPE

- Required: yes
- Value: `all`, `disk`, `stderr`, or `none`

### LOG_PATH

- Required: only if `LOG_TYPE` is `all` or `disk`
- Value: path e.g., `logs/lighthouse.log`

### DB_NAME

- Required: yes
- Value: database name e.g., `lighthouse`

### DB_USERNAME

- Required: yes
- Value: username e.g., `anon`

### DB_PASSWORD

- Required: no
- Value: password e.g., `anon`
- Default: (empty)

### DB_PROTOCOL

- Required: no
- Value: `tcp` or `unix`
- Default: `tcp`

### DB_ADDRESS

- Required: no
- Value: host and port combination e.g., `127.0.0.1:3306`, or a unix socket path e.g. `/run/mysqld/mysqld.sock`
- Default: `127.0.0.1:3306`

### DB_CHARSET

- Required: no
- Value: character set e.g., `utf8mb4`
- Default: `utf8mb4`

### DB_TIMEZONE

- Required: no
- Value: timezone e.g., `UTC`, `Europe/London`
- Default: `UTC`

### LIVE_RELOAD

- Required: no
- Value: `true`, `1`, or `yes` to enable, anything else to disable
- Default: `false`

### DEFAULT_ROOT_PASS

- Required: no
- Value: password e.g., `root`
- Default: `root`

### DEFAULT_ANON_PASS

- Required: no
- Value: password e.g., `anon`
- Default: `anon`

### TRUSTED_PROXIES

- Required: no
- Value: comma-separated list of IP addresses e.g. `127.0.0.1`

If running behind a reverse proxy then set this to the proxy's IP address(es). This will ensure the correct IP address is displayed in the audit log.

## Mail

Mail is optional, but without it the only way to see notifications is to manually check using the web interface. Once up and running go to settings and configure the SMTP settings.

## Tests

Run tests with `just test`.

## Contributions

Open an [issue](https://github.com/crdx/lighthouse/issues) or send a [pull request](https://github.com/crdx/lighthouse/pulls).

## Licence

[GPLv3](LICENCE).
