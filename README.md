# ercole-server
[![Build Status](https://travis-ci.org/ercole-io/ercole-server.svg?branch=master)](https://travis-ci.org/ercole-io/ercole-server) [![Join the chat at https://gitter.im/ercole-server/community](https://badges.gitter.im/ercole-server/community.svg)](https://gitter.im/ercole-server/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [![Docker Repository on Quay](https://quay.io/repository/amreo/ercole-server/status "Docker Repository on Quay")](https://quay.io/repository/amreo/ercole-server)

This is the server component for the Ercole project. Documentation available [here](https://ercole.netlify.com).

## Requirements

- Java Developement Kit 8.
- PostgreSQL server 9.6 or greater
- The ercole-web and ercole-packages should be available on your local or remote maven repository

## How to build

    mvn clean package

The resulting jar will be available under target/ercole-server-${version}.jar

## How to run

To run the binary:

    java -jar target/ercole-server-${version}.jar

You can customize the server parameters by copying the src/main/resources/application.properties file
in the same directory as your ercole-server jar file.
