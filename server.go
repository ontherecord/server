package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/golang/glog"
)

var (
	port = flag.Int("port", 8023, "The port on which to listen.")

	chain Chain
	nodes map[url.URL]bool
)

func handleMessagesNew(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Creating message")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprint(err))
		glog.Error(err)
		return
	}

	block := NewBlock(Message{
		Sender:   r.FormValue("sender"),
		Receiver: r.FormValue("receiver"),
		Room:     r.FormValue("room"),
		Text:     r.FormValue("text"),
	})
	glog.Infof("Created %+v", block)

	if err := chain.Add(block); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprint(err))
		glog.Error(err)
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

	// TODO: JSON?
	io.WriteString(w, fmt.Sprintf("%+v", chain))
}

func handleNodesRegister(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Registering node")
	glog.Infof("%+v", *r)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprint(err))
		glog.Error(err)
		return
	}
	glog.Infof("%+v", r.Form)
	glog.Infof("%+v", r.PostForm)

	glog.Infof("Registering address %q", r.FormValue("address"))

	nodeUrl, err := url.Parse(r.FormValue("address"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprint(err))
		glog.Error(err)
		return
	}

	nodes[*nodeUrl] = true
	glog.Infof("Registered node %q", nodeUrl.String())

	w.WriteHeader(http.StatusAccepted)
}

func handleNodesResolve(w http.ResponseWriter, r *http.Request) {
}

func handleNodesList(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Listing nodes")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// TODO: JSON?
	io.WriteString(w, fmt.Sprintf("%+v", nodes))
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
	http.HandleFunc("/nodes/resolve", handleNodesResolve)
	http.HandleFunc("/nodes/list", handleNodesList)

	glog.Infof("Running at http://localhost:%d", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		glog.Exit(err)
	}
}
