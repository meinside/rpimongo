# RPiMonGo

A simple web server for monitoring your Raspberry Pi system, built with Go.

It shows `hostname`, `uptime`, `CPU temperature`, `free disk spaces`, `memory split`, and `free memory`.

![rpimongo](https://user-images.githubusercontent.com/185988/60861114-387a4300-a254-11e9-8d7a-0a8f146a9462.jpg)

If you own a public domain name, you can run in on HTTPS with the help of [Let's Encrypt](https://letsencrypt.org/).

## 1. Install

```bash
$ go install github.com/meinside/rpimongo@latest
```

## 2. Configure

Copy the sample config file and edit it:

```bash
$ git clone https://github.com/meinside/rpimongo.git
$ cp rpimongo/config.json.sample /path/to/config.json
$ vi /path/to/config.json
```

Example of **config.json**:

```json
{
  "hostname": "my.domain.com",
  "serve_ssl": false,
  "port_http": 9999,
  "verbose": false
}
```

## 3. Run

Run with the path to your config file:

```bash
$ rpimongo /path/to/config.json
```

NOTE: If it is run on port 80 and/or 443, it needs root privilige:

```bash
$ sudo rpimongo /path/to/config.json
```

## 4. Access

Then it can be accessed through: `http://my.domain.com` (and also `https://my.domain.com` if you set "serve_ssl" as true)

## 5. Run it as a service

### For systemd:

```bash
$ sudo vi /lib/systemd/system/rpimongo.service
```

```
[Unit]
Description=RPiMonGo
After=syslog.target
After=network.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/dir/to/config
ExecStart=/somewhere/rpimongo/is/installed/rpimongo config.json
Restart=always
RestartSec=5
Environment=

[Install]
WantedBy=multi-user.target
```

Edit **WorkingDirectory**, and **ExecStart** values to yours.

Make it autostart on every boot:

```bash
$ sudo systemctl enable rpimongo.service
```

and start it manually with:

```bash
$ sudo systemctl start rpimongo.service
```

