package ploto

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

//Mssql mssql dialector
type Mssql struct {
}

//GetDialectDSN
/***config:{
	 	"clients": {
			"share":{
				"host": "127.0.0.1",
				"user": "sa",
				"password": "test123",
				"database": "test"
			},
			"test":{
				"host": "127.0.0.1",
				"port": 1433,
				"user": "sa",
				"password": "test123",
				"database": "test",
				"pool": {
					"maxIdleConns": 20,
					"maxLeftTime": 60000,
					"maxOpenConns": 50
				},
				"dialectOptions": {
					"dial timeout": "10"

				}
			}
		},
		"default": {
			"port": 1433,
			"dialect": "sqlserver",
			"pool": {
				"maxIdleConns": 2,
				"maxLeftTime": 60000,
				"maxOpenConns": 5
			},
			"dialectOptions": {
				"dial timeout": "3"
			}
		}
	}
**/
func (m Mssql) GetDialectDSN(database string, config map[string]interface{}) string {
	//https://github.com/denisenkom/go-mssqldb
	//sqlserver://username:password@host:port/instance?param1=value&param2=value

	params := []string{"sqlserver://", config["user"].(string), ":", url.QueryEscape(config["password"].(string)), "@", config["host"].(string), ":", fmt.Sprintf("%d", int(config["port"].(float64))), "?database=", config["database"].(string)}

	if _, ok := config["dialectOptions"]; ok {
		//存在
		options := config["dialectOptions"].(map[string]interface{})

		for k, v := range options {
			t := reflect.TypeOf(v)

			switch t.Kind() {
			case reflect.String:
				params = append(params, fmt.Sprintf("&%s=%s", url.QueryEscape(k), url.QueryEscape(v.(string))))
			case reflect.Float64:
				params = append(params, fmt.Sprintf("&%s=%d", url.QueryEscape(k), int(v.(float64))))
			case reflect.Bool:
				params = append(params, fmt.Sprintf("&%s=%t", url.QueryEscape(k), v.(bool)))

			}

		}
	}

	dnsPath := strings.Join(params, "")
	return dnsPath
}
