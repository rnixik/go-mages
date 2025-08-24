package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/rnixik/go-mages/internal/game"
	"github.com/rnixik/go-mages/internal/lobby"
	"github.com/rnixik/go-mages/internal/transport"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var addr = flag.String("addr", "127.0.0.1:8009", "http service address")
var serveFiles = flag.Bool("serveFiles", true, "use this app to serve static files (js, css, images)")
var appEnv = flag.String("env", "local", "application environment: local, production")

var indexPageContent []byte

func serveIndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Write(indexPageContent)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/favicon.ico")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	indexPageContentRaw, err := ioutil.ReadFile("web/index.html")
	if err != nil {
		log.Fatal("Read index.html error: ", err)
	}
	version, err := ioutil.ReadFile("version")
	if err != nil {
		log.Println("Cannot read file 'version': ", err)
	}
	indexPageContent = bytes.Replace(indexPageContentRaw, []byte("%APP_ENV%"), []byte(*appEnv), 1)
	indexPageContent = bytes.Replace(indexPageContent, []byte("%APP_VERSION%"), bytes.TrimSpace([]byte(version)), 2)

	newGameFunc := func(playersClients []lobby.ClientPlayer, room *lobby.Room, broadcastEventFunc func(event interface{})) lobby.GameEventsDispatcher {
		return game.NewGame(playersClients, room, broadcastEventFunc)
	}

	newBotFunc := func(botId uint64, room *lobby.Room, sendGameCommand func(client lobby.ClientPlayer, commandName string, commandData json.RawMessage)) lobby.ClientPlayer {
		return game.NewBotClient(botId, room, sendGameCommand)
	}

	matchMaker := game.NewMatchMaker()

	lobbyInstance := lobby.NewLobby(newGameFunc, newBotFunc, matchMaker, 2)
	go lobbyInstance.Run()
	http.HandleFunc("/", serveIndexPage)
	if *serveFiles {
		http.HandleFunc("/favicon.ico", faviconHandler)
		http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))
		http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
		http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/assets"))))
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		transport.ServeWebSocketRequest(lobbyInstance, w, r)
	})
	log.Printf("Listening http://%s", *addr)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
