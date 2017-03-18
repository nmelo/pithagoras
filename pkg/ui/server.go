package ui

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"runtime"

	"github.com/Pallinder/go-randomdata"
	"github.com/nmelo/pithagoras/pkg/db"
	"github.com/nmelo/pithagoras/pkg/wifi"
	"github.com/skycoin/skycoin/src/aether/wifi"
)

var (
	names = map[string]string{}
)

func Serve(ctx context.Context) {
	fmt.Println("Serving UI...")

	http.HandleFunc("/sessions", handleSessions)
	http.HandleFunc("/wifis", handleWifis)
	http.HandleFunc("/chat", handleChat)
	http.HandleFunc("/chat/clear", handleClearChat)

	var addr string
	if runtime.GOOS == "linux" {
		addr = ":80"
	} else {
		addr = ":8080"
	}
	go http.ListenAndServe(addr, nil)

	select {
	case <-ctx.Done():
		fmt.Println("Stoping UI...")
		return
	}
}

var sessionsTemplate = template.Must(template.New("sessions").Parse(`
<html>
<head/>
<body>
  <form action="/sessions" method="post">
	 <input type="submit" value="Clear">
  </form>

  <span>Sessions:</span>
  <ol>
  {{range .Sessions}}
    <li>{{.Key}} - {{.Date}}</li>
  {{else}}
  	No sessions
  {{end}}
  </ol>
</body>
</html>
`))

func handleSessions(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		err := db.ClearSessions()
		if err != nil {
			http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
			return
		}
		req.Method = http.MethodGet
		http.Redirect(w, req, "/sessions", http.StatusFound)
		return
	}

	sessions, err := db.ListSessions()
	if err != nil {
		http.Error(w, "failed to read sessions", http.StatusInternalServerError)
		return
	}

	if err := sessionsTemplate.Execute(w, struct {
		Sessions []db.Session
	}{
		Sessions: sessions,
	}); err != nil {
		fmt.Println(err)
		return
	}

}

var wifisTemplate = template.Must(template.New("wifis").Parse(`
<html>
<head/>
<body>
  <ol>
  {{range .Wifis}}
    <li>{{.ESSID}}</li>
  {{end}}
  </ol>
</body>
</html>
`))

func handleWifis(w http.ResponseWriter, req *http.Request) {

	nws, err := wifi.Scan()
	if err != nil {
		http.Error(w, "failed to scan", http.StatusInternalServerError)
		return
	}
	if err := wifisTemplate.Execute(w, struct {
		Wifis []network.WifiNetwork
	}{
		Wifis: nws,
	}); err != nil {
		fmt.Println(err)
		return
	}

}

var chatTemplate = template.Must(template.New("sessions").Parse(`
<html>
<head/>
<body>
  <form action="/chat/clear" method="post">
	 <input type="submit" value="Clear">
  </form>

  <form action="/chat" method="post">
	 Chat:
	  <input type="text" name="message">
	  <input type="submit" value="Send">
  </form>

  <span>Chats:</span>
  <ul>
  {{range .Messages}}
    <li>{{.Username}} - {{.Text}}</li>
  {{else}}
  	No messages
  {{end}}
  </ul>
</body>
</html>
`))

func handleChat(w http.ResponseWriter, req *http.Request) {
	var sessionID string
	c, err := req.Cookie("user")
	if err == http.ErrNoCookie {
		sessionID = GenerateID()
		names[sessionID] = randomdata.SillyName()
		fmt.Printf("Name: %s\n", names[sessionID])

		c = &http.Cookie{Name: "user", Value: sessionID, Path: "/chat"}
		http.SetCookie(w, c)
	} else {
		sessionID = c.Value
	}

	if req.Method == http.MethodPost {
		fmt.Print("Post")
		err := db.AddMessage(sessionID, names[sessionID], req.FormValue("message"))
		if err != nil {
			http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
			return
		}
		req.Method = http.MethodGet
		http.Redirect(w, req, "/chat", http.StatusFound)
		return
	}

	messages, err := db.ListMessages()
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
		return
	}

	if err := chatTemplate.Execute(w, struct {
		Messages []db.Message
	}{
		Messages: messages,
	}); err != nil {
		fmt.Println(err)
		return
	}

}

func handleClearChat(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		fmt.Print("Delete")
		err := db.ClearMessages()
		if err != nil {
			http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
			return
		}
		req.Method = http.MethodGet
		http.Redirect(w, req, "/chat", http.StatusFound)
		return
	}
}

func GenerateID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
