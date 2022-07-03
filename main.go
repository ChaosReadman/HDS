package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"crypto/tls"

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

	log.Println("info: ResponsWriter Start")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	w.WriteHeader(http.StatusOK)
	log.Println("info: Fprint response start")
	w.Write(retVal)
	log.Println("info: Fprint response end")
	log.Println("info: ResponsWriter End")
}

var pwd string

func main() {
	pwd, _ = os.Getwd()
	// colog の設定
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LDebug)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()

	/*
		// p12を読む
		// privとkeyを pkcs keystore からデコード
		// デコードできない・・・
		p12_data, err := os.ReadFile("auth/test.p12")
		if err != nil {
			log.Fatal(err)
		}

		pkey, cer, err := pkcs12.Decode(p12_data, "abcdefg")
		if err != nil {
			log.Println(err)
			return
		}
		// 認証情報をtls構造体にセット
		cert := tls.Certificate{}
		cert.PrivateKey = pkey
		cert.Certificate = [][]byte{cer.Raw}
	*/

	cer, err := tls.LoadX509KeyPair("auth/usercert.pem", "auth/userkey.pem")
	if err != nil {
		log.Println(err)
		return
	}

	// ハンドラの設定
	mux := http.NewServeMux()
	mux.HandleFunc("/", controll)

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{cer},
	}

	log.Println("info: Server Start....")

	srv := &http.Server{
		Addr:         ":443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	ListenErr := srv.ListenAndServeTLS()
	if err != nil {
		log.Printf("error : %s", ListenErr)
	}
	log.Printf("info: Server Started")
}

// PS C:\tmp\work\source\HDS> curl.exe --insecure -i POST 'https://localhost:8443/regist?aaa=1' -H 'Content-Type: application/json' -d '{"title":"タイトル７","page":"ぺーじ"}'
// curl: (6) Could not resolve host: POST
// HTTP/1.1 200 OK
// Access-Control-Allow-Headers: *
// Access-Control-Allow-Methods: GET, POST, PUT, DELETE, PATCH, OPTIONS
// Access-Control-Allow-Origin: *
// Content-Type: application/json; charset=utf-8
// Date: Thu, 30 Jun 2022 12:56:44 GMT
// Content-Length: 53

// {
//     "result": "OK",
//     "detail": "Registered"
// }
