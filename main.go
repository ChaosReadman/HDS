package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/comail/colog"
	//	"github.com/pkg/profile"
)

func controll(w http.ResponseWriter, req *http.Request) {
	//	defer profile.Start(profile.ProfilePath(".")).Stop()
	defer req.Body.Close()

	method := req.Method
	log.Println("debug: [method] " + method)
	if method == "OPTIONS" {
		// preflight処理
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Vary", "Access-Control-Request-Method")
		w.Header().Set("Vary", "Access-Control-Request-Headers")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{status: 200, result: 'ok'}")
		return
	}

	urlPath := pwd + "/serverRoot" + strings.ToLower(req.URL.Path) + "/retVal.json"

	// For Debug
	for k, v := range req.Header {
		log.Printf("debug: [header]" + k)
		log.Println("debug: " + strings.Join(v, ","))
	}

	// Open and Read retVal.json File
	log.Println(("debug: Read retVal.json from " + urlPath))
	retVal, error := os.ReadFile(urlPath)

	if error != nil {
		log.Println("error: error on Read ", urlPath)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		http.Error(w, "no ret file:"+urlPath, http.StatusNotFound)
		return
	}

	retValStr := string(retVal)
	log.Println("info: ResponsWriter Start")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	w.WriteHeader(http.StatusOK)
	log.Println("info: Fprint response start")
	fmt.Fprint(w, retValStr)
	log.Println("info: Fprint response end")
	log.Println("info: ResponsWriter End")
}

var pwd string

func main() {
	pwd, _ = os.Getwd()
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LWarning)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()

	log.Println("info: Server Start....")
	http.HandleFunc("/", controll)
	http.ListenAndServe(":5000", nil)
}

//curl -i POST '127.0.0.1:5000/recipe/newpage' -H 'Content-Type: application/json' -H  -d '{"title":"タイトル７","page":"ぺーじ"}'
