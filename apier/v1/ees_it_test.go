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

package v1

import (
	"io/ioutil"
	"net/rpc"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

var (
	eeSCfgPath   string
	eeSCfg       *config.CGRConfig
	eeSRPC       *rpc.Client
	eeSConfigDIR string //run tests for specific configuration

	sTestsEEs = []func(t *testing.T){
		testEEsPrepareFolder,
		testEEsInitCfg,
		testEEsInitDataDb,
		testEEsResetStorDb,
		testEEsStartEngine,
		testEEsRPCConn,
		testEEsAddCDRs,
		testEEsExportCDRs,
		testEEsVerifyExports,
		testEEsKillEngine,
		testEEsCleanFolder,
	}
)

//Test start here
func TestExportCDRs(t *testing.T) {
	switch *dbType {
	case utils.MetaInternal:
		eeSConfigDIR = "ees_internal"
	case utils.MetaMySQL:
		eeSConfigDIR = "ees_mysql"
	case utils.MetaMongo:
		eeSConfigDIR = "ees_mongo"
	case utils.MetaPostgres:
		t.SkipNow()
	default:
		t.Fatal("Unknown Database type")
	}
	for _, stest := range sTestsEEs {
		t.Run(eeSConfigDIR, stest)
	}
}

func testEEsPrepareFolder(t *testing.T) {
	for _, dir := range []string{"/tmp/testCSV"} {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal("Error removing folder: ", dir, err)
		}
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			t.Fatal("Error creating folder: ", dir, err)
		}
	}
}

func testEEsInitCfg(t *testing.T) {
	var err error
	eeSCfgPath = path.Join(*dataDir, "conf", "samples", eeSConfigDIR)
	eeSCfg, err = config.NewCGRConfigFromPath(eeSCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	eeSCfg.DataFolderPath = alsPrfDataDir // Share DataFolderPath through config towards StoreDb for Flush()
	config.SetCgrConfig(eeSCfg)
}

func testEEsInitDataDb(t *testing.T) {
	if err := engine.InitDataDb(eeSCfg); err != nil {
		t.Fatal(err)
	}
}

// Wipe out the cdr database
func testEEsResetStorDb(t *testing.T) {
	if err := engine.InitStorDb(eeSCfg); err != nil {
		t.Fatal(err)
	}
}

// Start CGR Engine
func testEEsStartEngine(t *testing.T) {
	if _, err := engine.StopStartEngine(eeSCfgPath, *waitRater); err != nil {
		t.Fatal(err)
	}
}

// Connect rpc client to rater
func testEEsRPCConn(t *testing.T) {
	var err error
	eeSRPC, err = newRPCClient(eeSCfg.ListenCfg()) // We connect over JSON so we can also troubleshoot if needed
	if err != nil {
		t.Fatal(err)
	}
}

func testEEsAddCDRs(t *testing.T) {
	//add a default charger
	chargerProfile := &ChargerWithCache{
		ChargerProfile: &engine.ChargerProfile{
			Tenant:       "cgrates.org",
			ID:           "Default",
			RunID:        utils.MetaRaw,
			AttributeIDs: []string{"*none"},
			Weight:       20,
		},
	}
	var result string
	if err := eeSRPC.Call(utils.APIerSv1SetChargerProfile, chargerProfile, &result); err != nil {
		t.Error(err)
	} else if result != utils.OK {
		t.Error("Unexpected reply returned", result)
	}
	storedCdrs := []*engine.CDR{
		{CGRID: "Cdr1",
			OrderID: 1, ToR: utils.VOICE, OriginID: "OriginCDR1", OriginHost: "192.168.1.1", Source: "test",
			RequestType: utils.META_RATED, Tenant: "cgrates.org",
			Category: "call", Account: "1001", Subject: "1001", Destination: "+4986517174963", SetupTime: time.Now(),
			AnswerTime: time.Now(), RunID: utils.MetaDefault, Usage: time.Duration(10) * time.Second,
			ExtraFields: map[string]string{"field_extr1": "val_extr1", "fieldextr2": "valextr2"}, Cost: 1.01,
		},
		{CGRID: "Cdr2",
			OrderID: 2, ToR: utils.VOICE, OriginID: "OriginCDR2", OriginHost: "192.168.1.1", Source: "test2",
			RequestType: utils.META_RATED, Tenant: "cgrates.org", Category: "call",
			Account: "1001", Subject: "1001", Destination: "+4986517174963", SetupTime: time.Now(),
			AnswerTime: time.Now(), RunID: utils.MetaDefault, Usage: time.Duration(5) * time.Second,
			ExtraFields: map[string]string{"field_extr1": "val_extr1", "fieldextr2": "valextr2"}, Cost: 1.01,
		},
		{CGRID: "Cdr3",
			OrderID: 3, ToR: utils.VOICE, OriginID: "OriginCDR3", OriginHost: "192.168.1.1", Source: "test2",
			RequestType: utils.META_RATED, Tenant: "cgrates.org", Category: "call",
			Account: "1001", Subject: "1001", Destination: "+4986517174963", SetupTime: time.Now(),
			AnswerTime: time.Now(), RunID: utils.MetaDefault, Usage: time.Duration(30) * time.Second,
			ExtraFields: map[string]string{"field_extr1": "val_extr1", "fieldextr2": "valextr2"}, Cost: 1.01,
		},
		{CGRID: "Cdr4",
			OrderID: 4, ToR: utils.VOICE, OriginID: "OriginCDR4", OriginHost: "192.168.1.1", Source: "test3",
			RequestType: utils.META_RATED, Tenant: "cgrates.org", Category: "call",
			Account: "1001", Subject: "1001", Destination: "+4986517174963", SetupTime: time.Now(),
			AnswerTime: time.Time{}, RunID: utils.MetaDefault, Usage: time.Duration(0) * time.Second,
			ExtraFields: map[string]string{"field_extr1": "val_extr1", "fieldextr2": "valextr2"}, Cost: 1.01,
		},
	}
	for _, cdr := range storedCdrs {
		var reply string
		if err := eeSRPC.Call(utils.CDRsV1ProcessCDR, &engine.CDRWithOpts{CDR: cdr}, &reply); err != nil {
			t.Error("Unexpected error: ", err.Error())
		} else if reply != utils.OK {
			t.Error("Unexpected reply received: ", reply)
		}
	}
	time.Sleep(100 * time.Millisecond)
}

func testEEsExportCDRs(t *testing.T) {
	attr := &utils.ArgExportCDRs{
		ExporterIDs: []string{"CSVExporter"},
	}
	var rply string
	if err := eeSRPC.Call(utils.APIerSv1ExportCDRs, &attr, &rply); err != nil {
		t.Error("Unexpected error: ", err.Error())
	}
	time.Sleep(1 * time.Second)
}

func testEEsVerifyExports(t *testing.T) {
	var files []string
	err := filepath.Walk("/tmp/testCSV/", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, utils.CSVSuffix) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	if len(files) != 1 {
		t.Errorf("Expected %+v, received: %+v", 1, len(files))
	}
	eCnt := "Cdr3,*raw,*voice,OriginCDR3,*rated,cgrates.org,call,1001,1001,+4986517174963,2020-08-30T14:40:32+03:00,2020-08-30T14:40:32+03:00,30s,-1\n" +
		"Cdr4,*raw,*voice,OriginCDR4,*rated,cgrates.org,call,1001,1001,+4986517174963,2020-08-30T14:40:32+03:00,0001-01-01T00:00:00Z,0s,0\n" +
		"Cdr1,*raw,*voice,OriginCDR1,*rated,cgrates.org,call,1001,1001,+4986517174963,2020-08-30T14:40:32+03:00,2020-08-30T14:40:32+03:00,10s,-1\n" +
		"Cdr2,*raw,*voice,OriginCDR2,*rated,cgrates.org,call,1001,1001,+4986517174963,2020-08-30T14:40:32+03:00,2020-08-30T14:40:32+03:00,5s,-1\n"
	if outContent1, err := ioutil.ReadFile(files[0]); err != nil {
		t.Error(err)
	} else if len(eCnt) != len(string(outContent1)) {
		t.Errorf("Expecting: \n<%q>, \nreceived: \n<%q>", eCnt, string(outContent1))
	}
}

func testEEsKillEngine(t *testing.T) {
	if err := engine.KillEngine(100); err != nil {
		t.Error(err)
	}
}

func testEEsCleanFolder(t *testing.T) {
	for _, dir := range []string{"/tmp/testCSV"} {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal("Error removing folder: ", dir, err)
		}
	}
}