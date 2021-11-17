# ploto
golang mysql/mssql 数据库简单封装


## DialectConfig 配置结构说明

配置支持多数据库连接，格式如下：

### mysql 
```json
{"mysql": {
		"clients": {
			"test":{
				"host": "127.0.0.1",
				"port": 3307,
				"user": "test",
				"password": "asfasdf@#sddfsdf",
				"database": "test"
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
				"parseTime":true,
				"multiStatements": true,
				"writeTimeout": "3000ms",
				"readTimeout": "3000ms",
				"timeout":"3000ms"
			}	
		}
	}}
```
更多dialectOptions参数见: https://github.com/go-sql-driver/mysql#parameters
### mssql

```
{"mssql": {
		"clients": {
	 
			"test":{
				"host": "127.0.0.1",
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
	}}
```
更多dialectOptions 参数见:https://github.com/denisenkom/go-mssqldb#connection-parameters-and-dsn


## Using

```
package main

import (
    "encoding/json"
    "fmt"
    "github.com/feiin/ploto"
     _ "github.com/go-sql-driver/mysql"
)

func getConfig() (config Configs) {
    testConfig := `{"mysql": {
        "clients": {
            "test":{
                "host": "127.0.0.1",
                "port": 3306,
                "user": "root",
                "password": "root",
                "database": "test"
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
                "parseTime":true,
                "multiStatements": true,
                "writeTimeout": "3000ms",
                "readTimeout": "3000ms",
                "timeout":"3000ms"
            }   
        }
    }}`

    var conf Configs

    json.Unmarshal([]byte(testConfig), &conf)

    // fmt.Printf("conf %+v", conf)
    return conf

}

type User struct {
    Id          int64  `db:"id"`
    Name        string `db:"name"`
    CreatedTime string `db:"created_time"`
    UpdatedTime string `db:"updated_time"`
}

type Configs struct {
    Mysql ploto.DialectConfig `json:"mysql"`
   // Mssql ploto.DialectConfig `json:"mssql"`
}

func main() {

    configs := getConfig()
    mysqlDSN := ploto.Mysql{}
    db, err := ploto.Open(configs.Mysql, mysqlDSN, nil)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    
    var users []User
    err = db.Use("test").Query("select * from users where id<100").Scan(&users)
    if err != nil {
        panic(err)
    }
    fmt.Printf("users %+v", users)
}

```

