package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"
)

var (
	port        = flag.Int("port", 8023, "The port on which to listen.")
	resolveTime = flag.Duration("resolve", 1*time.Minute, "How often to resolve chains")

	chain *Chain
	nodes map[url.URL]bool
)

func handleError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	io.WriteString(w, err.Error())
	glog.Error(err)
}

func handleMessagesNew(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Creating message")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	block := NewBlock(Message{
		From: r.FormValue("from"),
		To:   r.FormValue("to"),
		Room: r.FormValue("room"),
		Text: r.FormValue("text"),
	})
	glog.Infof("Created %+v", block)

	block, err := chain.Add(block)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, fmt.Sprintf("%d", block.Index))
}

func handleMessagesList(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Returning chain of messages")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bytes, err := json.Marshal(chain)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleNodesRegister(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Registering node")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	glog.Infof("Registering address %q", r.FormValue("address"))

	nodeUrl, err := url.ParseRequestURI(r.FormValue("address"))
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(nodeUrl.Scheme, "http") {
		handleError(w, fmt.Errorf("address %q should start with 'http'", nodeUrl), http.StatusBadRequest)
		return
	}

	nodes[*nodeUrl] = true
	glog.Infof("Registered node %q", nodeUrl.String())

	w.WriteHeader(http.StatusAccepted)
}

func handleNodesList(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Listing nodes")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	keys := []string{}
	for node := range nodes {
		keys = append(keys, node.String())
	}

	bytes, err := json.Marshal(keys)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func resolveLoop() {
	glog.Infof("Starting resolve loop on %v timer", *resolveTime)
	c := time.Tick(*resolveTime)
	for now := range c {
		glog.Infof("Running resolve at %v", now)
		resolve()
	}
}

func resolve() {
	// TODO: lock 'nodes' while resolving
	for node := range nodes {
		glog.Infof("Resolving %q", node.String())
		node.Path = path.Join(node.Path, "/messages/list")
		resp, err := http.Get(node.String())

		if err != nil {
			glog.Error(err)
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		var c Chain
		err = json.Unmarshal(body, &c)

		if err != nil {
			glog.Error(err)
			continue
		}

		if !c.IsValid() {
			glog.Warning(fmt.Errorf("skipping invalid chain from %s", node.String()))
			continue
		}

		// TODO: now resolve.
	}
}

func main() {
	flag.Parse()

	chain = NewChain()
	nodes = make(map[url.URL]bool)

	// Chain handlers
	http.HandleFunc("/messages/new", handleMessagesNew)
	http.HandleFunc("/messages/list", handleMessagesList)

	// Node handlers
	http.HandleFunc("/nodes/register", handleNodesRegister)
	http.HandleFunc("/nodes/list", handleNodesList)

	go resolveLoop()

	glog.Infof("Running at http://localhost:%d", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		glog.Exit(err)
	}
}
