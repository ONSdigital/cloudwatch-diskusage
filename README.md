cloudwatch-diskusage
====================

A custom CloudWatch metric for reporting available filesystem storage.

### Getting started

* `go get github.com/ONSdigital/cloudwatch-diskusage`
* `cd $GOPATH/src/github.com/ONSdigital/cloudwatch-diskusage && go build`
* `FILESYSTEMS=<FILESYSTEMS> INSTANCE_ID=<INSTANCE_ID> NAMESPACE=<NAMESPACE> REGION=<REGION> ./cloudwatch-diskusage`

### Configuration

| Envrionment variable | Default | Description
| -------------------- | ------- | -----------
| FILESYSTEMS          |         | Colon seperated list of filesystems
| INSTANCE_ID          |         | Identifier of the instance
| NAMESPACE            |         | Namespace of the metric
| PERSISTENT           | false   | Run persistently
| REGION               |         | Region of the metric

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2016, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
