package grafana

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	gapi "github.com/jwierzbo/terraform-provider-grafanads/pkg/api"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var resourceTests = []struct {
	resource   string
	config     string
	attrChecks map[string]string
}{
	{
		"grafanads_data_source_generic.mongoatlas",
		`
		resource "grafanads_data_source_generic" "mongoatlas" {
			type          = "grafana-mongodb-atlas-datasource"
			name          = "mongoatlas-provider-test"
			org_id         = 9
			json_data_string = {
				atlasPublicKey: "xxx",
			}
			secure_json_string = {
				atlasPrivateKey: "yyy",
			}
		}
		`,
		map[string]string{
			"type": "grafana-mongodb-atlas-datasource",
			"name": "mongoatlas-provider-test",
		},
	},
}

func TestAccDataSource_basic(t *testing.T) {
	var dataSource gapi.DataSourceGeneric

	// Iterate over the provided configurations for datasources
	for _, test := range resourceTests {

		// Always check that the resource was created and that `id` is a number
		checks := []resource.TestCheckFunc{
			testAccDataSourceCheckExists(test.resource, &dataSource),
			resource.TestMatchResourceAttr(
				test.resource,
				"id",
				regexp.MustCompile(`\d+`),
			),
		}

		// Add custom checks for specified attribute values
		for attr, value := range test.attrChecks {
			checks = append(checks, resource.TestCheckResourceAttr(
				test.resource,
				attr,
				value,
			))
		}

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccDataSourceCheckDestroy(&dataSource),
			Steps: []resource.TestStep{
				{
					Config: test.config,
					Check: resource.ComposeAggregateTestCheckFunc(
						checks...,
					),
				},
			},
		})
	}
}

func testAccDataSourceCheckExists(rn string, dataSource *gapi.DataSourceGeneric) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("resource id is malformed")
		}

		client := testAccProvider.Meta().(*gapi.Client)
		gotDataSource, err := client.DataSource(id)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}

		*dataSource = *gotDataSource

		return nil
	}
}

func testAccDataSourceCheckDestroy(dataSource *gapi.DataSourceGeneric) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*gapi.Client)
		_, err := client.DataSource(dataSource.Id)
		if err == nil {
			return fmt.Errorf("data source still exists")
		}
		return nil
	}
}
