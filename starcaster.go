// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
	_ "nuxui.org/nuxui/ui"

	"github.com/aymerick/raymond"
	screp "github.com/icza/screp/rep"
	"github.com/icza/screp/repparser"
)

const repPath = "C:\\Users\\Chris\\Documents\\StarCraft\\Maps\\Replays\\AutoSave\\"

//go:embed last5.hbs
var hbs string

type Home interface {
	nux.Component
}

type home struct {
	*nux.ComponentBase
}

func main() {

	m := http.NewServeMux()
	// fs := http.FileServer(http.Dir("./public"))
	// http.Handle("/", fs)

	// http.HandleFunc("/", serveFiles)

	m.HandleFunc("/", serveFiles)
	const addr = "localhost:8080"
	srv := http.Server{
		Handler:      m,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	raymond.RegisterHelper("race", func(r string, options *raymond.Options) string {
		if len(r) == 0 {
			return ""
		}
		return r[0:1]
	})

	go func() {
		fmt.Println("server started on ", addr)
		err := srv.ListenAndServe()
		log.Fatal("starcaster", "%s", err)
	}()

	nux.App().Init(manifest)
	nux.App().Run()
}

func (me *home) layout() string {
	return ""
}

func (me *home) style() string {
	return ""
}

func NewHome(manifest nux.Attr) Home {
	me := &home{}
	me.ComponentBase = nux.NewComponentBase(me, manifest)
	nux.InflateLayout(me, me.layout(), nux.InflateStyle(me.style()))
	return me
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	p := "." + r.URL.Path
	if p == "./" {
		testHandler(w, r)
		return
	}
	http.ServeFile(w, r, p)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	res := getTopReplaysData()

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)

	tpl, _ := raymond.Parse(hbs)
	ctx := map[string]interface{}{
		"matches": res,
	}
	//spew.Dump(ctx)
	result := tpl.MustExec(ctx)
	fmt.Fprint(w, result)
}

func getTopReplaysData() []map[string]interface{} {
	var rdlist []map[string]interface{}

	folder := getNewestFolder(repPath)
	repfiles := getNewestFiles(repPath+folder, 5)
	for _, fpath := range repfiles {
		repdata := getReplayData(repPath + folder + "\\" + fpath)
		rdlist = append(rdlist, repdata)
	}
	return rdlist
}

func getNewestFiles(path string, cnt int) []string {
	list := []string{}
	i := 0
	fmt.Println("reading file " + path)
	f, _ := os.Open(path)
	fis, _ := f.Readdir(-1)
	f.Close()
	sort.Sort(ByModTime(fis))

	for _, fi := range fis {
		fmt.Println("found file " + fi.Name())
		if i >= cnt {
			fmt.Println("too many files")
			break
		}
		if !fi.Mode().IsRegular() {
			fmt.Println("not regular file")
			continue
		}
		i++
		list = append(list, fi.Name())
	}
	return list
}

func getNewestFolder(path string) string {
	fmt.Println("reading folder " + path)
	f, _ := os.Open(path)
	fis, _ := f.Readdir(-1)
	f.Close()
	sort.Sort(ByModTime(fis))

	for _, fi := range fis {
		if fi.Mode().IsRegular() {
			fmt.Println("not dir")
			continue
		}
		fmt.Println("found dir " + fi.Name())
		return fi.Name()
	}
	return ""
}

func getReplayData(fileName string) map[string]interface{} {
	cfg := repparser.Config{
		Commands: true,
		MapData:  true,
	}
	r, err := repparser.ParseFileConfig(fileName, cfg)
	if err != nil {
		fmt.Printf("Failed to parse replay: %v\n", err)
		os.Exit(1)
	}
	var destination = os.Stdout
	return compileReplayInfo(destination, r)
}

func compileReplayInfo(out *os.File, rep *screp.Replay) map[string]interface{} {
	rep.Compute()
	var winner, loser *screp.Player
	winnerID := rep.Computed.WinnerTeam
	hasWinner := (winnerID != 0)

	for _, p := range rep.Header.Players {
		if p.Team == winnerID {
			winner = p
		} else {
			loser = p
		}
	}
	if !hasWinner {
		winner = rep.Header.Players[0]
	}

	engine := rep.Header.Engine.ShortName
	if rep.Header.Version != "" {
		engine = engine + " " + rep.Header.Version
	}
	mapName := rep.MapData.Name
	if mapName == "" {
		mapName = rep.Header.Map // But revert to Header.Map if the latter is not available.
	}

	d := rep.Header.Duration()

	ctx := map[string]interface{}{
		"winner":    winner,
		"loser":     loser,
		"len":       d.Truncate(time.Second).String(),
		"map":       mapName,
		"hasWinner": hasWinner,
	}

	return ctx
}

type ByModTime []os.FileInfo

func (fis ByModTime) Len() int {
	return len(fis)
}

func (fis ByModTime) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

func (fis ByModTime) Less(i, j int) bool {
	return fis[i].ModTime().After(fis[j].ModTime())
}
