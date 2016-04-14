package httpServices

import(
	"net/http"
	"fmt"
)

func init() {
	fmt.Println("helloWorld httpServices initialized.")	
	http.HandleFunc("/web/helloWorld/SayHello", sayHello)
}

func sayHello(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Say Hello To Go Core")
}