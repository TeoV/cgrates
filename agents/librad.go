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

package agents

import (
	"fmt"
	"net"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/utils"
	"github.com/cgrates/radigo"
)

// radReplyAppendAttributes appends attributes to a RADIUS reply based on predefined template
func radReplyAppendAttributes(reply *radigo.Packet, rplNM *config.NavigableMap) (err error) {
	for _, val := range rplNM.Values() {
		nmItms, isNMItems := val.([]*config.NMItem)
		if !isNMItems {
			return fmt.Errorf("cannot encode reply value: %s, err: not NMItems", utils.ToJSON(val))
		}
		// find out the first itm which is not an attribute
		var itm *config.NMItem
		if len(nmItms) == 1 {
			itm = nmItms[0]
		}
		if itm.Path[0] == MetaRadReplyCode { // Special case used to control the reply code of RADIUS reply
			if err = reply.SetCodeWithName(utils.IfaceAsString(itm.Data)); err != nil {
				return err
			}
			continue
		}
		var attrName, vendorName string
		if len(itm.Path) > 2 {
			vendorName, attrName = itm.Path[0], itm.Path[1]
		} else {
			attrName = itm.Path[0]
		}

		if err = reply.AddAVPWithName(attrName, utils.IfaceAsString(itm.Data), vendorName); err != nil {
			return err
		}
	}
	return
}

// newRADataProvider constructs a DataProvider
func newRADataProvider(req *radigo.Packet) (dP config.DataProvider) {
	dP = &radiusDP{req: req, cache: config.NewNavigableMap(nil)}
	return
}

// radiusDP implements engine.DataProvider, serving as radigo.Packet data decoder
// decoded data is only searched once and cached
type radiusDP struct {
	req   *radigo.Packet
	cache *config.NavigableMap
}

// String is part of engine.DataProvider interface
// when called, it will display the already parsed values out of cache
func (pk *radiusDP) String() string {
	return utils.ToIJSON(pk.req) // return ToJSON because Packet don't have a string method
}

// FieldAsInterface is part of engine.DataProvider interface
func (pk *radiusDP) FieldAsInterface(fldPath []string) (data interface{}, err error) {
	if len(fldPath) != 1 {
		return nil, utils.ErrNotFound
	}
	if data, err = pk.cache.FieldAsInterface(fldPath); err != nil {
		if err != utils.ErrNotFound { // item found in cache
			return
		}
		err = nil // cancel previous err
	} else {
		return // data found in cache
	}
	if len(pk.req.AttributesWithName(fldPath[0], "")) != 0 {
		data = pk.req.AttributesWithName(fldPath[0], "")[0].GetStringValue()
	}
	pk.cache.Set(fldPath, data, false, false)
	return
}

// FieldAsString is part of engine.DataProvider interface
func (pk *radiusDP) FieldAsString(fldPath []string) (data string, err error) {
	var valIface interface{}
	valIface, err = pk.FieldAsInterface(fldPath)
	if err != nil {
		return
	}
	return utils.IfaceAsString(valIface), nil
}

// AsNavigableMap is part of engine.DataProvider interface
func (pk *radiusDP) AsNavigableMap([]*config.FCTemplate) (
	nm *config.NavigableMap, err error) {
	return nil, utils.ErrNotImplemented
}

// RemoteHost is part of engine.DataProvider interface
func (pk *radiusDP) RemoteHost() net.Addr {
	return utils.NewNetAddr(pk.req.RemoteAddr().Network(), pk.req.RemoteAddr().String())
}

//radauthReq is used to authorize a request
//if User-Password avp is present use PAP auth
//if CHAP-Password is presented use CHAP auth
func radauthReq(req *radigo.Packet, aReq *AgentRequest) (bool, error) {
	// try to get UserPassword from Vars as slice of NMItems
	nmItems, err := aReq.Vars.FieldAsInterface([]string{utils.UserPassword})
	if err != nil {
		return false, err
	}
	userPassAvps := req.AttributesWithName("User-Password", utils.EmptyString)
	chapAVPs := req.AttributesWithName("CHAP-Password", utils.EmptyString)
	if len(userPassAvps) == 0 && len(chapAVPs) == 0 {
		return false, fmt.Errorf("cannot find User-Password or CHAP-Password AVP in request")
	}
	if len(userPassAvps) != 0 {
		if userPassAvps[0].StringValue != nmItems.([]*config.NMItem)[0].Data {
			return false, nil
		}
	} else {
		return radigo.AuthenticateCHAP([]byte(utils.IfaceAsString(nmItems.([]*config.NMItem)[0].Data)),
			req.Authenticator[:], chapAVPs[0].RawValue), nil
	}
	return true, nil
}
