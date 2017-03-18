package ui

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/nmelo/pithagoras/pkg/db"
)

type Server struct {
	Port uint
}

func Serve(ctx context.Context) {
	fmt.Println("Serving UI...")

	http.HandleFunc("/sessions", handleSessions)
	go http.ListenAndServe(":80", nil)

	select {
	case <-ctx.Done():
		fmt.Println("Stoping UI...")
		return
	}
}

func handleSessions(w http.ResponseWriter, req *http.Request) {

	sessions, err := db.ListSessions()
	if err != nil {
		http.Error(w, "failed to read sessions", http.StatusInternalServerError)
		return
	}

	if err := resultsTemplate.Execute(w, struct {
		Sessions []db.Session
	}{
		Sessions: sessions,
	}); err != nil {
		fmt.Println(err)
		return
	}

}

var resultsTemplate = template.Must(template.New("sessions").Parse(`
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
