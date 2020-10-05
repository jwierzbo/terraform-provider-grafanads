resource "grafanads_data_source_generic" "mongodbatlas" {
  type   = "grafana-mongodb-atlas-datasource"
  name   = "mongoatlas-datasource-provider-test"
  org_id = 9

  json_data_string = {
    atlasPublicKey = "XXX"
  }

  secure_json_string = {
    atlasPrivateKey = "YYY"
  }
}