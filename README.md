# RPiMonGo

Raspberry Pi Monitoring with Go

## 1-1. install

```bash
$ go get -d github.com/meinside/rpimongo
```

## 1-2. build

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ go build
```

## 2. setup for production

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

and run it with:

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ ./rpimongo
```

**NOTE**: It listens on port 80 and 443, so it needs to be run with root privilege.

## 3. access

Then it can be accessed through: `http://my.domain.com` (and also `https://my.domain.com` if you set "serve_ssl" as true)

## 4. run it as a service

### For systemd

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

