# lighthouse

**lighthouse** is a network monitor designed keep you aware of what's happening on your home network.

In a nutshell, lighthouse periodically sends out [ARP](https://en.wikipedia.org/wiki/Address_Resolution_Protocol) requests to all hosts within a subnet and waits for responses. A responsive host is considered online, and a non-responsive host is considered offline after a (configurable) amount of time has passed.

Devices can be configured so a notification is triggered when they go offline or come online, and the whole network can be watched so notifications are sent when a new device joins the network.

lighthouse was born out of frustration with the [Fingbox](https://www.fing.com)'s increasing tendency to gate new features behind their premium offering while consistently upgrade-nagging whenever the app is opened.

## Deployment (binary)

The simplest way to get lighthouse running is with the compiled binary, a `.env` file, and a local MySQL or MariaDB database.

Permission is needed to send raw packets so lighthouse needs to run either as root (not recommended) or with capability `cap_net_raw`.

```bash
just make
sudo setcap cap_net_raw+ep dist/lighthouse
```

Set up the `.env` file (see [Configuration](#configuration)), then run lighthouse with `--env`.

```bash
dist/lighthouse --env .env
```

## Deployment (docker)

Deployment with docker is also an option. In this case you'll need the built image, a `.env` file, and a local MySQL or MariaDB database. Because lighthouse needs to run with `network_mode: host` (to be able to monitor the host network) the local database server can be installed directly on the host or as an additional service placed in `docker-compose.yml` (which should also be `network_mode: host`).

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

Configure the `.env.prod` file to connect to it (see [Configuration](#configuration)).

```
DB_USER=anon
DB_PASS=anon
DB_HOST=127.0.0.1:3306
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

`0.0.0.0` binds to all interfaces making lighthouse accessible to all hosts, while `127.0.0.1` will bind to the loopback interface and restrict it to the local system. Other values will bind to that specific interface.

The exact value to set here will depend on your network configuration.

### PORT

- Required: only if `MODE` is `production`
- Value: port, e.g., `1337`

The port to listen on. If not specified in development mode then lighthouse will listen on a random port.

### LOG_TYPE

- Required: yes
- Value: `all`, `disk`, `stderr`, or `none`

### LOG_PATH

- Required: only if `LOG_TYPE` is `all` or `disk`
- Value: path e.g., `logs/lighthouse.log`

### DB_NAME

- Required: yes
- Value: database name e.g., `lighthouse`

### DB_USER

- Required: yes
- Value: username e.g., `anon`

### DB_PASS

- Required: no
- Value: password e.g., `anon`

### DB_SOCK

- Required: only if `DB_HOST` is not set
- Value: path to a unix socket e.g., `/run/mysqld/mysqld.sock`

### DB_HOST

- Required: only if `DB_SOCK` is not set
- Value: host and port combination e.g., `127.0.0.1:3306`

### DB_CHARSET

- Required: yes
- Value: character set e.g., `utf8mb4`

### DB_TZ

- Required: yes
- Value: timezone e.g., `UTC`, `Europe/London`

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

Mail is optional, but without it the only way to see notifications is to manually check using the web interface. Once lighthouse is up and running go to settings and configure the SMTP settings.

## Tests

Run tests with `just test`.

## Contributions

Open an [issue](https://github.com/crdx/lighthouse/issues) or send a [pull request](https://github.com/crdx/lighthouse/pulls).

## Licence

[GPLv3](LICENCE).
