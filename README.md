# RPiMonGo

Raspberry Pi Monitoring with Go

## install

```
$ go get github.com/astaxie/beego
$ go get github.com/beego/bee
$ go get github.com/meinside/rpimongo
```

## setup

```
$ cd $GOPATH/src/github.com/meinside/rpimongo
$ cp conf/app.conf.sample conf/app.conf
$ vi conf/app.conf
```

## run

```
$ bee run
```

## run as service

`$ sudo vi /etc/init.d/rpimongo`

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

`$ sudo update-rc.d -f rpimongo defaults`

