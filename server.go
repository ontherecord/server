package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

var (
	port = flag.Int("port", 8023, "The port on which to listen.")

	chain  []Block
	nodeID uuid.UUID
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

}

func handleNodesResolve(w http.ResponseWriter, r *http.Request) {

}

func main() {
	flag.Parse()

	nodeID = uuid.New()
	glog.Infof("Node ID: %s", nodeID.String())

	// Create the genesis block and a node ID for this node.
	NewBlock(Message{
		Sender:   nodeID.String(),
		Receiver: nodeID.String(),
		Text:     "[genesis]",
	})

	// Chain handlers
	http.HandleFunc("/messages/new", handleMessagesNew)
	http.HandleFunc("/messages/list", handleMessagesList)

	// Node handlers
	http.HandleFunc("/nodes/register", handleNodesRegister)
	http.HandleFunc("/nodes/resolve", handleNodesResolve)

	glog.Infof("Running at http://localhost:%d", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		glog.Exit(err)
	}
}
