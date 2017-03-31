package leaderelection

import (
	"time"

	"strings"

	"fmt"

	"github.com/ovh/cds/engine/api/database"
	"github.com/ovh/cds/engine/log"
	"github.com/ovh/cds/sdk"
)

var (
	hostname, cdsname string
	host              *CDSHost
)

// Init initalize the embedded election leader system
func Init(h, s string) {
	hostname = h
	cdsname = s

	heartbeatTicker := time.NewTicker(2 * time.Second).C
	awolTicker := time.NewTicker(time.Second).C

	for {
		_db := database.DB()
		if _db == nil {
			log.Critical("Database unavailable")
		}
		db := database.DBMap(_db)
		select {
		case <-heartbeatTicker:
			heartbeat()

		case <-awolTicker:
			b, err := IsLeader()
			if err != nil {
				log.Critical("IsLeader error : %s", err)
			}
			if b {
				awol()
			}
		}
	}

}

func IsLeader() (bool, error) {
	return false, nil
}

func heartbeat() error {
	time.Sleep(1 * time.Second)
	_db := database.DB()
	if _db == nil {
		return fmt.Errorf("Database unavailable")
	}
	db := database.DBMap(_db)
	res, err := db.Exec("update leader_election_host set heartbeat = $2 where id = $1", host.ID, time.Now())
	if err != nil {
		return sdk.WrapError(err, "leaderelection.registerRoutine> Unable to heartbeat %s", host.Name)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return sdk.WrapError(err, "leaderelection.registerRoutine> Unable to heartbeat %s", host.Name)
	}
	//Insert if no rows has been updated
	if n == 0 {
		host = &CDSHost{
			Hostname:  hostname,
			Name:      cdsname,
			IsLeader:  false,
			Heartbeat: time.Now(),
		}
		if err := db.Insert(host); err != nil {
			return sdk.WrapError(err, "leaderelection.registerRoutine> Unable to insert host %s", host.Name)
		}
	}
	return nil
}

func electionRoutine() {
	for {
		time.Sleep(1 * time.Second)
		_db := database.DB()
		if _db == nil {
			continue
		}
		db := database.DBMap(_db)
		tx, err := db.Begin()
		if err != nil {
			log.Warning("leaderelection.electionRoutine> Unable to start transaction : %s", err)
			continue
		}
		//Load all hosts
		hosts, err := LoadAll(tx)
		if err != nil {
			log.Warning("leaderelection.electionRoutine> Unable to load hosts : %s", err)
			continue
		}
		//Choose a leader, by name
		var leaderName string
		for _, h := range hosts {
			if strings.Compare(h.Name, leaderName) > 0 {
				leaderName = h.Name
			}
		}
		//Update the leader
		if _, err := db.Exec("update leader_election_host set leader = true where name = $1", host.Name); err != nil {
			log.Critical("leaderelection.registerRoutine> Unable to heartbeat %s : %s", host.Name, err)
			continue
		}
		//Update the other
		if _, err := db.Exec("update leader_election_host set leader = false where name != $1", host.Name); err != nil {
			log.Critical("leaderelection.registerRoutine> Unable to heartbeat %s : %s", host.Name, err)
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Warning("leaderelection.electionRoutine> Unable to commit transaction : %s", err)
			continue
		}
	}
}

func awol(db *gorp.DBMap) error {

	db := database.DBMap(_db)
	//Load all hosts
	hosts, err := LoadAll(db)
	if err != nil {
		return sdk.WrapError(err, "leaderelection.awolRoutine> Unable to load hosts")
	}
	for _, h := range hosts {
		//If last Heartbeat is older than 5 seconds ago
		if h.Heartbeat.Before(time.Now().Add(-5 * time.Second)) {
			if _, err := db.Delete(&h); err != nil {
				return sdk.WrapError(err, "leaderelection.awolRoutine> Unable to delete hosts : %s")
			}
		}
	}
	return nil
}
