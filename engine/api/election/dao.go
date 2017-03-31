package leaderelection

import "github.com/go-gorp/gorp"

// LoadAll returns all cds hosts registered in the database
func LoadAll(db gorp.SqlExecutor) ([]CDSHost, error) {
	hosts := []CDSHost{}
	if _, err := db.Select(&hosts, "select * from cds_host"); err != nil {
		return nil, err
	}
	return hosts, nil
}

// Load return a cds host given a hostname and a cdsname
func Load(db gorp.SqlExecutor, hostname, cdsname string) (*CDSHost, error) {
	host := &CDSHost{}
	if err := db.SelectOne(host, "select * from cds_host where hostname = $1 and name = $2"); err != nil {
		return nil, err
	}
	return host, nil
}
