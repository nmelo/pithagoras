package wifi

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Wifi struct {
	Name string
}

func GetList() []Wifi {
	cmd := exec.Command("iw", "dev", "wlan0", "scan")
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wifis: %q\n", out.String())
	return nil
}
