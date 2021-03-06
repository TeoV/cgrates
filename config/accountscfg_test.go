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

package config

import (
	"reflect"
	"testing"

	"github.com/cgrates/cgrates/utils"
)

func TestAccountSCfgLoadFromJSONCfg(t *testing.T) {
	jsonCfg := &AccountSJsonCfg{
		Enabled:               utils.BoolPointer(true),
		Attributes_conns:      &[]string{utils.MetaInternal},
		Rates_conns:           &[]string{utils.MetaInternal},
		Thresholds_conns:      &[]string{utils.MetaInternal},
		Indexed_selects:       utils.BoolPointer(false),
		String_indexed_fields: &[]string{"*req.index1"},
		Prefix_indexed_fields: &[]string{"*req.index1"},
		Suffix_indexed_fields: &[]string{"*req.index1"},
		Nested_fields:         utils.BoolPointer(true),
		Max_iterations:        utils.IntPointer(1000),
		Max_usage:             utils.StringPointer("200h"),
	}
	usage, err := utils.NewDecimalFromUsage("200h")
	if err != nil {
		t.Error(err)
	}
	expected := &AccountSCfg{
		Enabled:             true,
		AttributeSConns:     []string{utils.ConcatenatedKey(utils.MetaInternal, utils.MetaAttributes)},
		RateSConns:          []string{utils.ConcatenatedKey(utils.MetaInternal, utils.MetaRateS)},
		ThresholdSConns:     []string{utils.ConcatenatedKey(utils.MetaInternal, utils.MetaThresholds)},
		IndexedSelects:      false,
		StringIndexedFields: &[]string{"*req.index1"},
		PrefixIndexedFields: &[]string{"*req.index1"},
		SuffixIndexedFields: &[]string{"*req.index1"},
		NestedFields:        true,
		MaxIterations:       1000,
		MaxUsage:            usage,
	}
	jsnCfg := NewDefaultCGRConfig()
	if err = jsnCfg.accountSCfg.loadFromJSONCfg(jsonCfg); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(expected, jsnCfg.accountSCfg) {
		t.Errorf("\nExpecting <%+v>,\n Received <%+v>", utils.ToJSON(expected), utils.ToJSON(jsnCfg.accountSCfg))
	}
}

func TestAccountsCfLoadConfigError(t *testing.T) {
	accountsJson := &AccountSJsonCfg{
		Max_usage: utils.StringPointer("invalid_Decimal"),
	}
	actsCfg := new(AccountSCfg)
	expected := "strconv.ParseInt: parsing \"invalid_Decimal\": invalid syntax"
	if err := actsCfg.loadFromJSONCfg(accountsJson); err == nil || err.Error() != expected {
		t.Errorf("Expected %+v, received %+v", expected, err)
	}
}

func TestAccountSCfgAsMapInterface(t *testing.T) {
	cfgJSONStr := `{
"accounts": {								
	"enabled": true,						
	"indexed_selects": false,			
	"attributes_conns": ["*internal:*attributes"],
	"rates_conns": ["*internal:*rates"],
	"thresholds_conns": ["*internal:*thresholds"],					
	"string_indexed_fields": ["*req.index1"],			
	"prefix_indexed_fields": ["*req.index1"],			
	"suffix_indexed_fields": ["*req.index1"],			
	"nested_fields": true,			
    "max_iterations": 100,
    "max_usage": "72h",
},	
}`

	eMap := map[string]interface{}{
		utils.EnabledCfg:             true,
		utils.IndexedSelectsCfg:      false,
		utils.AttributeSConnsCfg:     []string{utils.MetaInternal},
		utils.RateSConnsCfg:          []string{utils.MetaInternal},
		utils.ThresholdSConnsCfg:     []string{utils.MetaInternal},
		utils.StringIndexedFieldsCfg: []string{"*req.index1"},
		utils.PrefixIndexedFieldsCfg: []string{"*req.index1"},
		utils.SuffixIndexedFieldsCfg: []string{"*req.index1"},
		utils.NestedFieldsCfg:        true,
		utils.MaxIterations:          100,
	}
	usage, err := utils.NewDecimalFromUsage("72h")
	if err != nil {
		t.Error(err)
	}
	eMap[utils.MaxUsage] = usage

	if cgrCfg, err := NewCGRConfigFromJSONStringWithDefaults(cfgJSONStr); err != nil {
		t.Error(err)
	} else if rcv := cgrCfg.accountSCfg.AsMapInterface(); !reflect.DeepEqual(eMap, rcv) {
		t.Errorf("Expected: %+v\n Received: %+v", utils.ToJSON(eMap), utils.ToJSON(rcv))
	}
}

func TestAccountSCfgClone(t *testing.T) {
	usage, err := utils.NewDecimalFromUsage("24h")
	if err != nil {
		t.Error(err)
	}
	ban := &AccountSCfg{
		Enabled:             true,
		IndexedSelects:      false,
		AttributeSConns:     []string{"*req.index1"},
		RateSConns:          []string{"*req.index1"},
		ThresholdSConns:     []string{"*req.index1"},
		StringIndexedFields: &[]string{"*req.index1"},
		PrefixIndexedFields: &[]string{"*req.index1", "*req.index2"},
		SuffixIndexedFields: &[]string{"*req.index1"},
		NestedFields:        true,
		MaxIterations:       1000,
		MaxUsage:            usage,
	}
	rcv := ban.Clone()
	if !reflect.DeepEqual(ban, rcv) {
		t.Errorf("\nExpected: %+v\nReceived: %+v", utils.ToJSON(ban), utils.ToJSON(rcv))
	}
	if (rcv.AttributeSConns)[0] = ""; (ban.AttributeSConns)[0] != "*req.index1" {
		t.Errorf("Expected clone to not modify the cloned")
	}
	if (rcv.RateSConns)[0] = ""; (ban.RateSConns)[0] != "*req.index1" {
		t.Errorf("Expected clone to not modify the cloned")
	}
	if (rcv.ThresholdSConns)[0] = ""; (ban.ThresholdSConns)[0] != "*req.index1" {
		t.Errorf("Expected clone to not modify the cloned")
	}
	if (*rcv.StringIndexedFields)[0] = ""; (*ban.StringIndexedFields)[0] != "*req.index1" {
		t.Errorf("Expected clone to not modify the cloned")
	}
	if (*rcv.PrefixIndexedFields)[0] = ""; (*ban.PrefixIndexedFields)[0] != "*req.index1" {
		t.Errorf("Expected clone to not modify the cloned")
	}
	if (*rcv.SuffixIndexedFields)[0] = ""; (*ban.SuffixIndexedFields)[0] != "*req.index1" {
		t.Errorf("Expected clone to not modify the cloned")
	}
}
