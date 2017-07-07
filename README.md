# RPiMonGo

A simple web server for monitoring your Raspberry Pi system, built with Go.

It shows `hostname`, `uptime`, `CPU temperature`, `free disk spaces`, `memory split`, and `free memory`.

If you own a public domain name, you can run in on HTTPS with the help of [Let's Encrypt](https://letsencrypt.org/).

## 1. Install

```bash
$ go get -d github.com/meinside/rpimongo
```

## 2. Build

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ go build
```

## 3. Setup

Copy the sample config file and edit it:

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ cp config.json.sample config.json
$ vi config.json
```

Example of **config.json**:

```json
{
  "hostname": "my.domain.com",
  "serve_ssl": true,
  "verbose": false
}
```

## 4. Run

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ ./rpimongo
```

**NOTE**: It listens on port 80 and 443, so it needs to be run with root privilege.

## 5. Access

Then it can be accessed through: `http://my.domain.com` (and also `https://my.domain.com` if you set "serve_ssl" as true)

## 6. Run it as a service

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
WorkingDirectory=/somewhere/go/src/github.com/meinside/rpimongo
ExecStart=/somewhere/go/src/github.com/meinside/rpimongo/rpimongo
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

