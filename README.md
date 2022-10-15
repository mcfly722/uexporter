# Universal Exporter

![Status: in progress](https://img.shields.io/badge/status-in%20progress-success.svg)
[![License: GPL3.0](https://img.shields.io/badge/License-GPL3.0-blue.svg)](https://www.gnu.org/licenses/gpl-3.0.html)

Universal Exporter is Prometheus client service that runs and parse commands by schedule.<br>
For parsing used JavaScript engine with own API.
### generate hash
```
echo -n somePassword | sha256sum
```


### run

```
go run .
```
### build
```
docker login
```
```
docker build -t mcfly722/uexporter .
```
```
docker push mcfly722/uexporter
```
