# RPiMonGo

Raspberry Pi Monitoring with Go

## 1. install

```
$ go get github.com/astaxie/beego
$ go get github.com/beego/bee
$ go get github.com/meinside/rpimongo
```

## 2-1. setup for production

```
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

```
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ bee run
```

## 2-2. setup for testing APIs (swagger)

```
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

```
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ bee run -downdoc=true -gendoc=true
```

## 3. access

Then it can be accessed through: http://some.where:8088

or through: http://127.0.0.1:8088/swagger/swagger-1/

## 4. if you wanna run it as a service

### For init.d

#### create init.d script

```
$ sudo vi /etc/init.d/rpimongo
```

Edit **RPIMONGO_DIR** to yours:

```
#!/bin/sh
### BEGIN INIT INFO
# Provides:          rpimongo
# Required-Start:    networking
# Required-Stop:     networking
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: RPiMonGo init.d script
# Description:       
# 
#  RPiMonGo(Raspberry Pi Monitoring with Go) init.d script
# 
#  last update: 2015.08.07.
# 
#  meinside@gmail.com
#                    
### END INIT INFO

# change following path
RPIMONGO_DIR=/somewhere/go/src/github.com/meinside/rpimongo

NAME=rpimongo
DAEMON=$RPIMONGO_DIR/$NAME
DESC="RPiMonGo service"

# exit if not executable
test -x $DAEMON || exit 0
  
case "$1" in
  start)
    echo -n "Starting $DESC: "
    start-stop-daemon --start --quiet --background --make-pidfile --pidfile /var/run/$NAME.pid --exec $DAEMON || true
    echo "$NAME."
    ;;
  stop)
    echo -n "Stopping $DESC: "
    start-stop-daemon --stop --quiet --pidfile /var/run/$NAME.pid || true
    echo "$NAME."
    ;;
  restart)
    echo -n "Restarting $DESC: "
    start-stop-daemon --stop --quiet --pidfile /var/run/$NAME.pid || true
    sleep 1
    start-stop-daemon --start --quiet --background --make-pidfile --pidfile /var/run/$NAME.pid --exec $DAEMON || true
    echo "$NAME."
    ;;
  *)
    N=/etc/init.d/$NAME
    echo "Usage: $N {start|stop|restart}" >&2
    exit 1
    ;;
esac

exit 0
```

#### setup & run

Run it:

`$ sudo service rpimongo start`

restart it:

`$ sudo service rpimongo restart`

or stop it:

`$ sudo service rpimongo stop`

If you want it to launch on boot time:

`$ sudo update-rc.d -f rpimongo defaults`

### For systemd

```
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

```
$ sudo systemctl enable rpimongo.service
```

and start it with:

```
$ sudo systemctl start rpimongo.service
```

## *(Optional)* 5. how to run it with Apache2 + reverse proxy

When used with apache2 and its reverse proxy, we can benefit from its functionalities like access logs.

### install apache2's proxy module and set it up

```
$ sudo apt-get install libapache2-mod-proxy-html
$ sudo a2enmod proxy_http
```

### create a proxy host

Create a site file:

```
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

```
$ sudo a2ensite some-host
$ sudo service apache2 restart
```

then it can be accessed through: http://my.raspberry.pi:8080

