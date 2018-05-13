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

	chain        []Block
	transactions []Transaction
	nodeID       uuid.UUID
)

func handleTransactionsNew(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Creating transaction")

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

	index := NewTransaction(r.FormValue("sender"), r.FormValue("receiver"), r.FormValue("room"), r.FormValue("message"))
	glog.Infof("Created %d", index)

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, fmt.Sprint(index))
}

func handleMine(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Mining")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	proof := NewProof(chain[len(chain)-1].Proof)

	NewTransaction("", nodeID.String(), "", "")
	block := NewBlock(proof)

	glog.Infof("Mined %+v", block)

	// TODO(dominic): JSON
	io.WriteString(w, fmt.Sprintf("%+v", block))
}

func handleChain(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Returning chain")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	io.WriteString(w, fmt.Sprintf("%+v", chain))
}

func handleNodesRegister(w http.ResponseWriter, r *http.Request) {

}

func handleNodesResolve(w http.ResponseWriter, r *http.Request) {

}

func main() {
	flag.Parse()

	// Create the genesis block and a node ID for this node.
	NewBlock(NewProof(42))
	nodeID = uuid.New()
	glog.Infof("Node ID: %s", nodeID.String())

	// Chain handlers
	http.HandleFunc("/transactions/new", handleTransactionsNew)
	http.HandleFunc("/mine", handleMine)
	http.HandleFunc("/chain", handleChain)

	// Node handlers
	http.HandleFunc("/nodes/register", handleNodesRegister)
	http.HandleFunc("/nodes/resolve", handleNodesResolve)

	glog.Infof("Running at http://localhost:%d", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		glog.Exit(err)
	}
}
