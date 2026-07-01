package handlers

import (
	"log"
	"net/http"
)

func (e *Env) TestAddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := e.XUIClient.AddUser(ctx, e.XUIInboundID, "6bfce152-4b82-46b7-9687-94ab1ceb95dc", 2); err != nil {
		log.Println("Error adding client:", err)
	}
}
