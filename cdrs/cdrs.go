/*
Rating system designed to be used in VoIP Carriers World
Copyright (C) 2013 ITsysCOM

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

package cdrs

import (
	"encoding/json"
	"fmt"
	utils "github.com/cgrates/cgrates/cgrcoreutils"
	"io/ioutil"
	"log/syslog"
	"net/http"
)

var (
	Logger utils.LoggerInterface
)

func init() {
	var err error
	Logger, err = syslog.New(syslog.LOG_INFO, "CGRateS")
	if err != nil {
		Logger = new(utils.StdLogger)
	}
}

type CDRVars struct {
	FSCdr map[string]string
}

func cdrHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	cdr := CDRVars{}
	if err := json.Unmarshal(body, &cdr); err == nil {
		new(FSCdr).New(body)
	} else {
		Logger.Err(fmt.Sprintf("CDRCAPTOR: Could not unmarshal cdr: %v", err))
	}
}

func startCaptiuringCDRs() {
	http.HandleFunc("/cdr", cdrHandler)
	http.ListenAndServe(":8080", nil)
}
