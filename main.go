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
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"

	"github.com/meinside/rpi-tools/status"
)

const (
	ConfigFilename = "config.json"

	StaticDirname      = "static"
	TimeoutSeconds     = 10
	IdleTimeoutSeconds = 60

	CacheDirname = "./acme"

	PortHttp  = 80
	PortHttps = 443

	DefaultPageTitle = "RPiMonGo: Raspberry Pi Monitoring with Go"
)

// Struct for config file
type Config struct {
	Title    string `json:"title"`
	Hostname string `json:"hostname"`
	ServeSSL bool   `json:"serve_ssl,omitempty"`
	Verbose  bool   `json:"verbose,omitempty"`
}

// Struct for json api result
type ApiResult struct {
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
func readValue(method string) (result string, err error) {
	switch method {
	case "hostname": // hostname
		return status.Hostname()
	case "uname": // uname -a
		return status.Uname()
	case "uptime": // uptime
		return status.Uptime()
	case "free_spaces": // df -h
		return status.FreeSpaces()
	case "memory_split": // vcgencmd get_mem arm && vcgencmd get_mem gpu
		splits, err := status.MemorySplit()
		return strings.Join(splits, "\n"), err
	case "free_memory": // free -o -h
		return status.FreeMemory()
	case "cpu_temperature": // vcgencmd measure_temp
		return status.CpuTemperature()
	case "cpu_info": //cat /proc/cpuinfo
		return status.CpuInfo()
	default:
		return "Error", fmt.Errorf("No such method: %s", method)
	}
}

// Read config file
func readConfig() (conf Config, err error) {
	_, filename, _, _ := runtime.Caller(0) // = __FILE__

	if file, err := ioutil.ReadFile(filepath.Join(path.Dir(filename), ConfigFilename)); err == nil {
		var conf Config
		if err := json.Unmarshal(file, &conf); err == nil {
			if conf.Title == "" {
				conf.Title = DefaultPageTitle
			}
			return conf, nil
		} else {
			return Config{}, err
		}
	} else {
		return Config{}, err
	}
}

// Render html template
func renderTemplate(w http.ResponseWriter, tmplName string, conf Config) {
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
func renderApiResult(w http.ResponseWriter, actionName string) {
	w.Header().Set("Content-Type", "application/json")

	if result, err := readValue(actionName); err == nil {
		json.NewEncoder(w).Encode(ApiResult{
			Result: "ok",
			Value:  result,
		})
	} else {
		json.NewEncoder(w).Encode(ApiResult{
			Result: "error",
			Value:  err.Error(),
		})
	}
}

func main() {
	if conf, err := readConfig(); err == nil {
		// route rules
		router := mux.NewRouter()
		router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(StaticDirname))))
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w, "index.html", conf)
		})
		router.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w, "links.html", conf)
		})
		router.HandleFunc("/api/{action}.json", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			renderApiResult(w, vars["action"])
		})

		// start HTTPS server
		if conf.ServeSSL {
			manager := autocert.Manager{
				Prompt: autocert.AcceptTOS,
				HostPolicy: func(ctx context.Context, host string) error {
					if host == conf.Hostname {
						return nil
					}
					return fmt.Errorf("acme/autocert: host %s is not allowed", host)
				},
				Cache: autocert.DirCache(CacheDirname),
			}
			server := newServer(PortHttps, router)
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
		var server *http.Server
		if conf.ServeSSL {
			// redirect to HTTPS
			router = mux.NewRouter()
			router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				uri := "https://" + r.Host + r.URL.String()
				http.Redirect(w, r, uri, http.StatusFound)
			})
		}
		server = newServer(PortHttp, router)
		if conf.Verbose {
			log.Printf("> HTTP server starts listening...")
		}
		if err := server.ListenAndServe(); err != nil {
			panic(err)
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
