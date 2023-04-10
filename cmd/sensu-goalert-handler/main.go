package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"
)

var detailsTmpl = template.Must(template.New("main").Parse(
	"**Timestamp**\n: {{.Event.ISOTimestamp}}\n\n" +
		"**Command**\n\n" +
		"```bash\n{{.Event.Check.Command}}\n\n# Output\n{{.Event.Check.Output}}\n```\n\n" +
		"**Payload**\n\n" +
		"```json\n{{.JSON}}\n```\n"))

type meta struct {
	Name       string
	Namespace  string
	GoAlertURL string `json:"goalert_url"`
}

func (m meta) String() string {
	if m.Namespace == "default" {
		return m.Name
	}
	return fmt.Sprintf("%s.%s", m.Namespace, m.Name)
}

type event struct {
	Timestamp int64
	Entity    struct {
		Metadata meta
	}
	Check struct {
		Metadata meta
		Command  string
		Output   string
		State    string
	}
}

func (e event) ISOTimestamp() string {
	return time.Unix(e.Timestamp, 0).Format(time.RFC3339)
}

func main() {
	log.SetFlags(log.Lshortfile)

	urlStr := flag.String("url", os.Getenv("GOALERT_URL"), "")
	flag.Parse()

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("read stdin:", err)
	}

	var e event
	err = json.Unmarshal(data, &e)
	if err != nil {
		log.Fatal("parse input:", err)
	}

	if e.Check.State == "" {
		log.Fatal("event does not contain a check")
	}

	if e.Entity.Metadata.GoAlertURL != "" {
		*urlStr = e.Entity.Metadata.GoAlertURL
	}

	if e.Check.Metadata.GoAlertURL != "" {
		*urlStr = e.Check.Metadata.GoAlertURL
	}

	if *urlStr == "" {
		log.Fatal("url is required")
	}

	dataBuf := new(bytes.Buffer)
	err = json.Indent(dataBuf, data, "", "  ")
	if err != nil {
		log.Fatal("format event data:", err)
	}
	var renderData struct {
		Event event
		JSON  string
	}
	renderData.JSON = dataBuf.String()
	renderData.Event = e

	dataBuf.Reset()
	err = detailsTmpl.Execute(dataBuf, renderData)
	if err != nil {
		log.Fatal("render alert details:", err)
	}

	alert := make(url.Values)
	if e.Check.State == "passing" {
		alert.Set("action", "close")
	}
	key := e.Entity.Metadata.String() + "/" + e.Check.Metadata.String()
	alert.Set("dedup", key)
	alert.Set("summary", key+": "+e.Check.Output)
	alert.Set("details", dataBuf.String())

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("POST", *urlStr, strings.NewReader(alert.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		req = req.WithContext(ctx)
		resp, err := http.DefaultClient.Do(req)
		cancel()
		if err != nil {
			log.Printf("ERROR: update %s: %v", key, err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode == 204 {
			// done
			return
		}
		log.Printf("ERROR: update %s: unexpected response '%s'; want 204", key, resp.Status)
	}

	log.Fatalf("failed to update state of %s", key)
}
