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
  "port_number": 8088,
  "verbose": false
}
```

and run it with:

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ ./rpimongo
```

## 3. access

Then it can be accessed through: `http://my.raspberry.pi:8088`

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
User=some_user
Group=some_user
WorkingDirectory=/somewhere/go/src/github.com/meinside/rpimongo
ExecStart=/somewhere/go/src/github.com/meinside/rpimongo/rpimongo
Restart=always
RestartSec=5
Environment=

[Install]
WantedBy=multi-user.target
```

and edit **User**, **Group**, **WorkingDirectory**, and **ExecStart** values.

Make it autostart on every boot:

```bash
$ sudo systemctl enable rpimongo.service
```

and start it with:

```bash
$ sudo systemctl start rpimongo.service
```

## *(Optional)* 5. how to run it with Apache2 + reverse proxy

When used with apache2 and its reverse proxy, we can benefit from its functionalities like access logs.

### install apache2's proxy module and set it up

```bash
$ sudo apt-get install libapache2-mod-proxy-html
$ sudo a2enmod proxy_http
```

### create a proxy host

Create a site file:

```bash
$ sudo vi /etc/apache2/sites-enabled/some-host
```

**ProxyPass** and **ProxyPassReverse** should direct to your running RPiMonGo server:

```
<VirtualHost *:8080>
    ServerAdmin root@localhost
    ServerName my.raspberry.pi
    ProxyRequests Off
    <Proxy *>
        Order deny,allow
        Allow from all
    </Proxy>
    ProxyPass / http://127.0.0.1:8088/
    ProxyPassReverse / http://127.0.0.1:8088/
</VirtualHost>
```

Enable it and restart Apache2:

```bash
$ sudo a2ensite some-host
$ sudo service apache2 restart
```

then it can be accessed through: `http://my.raspberry.pi:8080`

