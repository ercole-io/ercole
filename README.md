# Ercole
[![Build Status](https://travis-ci.com/ercole-io/ercole.png)](https://travis-ci.com/ercole-io/ercole)
[![Gitter](https://badges.gitter.im/ercole-io/community.svg)](https://gitter.im/ercole-io/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![codecov](https://codecov.io/gh/amreo/ercole-services/branch/master/graph/badge.svg)](https://codecov.io/gh/ercole-io/ercole)
[![Go Report Card](https://goreportcard.com/badge/github.com/ercole-io/ercole)](https://goreportcard.com/report/github.com/ercole-io/ercole)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ercole-io_ercole&metric=alert_status)](https://sonarcloud.io/dashboard?id=ercole-io_ercole)
[![API Documentation](https://img.shields.io/badge/API%20Documentation-read%20and%20try-brightgreen)](https://mrin9.github.io/OpenAPI-Viewer/#/load/https%3A%2F%2Fraw.githubusercontent.com%2Fercole-io%2Fercole%2Fmaster%2Fswagger.yml)
 
Ercole is a open-source software for proactive software asset management:
 
Ercole is made by multiple services:
* Alert: generate alerts, send notifications. Expose data for 3rd party usage (i.e. prometheus)
* API: provides REST APIs for the User Interface
* Chart: provides REST APIs for the charts
* Data: receives data from the agent
* Repo: provides a yum repository (proxy?) for the agent binaries

Documentation about Ercole available [here](https://ercole.io).
Documentation about this new version of Ercole [here](https://ercole.io/architecture.html#future-versions)

# Main functionalities

**Licensing always under control** Take care about your Oracle Database installation and prevent the usage of unathorized licenses.

**Proactive database optimization** All interesting Oracle advisory output pre-elaborated and in a single point.

**RMAN Backup policy** Plan your RMAN backup policy in the best way.

**PSU and RU advisor** Plan your PSU and RU patching lifecycle.

**Database server CPU and storage capacity** Find your over allocated DB server and use the licenses where you really need them.

**Auto filling of LMS Oracle audit file** Have you ever tried to fill this complicated file? Ercole does it in one click.

## Requirements

- [Go](https://golang.org/)

## How to build

    go build ./main.go -o ercole

## How to run the server

Run the binary: `./ercole serve`

You can customize parameters by copying the `config.toml` file in the same directory as your ercole binary or in `/opt/ercole/config.toml`.
 
