// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"nuxui.org/nuxui/nux"
	_ "nuxui.org/nuxui/ui"

	"github.com/aymerick/raymond"
	"github.com/icza/screp/rep"
	"github.com/icza/screp/repparser"
)

const repPath = "C:\\Users\\Chris\\Documents\\StarCraft\\Maps\\Replays\\LastReplay.rep"

func main() {
	m := http.NewServeMux()

	m.HandleFunc("/", testHandler)
	const addr = "localhost:8080"
	srv := http.Server{
		Handler:      m,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	go func() {
		fmt.Println("server started on ", addr)
		err := srv.ListenAndServe()
		log.Fatal(err)
	}()

	nux.App().Init(manifest)
	nux.App().Run()
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	getReplayData()
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	//fmt.Fprintf(w, "{ \"Version\": %q}", data)

	tpl, err := raymond.ParseFile("./last5.hbs")
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		result := tpl.MustExec(nil)
		fmt.Fprint(w, result)
	}
}

func getReplayData() string {
	cfg := repparser.Config{
		Commands: true,
		MapData:  true,
	}
	r, err := repparser.ParseFileConfig(repPath, cfg)
	if err != nil {
		fmt.Printf("Failed to parse replay: %v\n", err)
		os.Exit(1)
	}
	var destination = os.Stdout
	return getEngine(destination, r)
}

func getEngine(out *os.File, rep *rep.Replay) string {
	rep.Compute()

	engine := rep.Header.Engine.ShortName
	if rep.Header.Version != "" {
		engine = engine + " " + rep.Header.Version
	}
	mapName := rep.MapData.Name
	if mapName == "" {
		mapName = rep.Header.Map // But revert to Header.Map if the latter is not available.
	}

	return engine
}
