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
	"github.com/davecgh/go-spew/spew"
	screp "github.com/icza/screp/rep"
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

	raymond.RegisterHelper("race", func(r string, options *raymond.Options) string {
		return r[0:1]
	})

	go func() {
		fmt.Println("server started on ", addr)
		err := srv.ListenAndServe()
		log.Fatal(err)
	}()

	nux.App().Init(manifest)
	nux.App().Run()
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	res := getReplayData()
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	//fmt.Fprintf(w, "{ \"Version\": %q}", data)

	tpl, err := raymond.ParseFile("./last5.hbs")
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		result := tpl.MustExec(res)
		fmt.Fprint(w, result)
	}
}

func getReplayData() map[string]interface{} {
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
	return compileReplayInfo(destination, r)
}

func compileReplayInfo(out *os.File, rep *screp.Replay) map[string]interface{} {
	rep.Compute()
	//fmt.Printf("%+v\n", rep)
	var winner, loser *screp.Player
	winnerID := rep.Computed.WinnerTeam
	for _, p := range rep.Header.Players {
		if p.ID == winnerID {
			winner = p
		} else {
			loser = p
		}
	}
	spew.Dump(winner)
	spew.Dump(loser)

	engine := rep.Header.Engine.ShortName
	if rep.Header.Version != "" {
		engine = engine + " " + rep.Header.Version
	}
	mapName := rep.MapData.Name
	if mapName == "" {
		mapName = rep.Header.Map // But revert to Header.Map if the latter is not available.
	}

	ctx := map[string]interface{}{
		"winner": winner,
		"loser":  loser,
	}

	return ctx
}
