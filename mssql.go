package ploto

import (
	"fmt"
	"net/url"
	"strings"
)

// Mssql mssql dialector
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
func (m Mssql) GetDialectDSN(database string, config *DialectClientOption) string {
	//https://github.com/denisenkom/go-mssqldb
	//sqlserver://username:password@host:port/instance?param1=value&param2=value

	params := []string{"sqlserver://", config.User, ":", url.QueryEscape(config.Password), "@", config.Host, ":", fmt.Sprintf("%d", config.Port), "?database=", config.Database}

	if config.DialectOptions != nil {
		//存在
		options := config.DialectOptions

		for k, v := range options {
			params = append(params, fmt.Sprintf("&%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
		}
	}

	dnsPath := strings.Join(params, "")
	return dnsPath
}
