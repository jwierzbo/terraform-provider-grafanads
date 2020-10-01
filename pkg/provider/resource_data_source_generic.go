package grafana

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"

	gapi "github.com/jwierzbo/terraform-provider-grafana-datasource/pkg/api"
)

func ResourceDataSourceGeneric() *schema.Resource {
	return &schema.Resource{
		Create: CreateDataSourceGeneric,
		Update: UpdateDataSourceGeneric,
		Delete: DeleteDataSourceGeneric,
		Read:   ReadDataSourceGeneric,

		Schema: map[string]*schema.Schema{
			"access_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "proxy",
			},
			"org_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"basic_auth_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"basic_auth_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "",
				Sensitive: true,
			},
			"basic_auth_username": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "",
				Sensitive: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"json_data_string": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"json_data_bool": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeBool},
				Optional: true,
			},
			"json_data_int": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},

			"secure_json_string": {
				Type:      schema.TypeMap,
				Elem:      &schema.Schema{Type: schema.TypeString},
				Optional:  true,
				Sensitive: true,
			},
			"secure_json_bool": {
				Type:      schema.TypeMap,
				Elem:      &schema.Schema{Type: schema.TypeBool},
				Optional:  true,
				Sensitive: true,
			},
			"secure_json_int": {
				Type:      schema.TypeMap,
				Elem:      &schema.Schema{Type: schema.TypeInt},
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// CreateDataSourceGeneric creates a Grafana datasource
func CreateDataSourceGeneric(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	dataSource, err := makeDataSource(d)
	if err != nil {
		return err
	}

	id, err := client.NewDataSource(dataSource)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(id, 10))

	return ReadDataSourceGeneric(d, meta)
}

// UpdateDataSourceGeneric updates a Grafana datasource
func UpdateDataSourceGeneric(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	dataSource, err := makeDataSource(d)
	if err != nil {
		return err
	}

	return client.UpdateDataSource(dataSource)
}

// ReadDataSourceGeneric reads a Grafana datasource
func ReadDataSourceGeneric(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	idStr := d.Id()
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid id: %#v", idStr)
	}

	dataSource, err := client.DataSource(id)
	if err != nil {
		if err.Error() == "404 Not Found" {
			log.Printf("[WARN] removing datasource %s from state because it no longer exists in grafana", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("id", dataSource.Id)
	d.Set("access_mode", dataSource.Access)
	d.Set("basic_auth_enabled", dataSource.BasicAuth)
	d.Set("basic_auth_username", dataSource.BasicAuthUser)
	d.Set("basic_auth_password", dataSource.BasicAuthPassword)
	d.Set("database_name", dataSource.Database)
	d.Set("is_default", dataSource.IsDefault)
	d.Set("name", dataSource.Name)
	d.Set("password", dataSource.Password)
	d.Set("type", dataSource.Type)
	d.Set("url", dataSource.URL)
	d.Set("org_id", dataSource.OrgId)
	d.Set("username", dataSource.User)

	return nil
}

// DeleteDataSourceGeneric deletes a Grafana datasource
func DeleteDataSourceGeneric(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	idStr := d.Id()
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid id: %#v", idStr)
	}

	return client.DeleteDataSource(id)
}

func makeDataSource(d *schema.ResourceData) (*gapi.DataSourceGeneric, error) {
	idStr := d.Id()
	var id int64
	var err error
	if idStr != "" {
		id, err = strconv.ParseInt(idStr, 10, 64)
	}

	return &gapi.DataSourceGeneric{
		Id:                id,
		OrgId:             int64(d.Get("org_id").(int)),
		Name:              d.Get("name").(string),
		Type:              d.Get("type").(string),
		URL:               d.Get("url").(string),
		Access:            d.Get("access_mode").(string),
		Database:          d.Get("database_name").(string),
		User:              d.Get("username").(string),
		Password:          d.Get("password").(string),
		IsDefault:         d.Get("is_default").(bool),
		BasicAuth:         d.Get("basic_auth_enabled").(bool),
		BasicAuthUser:     d.Get("basic_auth_username").(string),
		BasicAuthPassword: d.Get("basic_auth_password").(string),
		JSONData:          makeJSONData(d),
		SecureJSONData:    makeSecureJSONData(d),
	}, err
}

func mergeMaps(input map[string]interface{}, output gapi.JsonData) {
	for k, v := range input {
		output[k] = v
	}
}

func makeJSONData(d *schema.ResourceData) gapi.JsonData {
	result := gapi.JsonData{}

	strings := d.Get("json_data_string").(map[string]interface{})
	bools := d.Get("json_data_bool").(map[string]interface{})
	ints := d.Get("json_data_int").(map[string]interface{})

	mergeMaps(strings, result)
	mergeMaps(bools, result)
	mergeMaps(ints, result)
	return result
}

func makeSecureJSONData(d *schema.ResourceData) gapi.JsonData {
	result := gapi.JsonData{}

	strings := d.Get("secure_json_string").(map[string]interface{})
	bools := d.Get("secure_json_bool").(map[string]interface{})
	ints := d.Get("secure_json_int").(map[string]interface{})

	mergeMaps(strings, result)
	mergeMaps(bools, result)
	mergeMaps(ints, result)
	return result
}
