package ploto

import (
	"encoding/json"
	"testing"
)

func getConfigMssql() (config Configs) {
	testConfig := `{"mssql": {
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
	}}`

	var conf Configs

	json.Unmarshal([]byte(testConfig), &conf)

	// fmt.Printf("conf %+v", conf)
	return conf

}

func TestMssqlDSNConfig(t *testing.T) {
	config := getConfigMssql()

	driver := &Dialect{
		Configs: config.Mssql,
	}

	clientConfig := driver.getClientConfig("test")

	t.Logf("config %+v", clientConfig)
	mssql := Mssql{}

	dnsPath := mssql.GetDialectDSN("test", clientConfig)

	t.Logf("dns %s", dnsPath)
	// t.Logf("dnspath %s", dnsPath)
	if dnsPath != "sqlserver://sa:test123@127.0.0.1:1433?database=test&dial+timeout=10" {
		t.Errorf("mysql GetDialectDSN error ")
	}
}
