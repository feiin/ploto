package ploto

import (
	"database/sql"
	// "reflect"
	// "strings"
	"time"
)

type Dialect struct {
	Clients map[string]*DB
	Configs DialectConfig
	logger  LoggerInterface
}

type DialectConfig struct {
	Clients map[string]interface{} `json:"clients"`
	Default map[string]interface{} `json:"default"`
}

type DialectDSN interface {
	GetDialectDSN(database string, config map[string]interface{}) string
}

//CreateClient  create the db pool for  the database
func (dialect *Dialect) CreateClient(database string, dsn DialectDSN) (db *DB, err error) {

	config := dialect.getClientConfig(database)

	dnsPath := dsn.GetDialectDSN(database, config)

	driverDB, err := sql.Open(config["dialect"].(string), dnsPath)
	if err != nil {
		dialect.logger.Error("connect to mysql database %s error", config["database"].(string))
		return nil, err
	}

	maxIdleConns := 10
	maxLeftTime := 7200
	maxOpenConns := 50

	pool := config["pool"].(map[string]interface{})

	if v, ok := pool["maxIdleConns"]; ok {
		maxIdleConns = int(v.(float64))
	}

	if v, ok := pool["maxLeftTime"]; ok {
		maxLeftTime = int(v.(float64))
	}

	if v, ok := pool["maxOpenConns"]; ok {
		maxOpenConns = int(v.(float64))
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	driverDB.SetMaxIdleConns(maxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	driverDB.SetMaxOpenConns(maxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	driverDB.SetConnMaxLifetime(time.Duration(maxLeftTime) * time.Millisecond)
	if err := driverDB.Ping(); err != nil {
		return nil, err
	}

	if _, ok := config["logging"]; ok {
		// logging := config["logging"].(bool)
		// db.LogMode(logging)
		// db.SetLogger(MyDbLogger{})
	}

	//set db to the clients
	db = &DB{DB: driverDB}

	// logger.Info("create mysql db %s client success", database)

	return db, err

}

//Use get the db's conn
func (dialect *Dialect) Use(database string) (db *DB) {

	client, _ := dialect.Clients[database]
	return client
}

//GetClientConfig get the client config
func (dialect *Dialect) getClientConfig(clientName string) (config map[string]interface{}) {

	config = make(map[string]interface{}, 10)
	for k, v := range dialect.Configs.Default {
		config[k] = v
	}

	clients := dialect.Configs.Clients

	if _, ok := clients[clientName]; !ok {
		return config
	}

	//存在
	client := clients[clientName].(map[string]interface{})

	for k, v := range client {
		config[k] = v
	}

	return config
}

//Close  Close the database
func (dialect *Dialect) Close() error {

	for k, v := range dialect.Clients {
		err := v.Close()
		if err != nil {
			dialect.logger.Info("close db %s error %+v", k, err)
			return err
		}
		dialect.logger.Info("close db %s success", k)

	}
	return nil
}
