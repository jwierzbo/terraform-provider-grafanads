# terraform-provider-grafanads

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

The terraform provider to manage every type of Grafana Datasource

Based on: https://github.com/grafana/grafana-api-golang-client


## Requirements

Tested with:
-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)


## Usage

Example code is located under [sample](sample) directory ([sample/main.tf](sample/main.tf))

```shell script
make install-local
cd sample
terraform init
terraform plan
terraform apply
```


## Tests

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

```
GRAFANA_URL=https://grafana.company.com GRAFANA_AUTH=XYZ make testacc
```
