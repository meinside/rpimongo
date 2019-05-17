package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"

	"github.com/meinside/rpi-tools/status"
)

// constants
const (
	configFilename = "config.json"

	StaticDirname      = "static"
	TimeoutSeconds     = 10
	IdleTimeoutSeconds = 60

	CacheDirname = "./acme"

	DefaultPortHTTP  = 80
	DefaultPortHTTPS = 443

	DefaultPageTitle = "RPiMonGo: Raspberry Pi Monitoring with Go"
)

// config is a struct for config file
type config struct {
	Title            string   `json:"title"`
	Hostname         string   `json:"hostname"`
	ServeSSL         bool     `json:"serve_ssl,omitempty"`
	PortHTTP         int      `json:"port_http,omitempty"`
	PortHTTPS        int      `json:"port_https,omitempty"`
	RedactedKeywords []string `json:"redacted_keywords,omitempty"`
	Verbose          bool     `json:"verbose,omitempty"`
}

// apiResult is a struct for json api result
type apiResult struct {
	Result string `json:"result"`
	Value  string `json:"value"`
}

// Load templates
var templates = template.Must(template.ParseFiles(
	"tpl/layout.html",
	"tpl/index.html",
	"tpl/links.html",
))

// Read system values with rpi-tools
func readValue(method string, redactedKeywords []string) (result string, err error) {
	switch method {
	case "hostname": // hostname
		result, err = status.Hostname()
	case "uname": // uname -a
		result, err = status.Uname()
	case "uptime": // uptime
		result, err = status.Uptime()
	case "free_spaces": // df -h
		result, err = status.FreeSpaces()
	case "memory_split": // vcgencmd get_mem arm && vcgencmd get_mem gpu
		var splits []string
		splits, err = status.MemorySplit()
		result = strings.Join(splits, "\n")
	case "free_memory": // free -h
		result, err = status.FreeMemory()
	case "cpu_temperature": // vcgencmd measure_temp
		result, err = status.CpuTemperature()
	case "cpu_info": //cat /proc/cpuinfo
		result, err = status.CpuInfo()
	default:
		result = "Error"
		err = fmt.Errorf("No such method: %s", method)
	}

	// redact keywords
	if err == nil {
		result = redact(result, redactedKeywords)
	}

	return
}

// redact given string
func redact(str string, keywords []string) string {
	for _, k := range keywords {
		str = strings.Replace(str, k, "*redacted*", -1)
	}

	return str
}

// Read config file
func readConfig() (conf config, err error) {
	var execFilepath string
	if execFilepath, err = os.Executable(); err == nil {
		var file []byte
		if file, err = ioutil.ReadFile(filepath.Join(filepath.Dir(execFilepath), configFilename)); err == nil {
			var conf config
			if err = json.Unmarshal(file, &conf); err == nil {
				if conf.Title == "" {
					conf.Title = DefaultPageTitle
				}
				return conf, nil
			}
		}
	}

	return config{}, err
}

// Render html template
func renderTemplate(w http.ResponseWriter, tmplName string, conf config) {
	w.Header().Set("Content-Type", "text/html")

	buffer := new(bytes.Buffer)
	if err := templates.ExecuteTemplate(buffer, tmplName, struct{}{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		if err := templates.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			"Title":   conf.Title,
			"Content": template.HTML(buffer.String()),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Render json api result
func renderAPIResult(w http.ResponseWriter, actionName string, conf config) {
	w.Header().Set("Content-Type", "application/json")

	if result, err := readValue(actionName, conf.RedactedKeywords); err == nil {
		json.NewEncoder(w).Encode(apiResult{
			Result: "ok",
			Value:  result,
		})
	} else {
		json.NewEncoder(w).Encode(apiResult{
			Result: "error",
			Value:  err.Error(),
		})
	}
}

func main() {
	if conf, err := readConfig(); err == nil {
		// route rules
		router := mux.NewRouter()
		// /static/
		router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(StaticDirname))))
		// index
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w, "index.html", conf)
		})
		// /links
		router.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w, "links.html", conf)
		})
		// /api/*.json
		router.HandleFunc("/api/{action}.json", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			renderAPIResult(w, vars["action"], conf)
		})

		// start HTTPS server
		var manager *autocert.Manager
		if conf.ServeSSL {
			manager = &autocert.Manager{
				Prompt: autocert.AcceptTOS,
				HostPolicy: func(ctx context.Context, host string) error {
					if host == conf.Hostname {
						return nil
					}
					return fmt.Errorf("acme/autocert: host %s is not allowed", host)
				},
				Cache: autocert.DirCache(CacheDirname),
			}

			port := conf.PortHTTPS
			if port <= 0 {
				port = DefaultPortHTTPS
			}

			server := newServer(port, router)
			server.TLSConfig = &tls.Config{GetCertificate: manager.GetCertificate}

			go func() {
				if conf.Verbose {
					log.Printf("> HTTPS server starts listening...")
				}
				if err := server.ListenAndServeTLS("", ""); err != nil {
					panic(err)
				}
			}()
		}

		// start HTTP server
		if manager == nil {
			port := conf.PortHTTP
			if port <= 0 {
				port = DefaultPortHTTP
			}

			server := newServer(port, router)

			if conf.Verbose {
				log.Printf("> HTTP server starts listening...")
			}

			if err := server.ListenAndServe(); err != nil {
				panic(err)
			}
		} else {
			if conf.Verbose {
				log.Printf("> HTTP server for 'http-01' challenge starts listening...")
			}

			// listening for `http-01` challenge
			if err := http.ListenAndServe(":http", manager.HTTPHandler(nil)); err != nil {
				panic(err)
			}
		}
	} else {
		panic(err)
	}
}

func newServer(port int, router *mux.Router) *http.Server {
	return &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: TimeoutSeconds * time.Second,
		ReadTimeout:  TimeoutSeconds * time.Second,
		IdleTimeout:  IdleTimeoutSeconds * time.Second,
	}
}
