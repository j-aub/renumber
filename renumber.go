package main

/*
import "net/http/pprof"
import "runtime/pprof"
import "os/signal"
*/

import (
	_ "embed"
	"io"
	"log"
	"net/http"
	"os/exec"
)

//go:embed template/header.html
var header string

//go:embed template/footer.html
var footer string

func renumberHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, header)
	defer io.WriteString(w, footer)

	if r.Method == "POST" {
		// parse the form as a stream to avoid reading it all at once
		r, err := r.MultipartReader()
		if err != nil {
			io.WriteString(w, "Bad form submission")
			return
		}

		// we assume the first part is the correct part
		p, err := r.NextPart()
		// first part is absent
		if err == io.EOF {
			io.WriteString(w, "Bad form submission")
			return
		}
		// general error
		if err != nil {
			io.WriteString(w, "Bad form submission")
			return
		}

		cmd := exec.Command("gawk", "{if(!/[0-9]+\\./) $0=\"0. \" $0; line++; sub(/^[0-9]+/,line)}1")
		// we feed the form to awk stdin
		stdin, err := cmd.StdinPipe()
		if err != nil {
			io.WriteString(w, "Failed to setup stdin pipe")
			return
		}
		// and read the renumbered result from awk's stdout
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			io.WriteString(w, "Failed to setup stdout pipe")
			return
		}

		// start awk
		if err := cmd.Start(); err != nil {
			io.WriteString(w, "Failed to start awk")
			return
		}

		go func() {
			// copy form body to awk stdin
			defer stdin.Close()
			if _, err := io.Copy(stdin, p); err != nil {
				io.WriteString(w, "Failed to copy form data to awk")
			}
		}()

		// copy awk stdout to response
		if _, err := io.Copy(w, stdout); err != nil {
			io.WriteString(w, "Failed to read renumbered list")
			return
		}

		// wait for awk to finish
		if err := cmd.Wait(); err != nil {
			io.WriteString(w, "Failed to wait for awk to finish")
			return
		}

	}
}

func main() {
	/*
	f, err := os.Create("./profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		pprof.StopCPUProfile()
		os.Exit(1)
	}()
	*/
	config := LoadConfig()

	http.HandleFunc("/", renumberHandler)
	// http.HandleFunc("/submit", renumberHandler)

	log.Printf("Listening on %v", config.Address)
	log.Fatal(http.ListenAndServe(config.Address, nil))
}
