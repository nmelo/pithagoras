package ui

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"runtime"

	"github.com/nmelo/pithagoras/pkg/db"
	"github.com/nmelo/pithagoras/pkg/wifi"
	"github.com/skycoin/skycoin/src/aether/wifi"
)

type Server struct {
	Port uint
}

func Serve(ctx context.Context) {
	fmt.Println("Serving UI...")

	http.HandleFunc("/sessions", handleSessions)
	http.HandleFunc("/wifis", handleWifis)

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
