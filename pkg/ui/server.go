package ui

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

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
	go http.ListenAndServe(":80", nil)

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
  <ol>
  {{range .Sessions}}
    <li>{{.Key}} - {{.Date}}</li>
  {{end}}
  </ol>
</body>
</html>
`))

func handleSessions(w http.ResponseWriter, req *http.Request) {

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
