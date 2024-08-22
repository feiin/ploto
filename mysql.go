package ploto

import (
	"fmt"
	"net/url"
	"strings"
)

// Mysql mysql dialector
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
func (m Mysql) GetDialectDSN(database string, config *DialectClientOption) string {
	//https://github.com/go-sql-driver/mysql
	//[driver[:password]@(host)][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]

	charset := "utf8mb4"
	if len(config.Charset) > 0 {
		charset = config.Charset
	}

	params := []string{config.User, ":", config.Password, "@tcp(", config.Host, ":", fmt.Sprintf("%d", config.Port), ")/", config.Database, "?charset=" + charset}

	if config.DialectOptions != nil {
		//存在
		for k, v := range config.DialectOptions {
			params = append(params, fmt.Sprintf("&%s=%s", k, url.QueryEscape(v)))
		}
	}

	dnsPath := strings.Join(params, "")
	return dnsPath
}
