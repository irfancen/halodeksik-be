package appdb

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconfig"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*sql.DB, error) {
	psqlInfo, err := getDataSourceName()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, err
}

func getDataSourceName() (string, error) {
	var (
		host     = appconfig.Config.DbHost
		port     = appconfig.Config.DbPort
		user     = appconfig.Config.DbUser
		password = appconfig.Config.DbPassword
		dbname   = appconfig.Config.DbName
	)

	portAsInt, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("bad port number: %v. Err: %v", port, err)
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s application_name=demo_practice sslmode=disable",
		host, portAsInt, user, password, dbname)

	return psqlInfo, nil
}
