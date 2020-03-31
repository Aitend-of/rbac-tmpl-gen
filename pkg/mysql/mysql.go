package mysql

import (
	"database/sql"

	tmplgen "github.com/aitend-of/rbac-tmpl-gen"
	_ "github.com/go-sql-driver/mysql"
)

// Open mysql DB
func Open(dataSourceName string) (*sql.DB, error) {
	return sql.Open("mysql", dataSourceName)
}

// Insert output YAML to DB
func InsertServiceTmpl(db *sql.DB, serviceTmpl *tmplgen.ServiceTmpl) {
	stmt, err := db.Prepare("INSERT INTO features(feature_id, feature_name, feature_descr, service_name, endpoint_path, endpoint_method) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	for _, currFeature := range serviceTmpl.Features {
		for _, endpoint := range currFeature.Endpoints {
			for path, method := range endpoint {
				_, err := stmt.Exec(currFeature.ID, currFeature.FeatureName, currFeature.Description, serviceTmpl.ServiceName, path, method)
				if err != nil {
					panic(err.Error())
				}
			}
		}
	}
}
