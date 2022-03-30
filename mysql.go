package ploto

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

//Mysql mysql dialector
type Mysql struct {
}

//GetDialectDSN
/***config:{
	 	"clients": {
			"share":{
				"host": "127.0.0.1",
				"user": "test",
				"password": "test123",
				"database": "test"
			},
			"test":{
				"host": "127.0.0.1",
				"port": 3307,
				"user": "test",
				"password": "test123",
				"database": "test",
				"pool": {
					"maxIdleConns": 20,
					"maxLeftTime": 60000,
					"maxOpenConns": 50
				},
				"dialectOptions": {
					"writeTimeout": "2000ms",
					"readTimeout": "2000ms",
					"timeout":"2000ms"
				}
			}
		},
		"default": {
			"port": 3306,
			"dialect": "mysql",
			"pool": {
				"maxIdleConns": 2,
				"maxLeftTime": 60000,
				"maxOpenConns": 5
			},
			"dialectOptions": {
				"writeTimeout": "3000ms",
				"readTimeout": "3000ms",
				"timeout":"3000ms"
			}
		}
	}
**/
func (m Mysql) GetDialectDSN(database string, config map[string]interface{}) string {
	//https://github.com/go-sql-driver/mysql
	//[driver[:password]@(host)][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]

	charset := "utf8mb4"
	if _, ok := config["charset"]; ok {
		charset = config["charset"].(string)
	}

	params := []string{config["user"].(string), ":", config["password"].(string), "@tcp(", config["host"].(string), ":", fmt.Sprintf("%d", int(config["port"].(float64))), ")/", config["database"].(string), "?charset=" + charset}

	if _, ok := config["dialectOptions"]; ok {
		//存在
		options := config["dialectOptions"].(map[string]interface{})

		for k, v := range options {
			t := reflect.TypeOf(v)

			switch t.Kind() {
			case reflect.String:
				params = append(params, fmt.Sprintf("&%s=%s", k, url.QueryEscape(v.(string))))
			case reflect.Float64:
				params = append(params, fmt.Sprintf("&%s=%d", k, int(v.(float64))))
			case reflect.Bool:
				params = append(params, fmt.Sprintf("&%s=%t", k, v.(bool)))

			}

		}
	}

	dnsPath := strings.Join(params, "")
	return dnsPath
}
