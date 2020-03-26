CREATE DATABASE rbac_template;

USE rbac_template;

CREATE TABLE features (
	feature_id 		VARCHAR(50),
    feature_name	VARCHAR(50),
    feature_descr	VARCHAR(150),
    service_name	VARCHAR(50),
    endpoint_path	VARCHAR(50),
    endpoint_method	VARCHAR(10)
)

