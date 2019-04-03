package commands

import (
	"sczg/config"

	log "github.com/sirupsen/logrus"
)

// StartArchiver archives old entries
// and removes them from "active" table
func StartArchiver(env *config.Env) {
	defer env.DB.Close()
	err := env.DB.ArchiveEntries(env.Cfg.Archive)
	if err != nil {
		log.Fatalf("archiving failed with err: { %v }", err)
	}
	err = env.DB.PurgeOutdatedActive(env.Cfg.Archive)
	if err != nil {
		log.Fatalf("purgind old entries failed with err {%v}", err)
	}
}
