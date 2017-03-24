package main

import (
	"bytes"
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

	"github.com/meinside/rpi-tools/status"
)

const (
	ConfigFilename = "config.json"

	StaticDirname  = "static"
	TimeoutSeconds = 10
)

// Struct for config file
type Config struct {
	PortNumber int  `json:"port_number"`
	Verbose    bool `json:"verbose"`
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
			return conf, nil
		} else {
			return Config{}, err
		}
	} else {
		return Config{}, err
	}
}

// Render html template
func renderTemplate(w http.ResponseWriter, tmplName string) {
	w.Header().Set("Content-Type", "text/html")

	buffer := new(bytes.Buffer)
	if err := templates.ExecuteTemplate(buffer, tmplName, struct{}{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		if err := templates.ExecuteTemplate(w, "layout.html", map[string]interface{}{
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
			renderTemplate(w, "index.html")
		})
		router.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w, "links.html")
		})
		router.HandleFunc("/api/{action}.json", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			renderApiResult(w, vars["action"])
		})

		// start server
		if conf.Verbose {
			log.Printf("Listening on port: %d...", conf.PortNumber)
		}
		server := &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf(":%d", conf.PortNumber),
			WriteTimeout: TimeoutSeconds * time.Second,
			ReadTimeout:  TimeoutSeconds * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
}
