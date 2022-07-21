package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Host        string
	Port        string
	Destination string
}

var client = http.DefaultClient

func main() {

	config := Config{
		Host:        "",
		Port:        "5969",
		Destination: "",
	}

	flag.StringVar(&config.Host, "h", "", "hostname for the rp server, leave empty for 0.0.0.0")
	flag.StringVar(&config.Port, "p", "5969", "port for the rp server, leave empty for 5969")
	flag.StringVar(&config.Destination, "d", "", "port for the rp server, must provide")

	flag.Parse()

	errs := startServer(config)

	for {
		<-errs
		os.Exit(-1)
	}
}

func addressFromHostAndPort(host string, port string) string {
	return fmt.Sprint(host, ":", port)
}

func startServer(c Config) chan error {
	errs := make(chan error)

	initClient()

	go func() {

		mux := getHandler(c.Destination)

		address := addressFromHostAndPort(c.Host, c.Port)
		log.Println("rp | service running on " + address)
		if err := http.ListenAndServe(address, mux); err != nil {
			log.Fatalf("EXITING! %v", err.Error())
			errs <- err
		}
	}()

	return errs
}

func initClient() {
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

func getHandler(destination string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		uri := req.RequestURI
		url := fmt.Sprintf("%s%s", destination, uri)

		defer req.Body.Close()

		r, err := http.NewRequest(req.Method, url, req.Body)
		if err != nil {
			handleError(w, err)
			return
		}

		r.Header = req.Header

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			handleError(w, err)
			return
		}

		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, mvals := range values {
				w.Header().Add(key, mvals)
			}
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		log.Println("rp | url: ", url)
	}
}

func handleError(w http.ResponseWriter, e error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	fmt.Fprintf(w, `{"title":"Error Occured","detail":"%s","status":"500"}`, e.Error())
}
