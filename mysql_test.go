package ploto

import (
	"encoding/json"
	"testing"
)

func getConfigMysql() (config Configs) {
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

func TestMysqlDSNConfig(t *testing.T) {
	config := getConfigMysql()

	driver := &Dialect{
		Configs: config.Mysql,
	}

	clientConfig := driver.getClientConfig("test")

	t.Logf("config %+v", clientConfig)
	mysql := Mysql{}

	dnsPath := mysql.GetDialectDSN("test", clientConfig)

	t.Logf("dns %s", dnsPath)
	// t.Logf("dnspath %s", dnsPath)
	if dnsPath != "test:asfasdf%40%23sddfsdf@tcp(127.0.0.1:3307)/test?charset=utf8mb4&multiStatements=true&writeTimeout=3000ms&readTimeout=3000ms&timeout=3000ms&parseTime=true" {
		t.Errorf("mysql GetDialectDSN error ")
	}
}
