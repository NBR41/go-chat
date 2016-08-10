package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

const listenAddr = "localhost:4000"

/*
example call to websocket
var sock = new WebSocket("ws://localhost:4000/");
sock.onmessage = function(m) { console.log("Received:", m.data); }
sock.send("Hello!\n")

resources :
http://sii-rennes.developpez.com/articles/un-chat-en-html5-avec-les-websockets/
https://talks.golang.org/2012/chat.slide
https://vimeo.com/53221560
*/

var partner = make(chan io.ReadWriteCloser)

func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func socketHandler(ws *websocket.Conn) {
	log.Println("receive connexion")
	r, w := io.Pipe()
	go func() {
		_, err := io.Copy(io.MultiWriter(w, NewBotChatWriter(chain)), ws)
		w.CloseWithError(err)
	}()
	s := socket{r, ws, make(chan bool)}
	go match(s)
	<-s.done
}

func match(c io.ReadWriteCloser) {
	log.Println("start match")
	fmt.Fprint(c, NewServerChatResponse("Waiting for a partner..."))
	select {
	case partner <- c:
		// now handled by the other goroutine
	case p := <-partner:
		fmt.Println("start partner")
		chat(p, c)
	case <-time.After(5 * time.Second):
		fmt.Println("start bot chatting")
		chat(Bot(), c)
	}
}

func chat(a, b io.ReadWriteCloser) {
	log.Println("start chat")
	fmt.Fprintln(a, NewServerChatResponse("Found one! Say hi."))
	fmt.Fprintln(b, NewServerChatResponse("Found one! Say hi."))
	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	log.Println("close chat")
	a.Close()
	b.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}
