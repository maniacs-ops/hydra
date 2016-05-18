# ![Ory/Hydra](dist/logo.png)

[![Join the chat at https://gitter.im/ory-am/hydra](https://badges.gitter.im/ory-am/hydra.svg)](https://gitter.im/ory-am/hydra?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/ory-am/hydra.svg?branch=master)](https://travis-ci.org/ory-am/hydra)
[![Coverage Status](https://coveralls.io/repos/ory-am/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/hydra?branch=master)



Hydra is being developed at [Ory](https://ory.am). Join our [mailinglist](http://eepurl.com/bKT3N9) to stay on top of new developments.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is Hydra?](#what-is-hydra)
- [Motivation](#motivation)
- [Quickstart](#quickstart)
  - [Installation](#installation)
    - [CLI Client](#cli-client)
    - [CLI Client using Docker (not recommended)](#cli-client-using-docker-not-recommended)
  - [Run minimal installation](#run-minimal-installation)
- [Documentation](#documentation)
  - [Guide](#guide)
  - [REST API Documentation](#rest-api-documentation)
  - [CLI Documentation](#cli-documentation)
  - [Develop](#develop)
- [Frequently Asked Questions](#frequently-asked-questions)
  - [Deploy using buildpacks (Heroku, Cloud Foundry, ...)](#deploy-using-buildpacks-heroku-cloud-foundry-)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Hydra?

1. Hydra is an OAuth2 and OpenID Connect provider built for availability. The distributed in-memory architecture allows for heavy duty workloads.
2. Hydra works with every Identity Provider. The deprecated php-3.0 authentication service your intern wrote? It works with that too, don't worry.
3. Hydra does not use any templates, it is up to you what your front end should look like.
4. Hydra comes with two factor authentication, key management, social log on, policy management and access control.

## Motivation

At first, there was the monolith. The monolith worked well with the customized joomla authentication module. Then, the web evolved into an elastic cloud that serves thousands of different user agents in every part of the world. Hydra is driven by the need for an easy scalable, in memory OAuth2 and OpenID Connect provider, that integrates with every Identity Provider you can imagine. 

Hydra uses pub/sub to always have the latest data available in memory. Hydra scales effortlessly on every platform you can imagine, including Heroku, Cloud Foundry, Docker, Google Container Engine and many more.

## Quickstart

This section is a quickstart guide to working with Hydra. In-depth docs are available as well:

* The documentation is available on [GitBook](https://ory-am.gitbooks.io/hydra/content/).
* The REST API documentation is available at [Apiary](http://docs.hdyra.apiary.io).

### Installation

**Starting the host** is easiest with docker. The host process handles HTTP requests and is backed by a database.
Read how to install docker on [Linux](https://docs.docker.com/linux/), [OSX](https://docs.docker.com/mac/) or
[Windows](https://docs.docker.com/windows/).

The easiest way to start docker is without a database. Hydra will keep all changes in memory. But be aware! Restarting, scaling
or stopping the container will make you **lose all data**.

```
$ docker run -d -p 4444:4444 oryam/hydra --name my-hydra
ec91228cb105db315553499c81918258f52cee9636ea2a4821bdb8226872f54b
```

**The CLI client** is available at [gobuild.io](https://gobuild.io/ory-am/hydra) or through
the [releases tab](https://github.com/ory-am/hydra/releases).

There is currently no installer which adds the client to your path automatically. You have to set up the path yourself.
If you do not understand what that means, ask on our [Gitter channel](https://gitter.im/ory-am/hydra).

If you wish to compile the CLI yourself, you need to install and set up [Go](https://golang.org/) and add `$GOPATH/bin`
to your `$PATH`. Here is a [comprehensive Go installation guide](https://github.com/ory-am/workshop-dbg#googles-go-language) with screenshots.

```
go install github.com/ory-am/hydra
hydra
```

Alternatively, you can use the CLI in Docker (not recommended):

```
$ docker exec -i -t <hydra-container-id> /bin/bash
# e.g. docker exec -i -t ec /bin/bash

root@ec91228cb105:/go/src/github.com/ory-am/hydra# hydra
Hydra is a twelve factor OAuth2 and OpenID Connect provider

Usage:
  hydra [command]

[...]
```

### Run the example

![Run the example](dist/run-the-example.gif)

Install the [CLI and Docker Toolbox](#installation). Make sure you install Docker Compose. On OSX and Windows,
open the Docker Quickstart Terminal. On Linux open any terminal.

**On OSX and Windows** using the Docker Quickstart Terminal.
```
$ go get github.com/ory-am/hydra
$ cd $GOPATH/src/github.com/ory-am/hydra
$ DOCKER_IP=$(docker-machine ip default) docker-compose up
Starting hydra_hydra
Starting hydra_consent
[...]
```

**On Linux.**
```
$ go get github.com/ory-am/hydra
$ cd $GOPATH/src/github.com/ory-am/hydra
$ DOCKER_IP=localhost docker-compose up
Starting hydra_hydra
Starting hydra_consent
[...]
mhydra   | mtime="2016-05-17T18:09:28Z" level=warning msg="Generated system secret: MnjFP5eLIr60h?hLI1h-!<4(TlWjAHX7"
[...]
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="Temporary root client created."
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_id: d9227bd5-5d47-4557-957d-2fd3bee11035"
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_secret: ,IvxGt02uNjv1ur9"
[...]
```

You have now a running hydra docker container! It is not backed by any database and runs completely in memory. Rebooting
or any other sort of disruption will purge all data.

*TBD: Provision with RethinkDB.*

Hydra can be managed with the hydra cli client. The client hast to log on before it is allowed to do anything.
When hydra detects a new installation, a new temporary root client is created. The client credentials are printed by
`docker compose up`:

```
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_id: d9227bd5-5d47-4557-957d-2fd3bee11035"
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_secret: ,IvxGt02uNjv1ur9"
```

The system secret is a global secret assigned to every hydra instance. It is used to encrypt data at rest. You can
set the system secret through the `$SYSTEM_SECRET` environment variable. When no secret is set, hydra generates one:

```
time="2016-05-15T14:56:34Z" level=warning msg="Generated system secret: (.UL_&77zy8/v9<sUsWLKxLwuld?.82B"
```

**Important note:** Please be aware that logging passwords should never be done on a production server. Either prune
the logs, set the required parameters, or replace the credentials with other ones.

Now you know which credentials you need to use. Next, we log in.

**Note:** If you are using docker toolbox, please use the ip address provided by `docker-machine ip default` as cluster url host.

```
$ hydra connect
Cluster URL: https://localhost:4444
Client ID: d9227bd5-5d47-4557-957d-2fd3bee11035
Client Secret: ,IvxGt02uNjv1ur9
Done.
```

Great! You are now connected to Hydra and can start by creating a new client:

```
$ hydra clients create --skip-tls-verify
Warning: Skipping TLS Certificate Verification.
Client ID: c003830f-a090-4721-9463-92424270ce91
Client Secret: Z2pJ0>Tp7.ggn>EE&rhnOzdt1
```

**Important note:** Hydra is using self signed TLS certificates for HTTPS, if no certificate was provided. This should
never be done in production. To skip the TLS verification step on the client, provide the `--skip-tls-verify` flag.

Why not issue an access token for your client?

```
$ hydra token client --skip-tls-verify
Warning: Skipping TLS Certificate Verification.
JLbnRS9GQmzUBT4x7ESNw0kj2wc0ffbMwOv3QQZW4eI.qkP-IQXn6guoFew8TvaMFUD-SnAyT8GmWuqGi3wuWXg
```

Let's try this with the authorize code grant!

```
$ hydra token user --skip-tls-verify
Warning: Skipping TLS Certificate Verification.
If your browser does not open automatically, navigate to: https://192.168.99.100:4444/oauth2/auth?client_id=d9227bd5-5d47-4557-957d-2fd3bee11035&response_type=code&scope=core+hydra&state=sbnwdelqzxyedwtqinxnolbr&nonce=sffievieeesltbjkwxyhycyq
Setting up callback listener on http://localhost:4445/callback
Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.
```

![OAuth2 Flow](dist/oauth2-flow.gif)

Great! You installed hydra, connected the CLI, created a client and completed two authentication flows!
Your next stop should be the [Guide](#guide).

## Documentation

### Guide

The Guide is available on [GitBook](https://ory-am.gitbooks.io/hydra/content/).

### REST API Documentation

The REST API is documented at [Apiary](http://docs.hdyra.apiary.io).

### CLI Documentation

The CLI help is verbose. To see it, run `hydra -h` or `hydra [command] -h`.

### Develop

Unless you want to test Hydra against a database, developing with Hydra is as easy as

```
go get github.com/ory-am/hydra
cd $GOPATH/src/github.com/ory-am/hydra
git checkout -b develop
go test ./... -race
go run main.go
```

## Frequently Asked Questions

### Deploy using buildpacks (Heroku, Cloud Foundry, ...)

Hydra runs pretty much out of the box when using a Platform as a Service (PaaS).
Here are however a few notes which might assist you in your task:
* Heroku (and probably Cloud Foundry as well) *force* TLS termination, meaning that Hydra must be configured with `DANGEROUSLY_FORCE_HTTP=force`.
* Using bash, you can easily add multi-line environment variables to Heroku using `heroku config:set JWT_PUBLIC_KEY="$(my-public-key.pem)"`.
  This does not work on Windows!
