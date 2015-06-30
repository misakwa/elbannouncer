ELB Announcer
=============

Simple sidekick for adding/removing instances to/from elb after balanced services.

```sh
$ elbannouncer link my-elb # links current instance to my-elb
$ elbannouncer unlink my-elb # unlinks current instance from my-elb
```

Sample fleetctl unit files
--------------------------
Sample apache service running on docker
```ini
[Unit]
Description=My Balanced Apache Frontend
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker kill apache1
ExecStartPre=-/usr/bin/docker rm apache1
ExecStartPre=/usr/bin/docker pull coreos/apache
ExecStart=/usr/bin/docker run -rm --name apache1 -p 80:80 coreos/apache /usr/sbin/apache2ctl -D FOREGROUND
ExecStop=/usr/bin/docker stop apache1

[X-Fleet]
Conflicts=my-balanced-service@*.service
```

Sample Unit file for sidekick

```ini
[Unit]
Description=ELB Announcer; Add service instance to elb
Documentation=https://github.com/misakwa/elbannouncer
BindsTo=my-balanced-service@%i.service
After=my-balanced-service@%i.service
Wants=my-balanced-service@%i.service
PartOf=my-balanced-service@%i.service

[Service]
TimeoutStartSec=0
ExecStart=/usr/bin/elbannouncer link apache-frontend
ExecStop=/usr/bin/elbannouncer unlink apache-frontend

[X-Fleet]
MachineOf=my-balanced-service@%i.service
```

```sh
$ fleetctl submit my-balanced-service@service
$ fleetctl start my-balanced-service@i
$ # Sidekick should start automatically after the service
```
