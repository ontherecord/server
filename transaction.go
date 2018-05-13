package main

type Transaction struct {
	Sender, Receiver, Room string
	Message                string
}

// NewTransaction creates a transaction that will go into the next block, and returns the index of that next block.
func NewTransaction(sender, receiver, room, message string) Index {
	transactions = append(transactions, Transaction{
		Sender:   sender,
		Receiver: receiver,
		Room:     room,
		Message:  message,
	})

	return chain[len(chain)-1].Index + 1
}
