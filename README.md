# ploto

A go Library for scan database/sql rows to struct、slice、other types.  And it support multiple databases connection management

It's not an ORM. works with `database/sql`

## 功能

- Scan rows, 支持struct,slice,map,其他基本类型
- 多数据库配置连接管理


## 说明

仅对database/sql的DB.Query,DB.QueryContext进行封装，其他使用保持不变，Query结果Scan支持*Slice、*Struct、*Map、*int等基本类型.



## 数据库配置

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
			"dialect": "sqlserver", //or mssql
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

### 配合多数据库管理一起使用

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


### 只用Scan功能

> 支持对rows结果转化到struct,slice，int等

```golang

package main

import (
	"database/sql"
	"fmt"
	"github.com/feiin/ploto"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@/database")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	//scan rows to slices
	var users []User
	rows, err = db.Query("select * from users where id<100")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := ploto.Scan(rows, &user)
		users = append(users, user)
	}


	//ScanResult等同上代码
	var users []User
	rows, err = db.Query("select * from users where id<100")
	if err != nil {
		panic(err)
	}

	//No need to Close
	err := ploto.ScanResult(rows, &users)

	//.....
	// select count(1) as cnt from users

	if rows.Next() {
		var a int64
		ploto.Scan(rows,&a)
	}
	//.....

	// select * from users where id=1

	if rows.Next() {
		var user User 
		ploto.Scan(rows,&user)
	}
	//.....
}

```
