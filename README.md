# drone-facebook

[![GoDoc](https://godoc.org/github.com/appleboy/drone-facebook?status.svg)](https://godoc.org/github.com/appleboy/drone-facebook)
[![Build Status](https://cloud.drone.io/api/badges/appleboy/drone-facebook/status.svg)](https://cloud.drone.io/appleboy/drone-facebook)
[![Build status](https://ci.appveyor.com/api/projects/status/aexij85gjg3dsesl?svg=true)](https://ci.appveyor.com/project/appleboy/drone-facebook)
[![codecov](https://codecov.io/gh/appleboy/drone-facebook/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-facebook)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-facebook)](https://goreportcard.com/report/github.com/appleboy/drone-facebook)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-facebook.svg)](https://hub.docker.com/r/appleboy/drone-facebook/)
[![](https://images.microbadger.com/badges/image/appleboy/drone-facebook.svg)](https://microbadger.com/images/appleboy/drone-facebook "Get your own image badge on microbadger.com")

[Drone](https://github.com/drone/drone) plugin for sending [Facebook Messages](https://developers.facebook.com/docs/messenger-platform). For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/appleboy/drone-facebook/).

## Feature

* [x] Send with Text Message (Support [handlebar](https://github.com/aymerick/raymond) template).
* [x] Send with New Image.
* [x] Send with New Audio.
* [x] Send with New Video.
* [x] Send with New File.
* [x] Support [prometheus](https://prometheus.io) metrics API.
* [x] Automatically install TLS certificates from [Let's Encrypt](https://letsencrypt.org/).

## Build

Build the binary with the following commands:

```
$ make build
```

## Testing

Test the package with the following command:

```
$ make test
```

## Docker

Build the docker image with the following commands:

```
$ make docker
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e PLUGIN_FB_PAGE_TOKEN=xxxxxxx \
  -e PLUGIN_FB_VERIFY_TOKEN=xxxxxxx \
  -e PLUGIN_TO=xxxxxxx \
  -e PLUGIN_MESSAGE=test \
  -e PLUGIN_IMAGES=http://example.com/test.png \
  -e PLUGIN_AUDIOS=http://example.com/test.mp3 \
  -e PLUGIN_VIDEOS=http://example.com/test.mp4 \
  -e PLUGIN_FILES=http://example.com/test.pdf \
  -e PLUGIN_ONLY_MATCH_EMAIL=false \
  -e DRONE_REPO_OWNER=appleboy \
  -e DRONE_REPO_NAME=go-hello \
  -e DRONE_COMMIT_SHA=e5e82b5eb3737205c25955dcc3dcacc839b7be52 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=appleboy \
  -e DRONE_COMMIT_AUTHOR_EMAIL=appleboy@gmail.com \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/appleboy/go-hello \
  -e DRONE_JOB_STARTED=1477550550 \
  -e DRONE_JOB_FINISHED=1477550750 \
  -e DRONE_TAG=1.0.0 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-facebook
```
