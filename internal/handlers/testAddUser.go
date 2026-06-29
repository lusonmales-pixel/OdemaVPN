package handlers

import (
	"log"
	"net/http"
)

func (e *Env) TestAddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := e.XUIClient.AddUser(ctx, 13, "2b9b6323-6375-4107-b103-2476dbc22a85", 1)
	if err != nil {
		log.Println("Error:", err)
	}
}
