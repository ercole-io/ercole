# ercole-server
[![Build Status](https://travis-ci.org/ercole-io/ercole-server.svg?branch=master)](https://travis-ci.org/ercole-io/ercole-server) [![Join the chat at https://gitter.im/ercole-server/community](https://badges.gitter.im/ercole-server/community.svg)](https://gitter.im/ercole-server/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

This is the server component for the Ercole project. Documentation available [here](https://ercole.io).

## Ercole server installation

### Requirements

|     Component     | Prerequisite                                 |
|:-----------------:|----------------------------------------------|
| Operating system  | Centos, RedHat, OracleLinux 7                |
| RAM               | 4GB                                          |
| Filesystem        | 50GB (minimum)                               |
| CPU               | 2 VirtualCPU                                 |
| Database          | PostgreSQL >= 9.6                            |
| Software          | java-11-openjdk                              |
| Network           | port 9080 (HTTP) or 443 (HTTPS) open         | 

### Installation 

* Postgresql db creation

```
$ psql
postgres-# create database ercole; 
postgres-# create user ercole with password 'ercole';    
postgres-# alter database ercole owner to ercole;
```

* Modify pg_hba.conf

```
vi <Postgresql data directory>/pg_hba.conf  <-- ex. /var/lib/pgsql/9.6/data/pg_hba.conf
```

```
# TYPE  DATABASE        USER            ADDRESS                 METHOD
# "local" is for Unix domain socket connections only
local   all             all                                     md5
# IPv4 local connections:
host    all             all             127.0.0.1/32            md5
# IPv6 local connections:
host    all             all             ::1/128                 ident
# Allow replication connections from localhost, by a user with the
# replication privilege.
local   replication     all                                     peer
host    replication     all             127.0.0.1/32            ident
host    replication     all             ::1/128                 ident
```

* OS user creation

```
useradd -s /bin/bash -g users -d /home/ercole -m ercole 
mkdir -p /opt/ercole-server/{log,conf} 
chown ercole.users /opt/ercole-server/log
```

* Install rpm Ercole Server 

```
yum install "rpm_ercole_server" (ex. ercole-server-1.5.0n-1.el7.x86_64.rpm)
```

* Configure and start Ercole Server

In order to configure ercole server you have to customize the file /opt/ercole-server/application.properties with the parameters different from the default.

Main parameter are:

| Parameter | Description | Default |
|----------------------------|------------------------------|-----------------------------------------|
| spring.datasource.url | Postgres database connection | jdbc:postgresql://localhost:5432/ercole |
| spring.datasource.username | DB user | ercole |
| spring.datasource.password | DB user password | ercole |
| user.normal.name | Ercole server user | user |
| user.normal.password | Ercole server user password | password |
| agent.user | Ercole agent user | user |
| agent.password | Ercole agent user password | password |
| agent.password | Ercole agent user password | password |
| server.port | Ercole server port | 9080 |

* systemctl daemon-reload
* systemctl start ercole.service
* systemctl enable ercole.service
