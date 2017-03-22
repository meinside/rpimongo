# RPiMonGo

Raspberry Pi Monitoring with Go

## 1-1. install

```bash
$ go get github.com/astaxie/beego
$ go get github.com/beego/bee
$ go get -d github.com/meinside/rpimongo
```

## 1-2. build

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ go build
```

## 2-1. setup for production

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ cp conf/app.conf.sample conf/app.conf
$ vi conf/app.conf
```

Example of **conf/app.conf**:

```
appname = My RPiMonGo Server
httpport = 8088
runmode = prod
```

and run it with:

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ bee run
```

## 2-2. setup for testing APIs (swagger)

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ cp conf/app.conf.sample conf/app.conf
$ vi conf/app.conf
```

Example of **conf/app.conf**:

```
appname = My RPiMonGo Server
httpaddr = 127.0.0.1
httpport = 8088
runmode = dev
EnableDocs = true
```

and run it with:

```bash
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ bee run -downdoc=true -gendoc=true
```

## 3. access

Then it can be accessed through: http://some.where:8088

or through: http://127.0.0.1:8088/swagger/swagger-1/

## 4. if you wanna run it as a service

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

then it can be accessed through: http://my.raspberry.pi:8080

