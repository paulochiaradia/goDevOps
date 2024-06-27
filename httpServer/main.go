package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	auth "github.com/abbot/go-http-auth"
)

func Secret(user, healm string) string {
	if user == "admin" {
		return "$1$KHE.xNk3$zKy4hTxRSWaY1kuqDipSG."
	}
	return ""
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Uso: go run main.go <diretorio> <porta>")
	}
	diretorio := os.Args[1]
	authenticator := auth.NewBasicAuthenticator("meuserver.com", Secret)
	http.HandleFunc("/", authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		http.FileServer(http.Dir(diretorio)).ServeHTTP(w, &r.Request)
	}))

	porta := os.Args[2]
	fmt.Print("Iniciando o Serividor")
	log.Fatal(http.ListenAndServe(":"+porta, nil))
}
