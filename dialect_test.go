package ploto

import (
	"encoding/json"
	"testing"
)

type Configs struct {
	Mysql DialectConfig `json:"mysql"`
	Mssql DialectConfig `json:"mssql"`
}

func getConfig() (config Configs) {
	testConfig := `{"mysql": {
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
	}}`

	var conf Configs

	json.Unmarshal([]byte(testConfig), &conf)

	// fmt.Printf("conf %+v", conf)
	return conf

}

func TestGetClientConfig(t *testing.T) {
	config := getConfig()

	driver := &Dialect{
		Configs: config.Mysql,
	}

	clientConfig := driver.getClientConfig("test")
	t.Logf("config %+v", clientConfig)
	if int(clientConfig["port"].(float64)) != 3307 {

		t.Errorf("get port failed %v", clientConfig["port"])
	}

}
