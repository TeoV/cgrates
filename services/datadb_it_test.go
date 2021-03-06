// +build integration

/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/
package services

import (
	"path"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/cgrates/rpcclient"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/cores"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/servmanager"
	"github.com/cgrates/cgrates/utils"
)

func TestDataDBReload(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()

	utils.Logger, _ = utils.Newlogger(utils.MetaSysLog, cfg.GeneralCfg().NodeID)
	utils.Logger.SetLogLevel(7)

	shdChan := utils.NewSyncedChan()
	shdWg := new(sync.WaitGroup)
	chS := engine.NewCacheS(cfg, nil, nil)
	filterSChan := make(chan *engine.FilterS, 1)
	filterSChan <- nil
	close(chS.GetPrecacheChannel(utils.CacheAttributeProfiles))
	close(chS.GetPrecacheChannel(utils.CacheAttributeFilterIndexes))
	server := cores.NewServer(nil)
	srvDep := map[string]*sync.WaitGroup{utils.DataDB: new(sync.WaitGroup)}
	srvMngr := servmanager.NewServiceManager(cfg, shdChan, shdWg)
	cM := engine.NewConnManager(cfg, nil)
	db := NewDataDBService(cfg, cM, srvDep)
	anz := NewAnalyzerService(cfg, server, filterSChan, shdChan, make(chan rpcclient.ClientConnector, 1), srvDep)
	srvMngr.AddServices(NewAttributeService(cfg, db,
		chS, filterSChan, server, make(chan rpcclient.ClientConnector, 1), anz, srvDep),
		NewLoaderService(cfg, db, filterSChan, server, make(chan rpcclient.ClientConnector, 1), nil, anz, srvDep), db)
	if err := srvMngr.StartServices(); err != nil {
		t.Error(err)
	}
	if db.IsRunning() {
		t.Errorf("Expected service to be down")
	}
	var reply string
	cfg.AttributeSCfg().Enabled = true
	if err := cfg.V1ReloadConfig(&config.ReloadArgs{
		Path:    path.Join("/usr", "share", "cgrates", "conf", "samples", "tutmongo"),
		Section: config.DATADB_JSN,
	}, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("Expecting OK ,received %s", reply)
	}
	time.Sleep(10 * time.Millisecond) //need to switch to gorutine
	if !db.IsRunning() {
		t.Errorf("Expected service to be running")
	}
	getDm := db.GetDM()
	if !reflect.DeepEqual(getDm, db.dm) {
		t.Errorf("\nExpecting <%+v>,\n Received <%+v>", db.dm, getDm)
	}
	oldcfg := &config.DataDbCfg{
		DataDbType: utils.Mongo,
		DataDbHost: "127.0.0.1",
		DataDbPort: "27017",
		DataDbName: "10",
		DataDbUser: "cgrates",
		Opts: map[string]interface{}{
			utils.QueryTimeoutCfg:            "10s",
			utils.RedisClusterOnDownDelayCfg: "0",
			utils.RedisClusterSyncCfg:        "5s",
			utils.RedisClusterCfg:            false,
			utils.RedisSentinelNameCfg:       "",
			utils.RedisTLS:                   false,
			utils.RedisClientCertificate:     "",
			utils.RedisClientKey:             "",
			utils.RedisCACertificate:         "",
		},
		RmtConns: []string{},
		RplConns: []string{},
		Items: map[string]*config.ItemOpt{
			utils.MetaAccounts: {
				Replicate: false,
				Remote:    false},
			utils.MetaReverseDestinations: {
				Replicate: false,
				Remote:    false},
			utils.MetaDestinations: {
				Replicate: false,
				Remote:    false},
			utils.MetaRatingPlans: {
				Replicate: false,
				Remote:    false},
			utils.MetaRatingProfiles: {
				Replicate: false,
				Remote:    false},
			utils.MetaActions: {
				Replicate: false,
				Remote:    false},
			utils.MetaActionPlans: {
				Replicate: false,
				Remote:    false},
			utils.MetaAccountActionPlans: {
				Replicate: false,
				Remote:    false},
			utils.MetaActionTriggers: {
				Replicate: false,
				Remote:    false},
			utils.MetaSharedGroups: {
				Replicate: false,
				Remote:    false},
			utils.MetaTimings: {
				Replicate: false,
				Remote:    false},
			utils.MetaResourceProfile: {
				Replicate: false,
				Remote:    false},
			utils.MetaStatQueues: {
				Replicate: false,
				Remote:    false},
			utils.MetaResources: {
				Replicate: false,
				Remote:    false},
			utils.MetaStatQueueProfiles: {
				Replicate: false,
				Remote:    false},
			utils.MetaThresholds: {
				Replicate: false,
				Remote:    false},
			utils.MetaThresholdProfiles: {
				Replicate: false,
				Remote:    false},
			utils.MetaFilters: {
				Replicate: false,
				Remote:    false},
			utils.MetaRouteProfiles: {
				Replicate: false,
				Remote:    false},
			utils.MetaAttributeProfiles: {
				Replicate: false,
				Remote:    false},
			utils.MetaDispatcherHosts: {
				Replicate: false,
				Remote:    false},
			utils.MetaChargerProfiles: {
				Replicate: false,
				Remote:    false},
			utils.MetaDispatcherProfiles: {
				Replicate: false,
				Remote:    false},
			utils.MetaLoadIDs: {
				Replicate: false,
				Remote:    false},
			utils.MetaIndexes: {
				Replicate: false,
				Remote:    false},
			utils.MetaRateProfiles: {
				Replicate: false,
				Remote:    false},
			utils.MetaActionProfiles: {
				Replicate: false,
				Remote:    false},
			utils.MetaAccountProfiles: {
				Replicate: false,
				Remote:    false},
		},
	}
	if !reflect.DeepEqual(oldcfg, db.oldDBCfg) {
		t.Errorf("Expected %s \n received:%s", utils.ToJSON(oldcfg), utils.ToJSON(db.oldDBCfg))
	}

	err := db.Reload()
	if err != nil {
		t.Errorf("\nExpecting <nil>,\n Received <%+v>", err)
	}
	cfg.AttributeSCfg().Enabled = false
	cfg.GetReloadChan(config.DATADB_JSN) <- struct{}{}
	time.Sleep(10 * time.Millisecond)
	if db.IsRunning() {
		t.Errorf("Expected service to be down")
	}
	shdChan.CloseOnce()
	time.Sleep(10 * time.Millisecond)
}
func TestDataDBReload2(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()

	utils.Logger, _ = utils.Newlogger(utils.MetaSysLog, cfg.GeneralCfg().NodeID)
	utils.Logger.SetLogLevel(7)

	shdChan := utils.NewSyncedChan()
	shdWg := new(sync.WaitGroup)
	chS := engine.NewCacheS(cfg, nil, nil)
	filterSChan := make(chan *engine.FilterS, 1)
	filterSChan <- nil
	close(chS.GetPrecacheChannel(utils.CacheAttributeProfiles))
	close(chS.GetPrecacheChannel(utils.CacheAttributeFilterIndexes))
	server := cores.NewServer(nil)
	srvDep := map[string]*sync.WaitGroup{utils.DataDB: new(sync.WaitGroup)}
	srvMngr := servmanager.NewServiceManager(cfg, shdChan, shdWg)
	cM := engine.NewConnManager(cfg, nil)
	db := NewDataDBService(cfg, cM, srvDep)
	anz := NewAnalyzerService(cfg, server, filterSChan, shdChan, make(chan rpcclient.ClientConnector, 1), srvDep)
	srvMngr.AddServices(NewAttributeService(cfg, db,
		chS, filterSChan, server, make(chan rpcclient.ClientConnector, 1), anz, srvDep),
		NewLoaderService(cfg, db, filterSChan, server, make(chan rpcclient.ClientConnector, 1), nil, anz, srvDep), db)
	if err := srvMngr.StartServices(); err != nil {
		t.Error(err)
	}
	if db.IsRunning() {
		t.Errorf("Expected service to be down")
	}
	var reply string
	cfg.AttributeSCfg().Enabled = true
	if err := cfg.V1ReloadConfig(&config.ReloadArgs{
		Path:    path.Join("/usr", "share", "cgrates", "conf", "samples", "tutmongo"),
		Section: config.DATADB_JSN,
	}, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("Expecting OK ,received %s", reply)
	}
	time.Sleep(10 * time.Millisecond) //need to switch to gorutine
	if !db.IsRunning() {
		t.Errorf("Expected service to be running")
	}

	err := db.Reload()
	if err != nil {
		t.Errorf("\nExpecting <nil>,\n Received <%+v>", err)
	}
	if err := cfg.V1ReloadConfig(&config.ReloadArgs{
		Path:    path.Join("/usr", "share", "cgrates", "conf", "samples", "tutmysql"),
		Section: config.DATADB_JSN,
	}, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("Expecting OK ,received %s", reply)
	}

	err = db.Reload()
	if err != nil {
		t.Errorf("\nExpecting <nil>,\n Received <%+v>", err)
	}
	cfg.AttributeSCfg().Enabled = false
	cfg.GetReloadChan(config.DATADB_JSN) <- struct{}{}
	time.Sleep(10 * time.Millisecond)
	if db.IsRunning() {
		t.Errorf("Expected service to be down")
	}
	shdChan.CloseOnce()
	time.Sleep(10 * time.Millisecond)
}

func TestDataDBReload3(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()

	utils.Logger, _ = utils.Newlogger(utils.MetaSysLog, cfg.GeneralCfg().NodeID)
	utils.Logger.SetLogLevel(7)
	cfg.DataDbCfg().DataDbType = ""
	shdChan := utils.NewSyncedChan()
	shdWg := new(sync.WaitGroup)
	chS := engine.NewCacheS(cfg, nil, nil)
	filterSChan := make(chan *engine.FilterS, 1)
	filterSChan <- nil
	close(chS.GetPrecacheChannel(utils.CacheAttributeProfiles))
	close(chS.GetPrecacheChannel(utils.CacheAttributeFilterIndexes))
	server := cores.NewServer(nil)
	srvDep := map[string]*sync.WaitGroup{utils.DataDB: new(sync.WaitGroup)}
	srvMngr := servmanager.NewServiceManager(cfg, shdChan, shdWg)
	cM := engine.NewConnManager(cfg, nil)
	db := NewDataDBService(cfg, cM, srvDep)
	anz := NewAnalyzerService(cfg, server, filterSChan, shdChan, make(chan rpcclient.ClientConnector, 1), srvDep)
	srvMngr.AddServices(NewAttributeService(cfg, db,
		chS, filterSChan, server, make(chan rpcclient.ClientConnector, 1), anz, srvDep),
		NewLoaderService(cfg, db, filterSChan, server, make(chan rpcclient.ClientConnector, 1), nil, anz, srvDep), db)
	if err := srvMngr.StartServices(); err != nil {
		t.Error(err)
	}
	cfg.AttributeSCfg().Enabled = true
	err := db.Start()
	if err == nil {
		t.Errorf("\nExpecting <%+v>,\n Received <%+v>", "unsupported db_type <>", err)
	}
	cfg.AttributeSCfg().Enabled = false
	cfg.GetReloadChan(config.DATADB_JSN) <- struct{}{}
	time.Sleep(10 * time.Millisecond)
	if db.IsRunning() {
		t.Errorf("Expected service to be down")
	}
	shdChan.CloseOnce()
	time.Sleep(10 * time.Millisecond)
}
func TestDataDBReload4(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()

	utils.Logger, _ = utils.Newlogger(utils.MetaSysLog, cfg.GeneralCfg().NodeID)
	utils.Logger.SetLogLevel(7)
	cfg.DataDbCfg().DataDbType = ""
	shdChan := utils.NewSyncedChan()
	shdWg := new(sync.WaitGroup)
	chS := engine.NewCacheS(cfg, nil, nil)
	filterSChan := make(chan *engine.FilterS, 1)
	filterSChan <- nil
	close(chS.GetPrecacheChannel(utils.CacheAttributeProfiles))
	close(chS.GetPrecacheChannel(utils.CacheAttributeFilterIndexes))
	server := cores.NewServer(nil)
	srvDep := map[string]*sync.WaitGroup{utils.DataDB: new(sync.WaitGroup)}
	srvMngr := servmanager.NewServiceManager(cfg, shdChan, shdWg)
	cM := engine.NewConnManager(cfg, nil)
	db := NewDataDBService(cfg, cM, srvDep)
	anz := NewAnalyzerService(cfg, server, filterSChan, shdChan, make(chan rpcclient.ClientConnector, 1), srvDep)
	srvMngr.AddServices(NewAttributeService(cfg, db,
		chS, filterSChan, server, make(chan rpcclient.ClientConnector, 1), anz, srvDep),
		NewLoaderService(cfg, db, filterSChan, server, make(chan rpcclient.ClientConnector, 1), nil, anz, srvDep), db)
	if err := srvMngr.StartServices(); err != nil {
		t.Error(err)
	}
	cfg.SessionSCfg().Enabled = true
	err := db.Start()
	if err != nil {
		t.Errorf("\nExpecting <nil>,\n Received <%+v>", err)
	}
	cfg.SessionSCfg().Enabled = false
	cfg.GetReloadChan(config.DATADB_JSN) <- struct{}{}
	time.Sleep(10 * time.Millisecond)
	if db.IsRunning() {
		t.Errorf("Expected service to be down")
	}
	shdChan.CloseOnce()
	time.Sleep(10 * time.Millisecond)
}

func TestDataDBReload5(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()

	utils.Logger, _ = utils.Newlogger(utils.MetaSysLog, cfg.GeneralCfg().NodeID)
	utils.Logger.SetLogLevel(7)

	shdChan := utils.NewSyncedChan()
	shdWg := new(sync.WaitGroup)
	chS := engine.NewCacheS(cfg, nil, nil)
	filterSChan := make(chan *engine.FilterS, 1)
	filterSChan <- nil
	close(chS.GetPrecacheChannel(utils.CacheAttributeProfiles))
	close(chS.GetPrecacheChannel(utils.CacheAttributeFilterIndexes))
	server := cores.NewServer(nil)
	srvDep := map[string]*sync.WaitGroup{utils.DataDB: new(sync.WaitGroup)}
	srvMngr := servmanager.NewServiceManager(cfg, shdChan, shdWg)
	cM := engine.NewConnManager(cfg, nil)
	db := NewDataDBService(cfg, cM, srvDep)
	anz := NewAnalyzerService(cfg, server, filterSChan, shdChan, make(chan rpcclient.ClientConnector, 1), srvDep)
	srvMngr.AddServices(NewAttributeService(cfg, db,
		chS, filterSChan, server, make(chan rpcclient.ClientConnector, 1), anz, srvDep),
		NewLoaderService(cfg, db, filterSChan, server, make(chan rpcclient.ClientConnector, 1), nil, anz, srvDep), db)
	if err := srvMngr.StartServices(); err != nil {
		t.Error(err)
	}
	if db.IsRunning() {
		t.Errorf("Expected service to be down")
	}
	var reply string
	cfg.AttributeSCfg().Enabled = true
	if err := cfg.V1ReloadConfig(&config.ReloadArgs{
		Path:    path.Join("/usr", "share", "cgrates", "conf", "samples", "tutmongo"),
		Section: config.DATADB_JSN,
	}, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("Expecting OK ,received %s", reply)
	}
	time.Sleep(10 * time.Millisecond) //need to switch to gorutine
	if !db.IsRunning() {
		t.Errorf("Expected service to be running")
	}

	err := db.Reload()
	if err != nil {
		t.Errorf("\nExpecting <nil>,\n Received <%+v>", err)
	}
	if err := cfg.V1ReloadConfig(&config.ReloadArgs{
		Path:    path.Join("/usr", "share", "cgrates", "conf", "samples", "tutmysql"),
		Section: config.DATADB_JSN,
	}, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("Expecting OK ,received %s", reply)
	}
	err = db.Reload()
	if err != nil {
		t.Errorf("\nExpecting <nil>,\n Received <%+v>", err)
	}

	cfg.DataDbCfg().DataDbType = "bad_type"
	err = db.Reload()
	if err == nil {
		t.Errorf("\nExpecting <unsupported db_type <bad_type>>,\n Received <%+v>", err)
	}

	cfg.AttributeSCfg().Enabled = false
	cfg.GetReloadChan(config.DATADB_JSN) <- struct{}{}
	time.Sleep(10 * time.Millisecond)
	if db.IsRunning() {
		t.Errorf("Expected service to be down")
	}
	shdChan.CloseOnce()
	time.Sleep(10 * time.Millisecond)
}

func TestDataDBReload6(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()

	utils.Logger, _ = utils.Newlogger(utils.MetaSysLog, cfg.GeneralCfg().NodeID)
	utils.Logger.SetLogLevel(7)

	shdChan := utils.NewSyncedChan()
	shdWg := new(sync.WaitGroup)
	chS := engine.NewCacheS(cfg, nil, nil)
	filterSChan := make(chan *engine.FilterS, 1)
	filterSChan <- nil
	close(chS.GetPrecacheChannel(utils.CacheAttributeProfiles))
	close(chS.GetPrecacheChannel(utils.CacheAttributeFilterIndexes))
	server := cores.NewServer(nil)
	srvDep := map[string]*sync.WaitGroup{utils.DataDB: new(sync.WaitGroup)}
	srvMngr := servmanager.NewServiceManager(cfg, shdChan, shdWg)
	cM := engine.NewConnManager(cfg, nil)
	db := NewDataDBService(cfg, cM, srvDep)
	anz := NewAnalyzerService(cfg, server, filterSChan, shdChan, make(chan rpcclient.ClientConnector, 1), srvDep)
	srvMngr.AddServices(NewAttributeService(cfg, db,
		chS, filterSChan, server, make(chan rpcclient.ClientConnector, 1), anz, srvDep),
		NewLoaderService(cfg, db, filterSChan, server, make(chan rpcclient.ClientConnector, 1), nil, anz, srvDep), db)
	if err := srvMngr.StartServices(); err != nil {
		t.Error(err)
	}
	if db.IsRunning() {
		t.Errorf("Expected service to be down")
	}
	var reply string
	cfg.AttributeSCfg().Enabled = true
	if err := cfg.V1ReloadConfig(&config.ReloadArgs{
		Path:    path.Join("/usr", "share", "cgrates", "conf", "samples", "tutmongo"),
		Section: config.DATADB_JSN,
	}, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("Expecting OK ,received %s", reply)
	}
	time.Sleep(10 * time.Millisecond) //need to switch to gorutine
	if !db.IsRunning() {
		t.Errorf("Expected service to be running")
	}

	err := db.Reload()
	if err != nil {
		t.Errorf("\nExpecting <nil>,\n Received <%+v>", err)
	}
	if err := cfg.V1ReloadConfig(&config.ReloadArgs{
		Path:    path.Join("/usr", "share", "cgrates", "conf", "samples", "tutmongo"),
		Section: config.DATADB_JSN,
	}, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("Expecting OK ,received %s", reply)
	}
	time.Sleep(10 * time.Millisecond)
	err = db.Reload()
	if err != nil {
		t.Errorf("\nExpecting <nil>,\n Received <%+v>", err)
	}
	db.cfg.DataDbCfg().DataDbType = utils.Mongo
	db.cfg.DataDbCfg().Opts = map[string]interface{}{
		utils.QueryTimeoutCfg: false,
	}
	err = db.Reload()
	if err == nil {
		t.Errorf("\nExpecting <cannot convert field: false to time.Duration>,\n Received <%+v>", err)
	}

	shdChan.CloseOnce()
	time.Sleep(10 * time.Millisecond)
}

func TestDataDBReloadVersion(t *testing.T) {
	cfg, err := config.NewCGRConfigFromPath(path.Join("/usr", "share", "cgrates", "conf", "samples", "tutmongo"))
	if err != nil {
		t.Fatal(err)
	}
	dbConn, err := engine.NewDataDBConn(cfg.DataDbCfg().DataDbType,
		cfg.DataDbCfg().DataDbHost, cfg.DataDbCfg().DataDbPort,
		cfg.DataDbCfg().DataDbName, cfg.DataDbCfg().DataDbUser,
		cfg.DataDbCfg().DataDbPass, cfg.GeneralCfg().DBDataEncoding,
		cfg.DataDbCfg().Opts)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		dbConn.Flush("")
		dbConn.Close()
	}()

	err = dbConn.SetVersions(engine.Versions{
		utils.StatS:          4,
		utils.Accounts:       3,
		utils.Actions:        2,
		utils.ActionTriggers: 2,
		utils.ActionPlans:    3,
		utils.SharedGroups:   2,
		utils.Thresholds:     4,
		utils.Routes:         2,
		// old version for Attributes
		utils.Attributes:          5,
		utils.Timing:              1,
		utils.RQF:                 5,
		utils.Resource:            1,
		utils.Subscribers:         1,
		utils.Destinations:        1,
		utils.ReverseDestinations: 1,
		utils.RatingPlan:          1,
		utils.RatingProfile:       1,
		utils.Chargers:            2,
		utils.Dispatchers:         2,
		utils.LoadIDsVrs:          1,
		utils.RateProfiles:        1,
		utils.ActionProfiles:      1,
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	utils.Logger, _ = utils.Newlogger(utils.MetaSysLog, cfg.GeneralCfg().NodeID)
	utils.Logger.SetLogLevel(7)

	shdChan := utils.NewSyncedChan()
	shdWg := new(sync.WaitGroup)
	chS := engine.NewCacheS(cfg, nil, nil)
	filterSChan := make(chan *engine.FilterS, 1)
	filterSChan <- nil
	close(chS.GetPrecacheChannel(utils.CacheAttributeProfiles))
	close(chS.GetPrecacheChannel(utils.CacheAttributeFilterIndexes))
	server := cores.NewServer(nil)
	srvDep := map[string]*sync.WaitGroup{utils.DataDB: new(sync.WaitGroup)}
	srvMngr := servmanager.NewServiceManager(cfg, shdChan, shdWg)
	cM := engine.NewConnManager(cfg, nil)
	db := NewDataDBService(cfg, cM, srvDep)
	anz := NewAnalyzerService(cfg, server, filterSChan, shdChan, make(chan rpcclient.ClientConnector, 1), srvDep)
	srvMngr.AddServices(NewAttributeService(cfg, db,
		chS, filterSChan, server, make(chan rpcclient.ClientConnector, 1), anz, srvDep),
		NewLoaderService(cfg, db, filterSChan, server, make(chan rpcclient.ClientConnector, 1), nil, anz, srvDep), db)
	srvMngr.StartServices()
	<-shdChan.Done()
	db.dm = nil
	err = db.Reload()
	if err == nil || err.Error() != "can't conver DataDB of type mongo to MongoStorage" {
		t.Fatal(err)
	}
	shdChan.CloseOnce()
	time.Sleep(10 * time.Millisecond)
}
