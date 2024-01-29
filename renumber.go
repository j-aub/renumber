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

// take all data from reader, renumber it, write it to writer
func renumber(w io.Writer, r io.Reader) {
	cmd := exec.Command("gawk", "{if(!/^[0-9]+\\./) $0=\"0. \" $0; line++; sub(/^[0-9]+/,line)}1")
	// we'll feed the form to awk's stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Print("Failed to setup stdin pipe")
		return
	}
	// and read the renumbered result from awk's stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Print("Failed to setup stdout pipe")
		return
	}

	// start awk
	if err := cmd.Start(); err != nil {
		log.Print("Failed to start awk")
		return
	}

	// if anything pipe related fails from here on now we don't want a
	// zombie process
	defer cmd.Wait()

	go func() {
		// copy form body to awk stdin
		defer stdin.Close()
		if _, err := io.Copy(stdin, r); err != nil {
			log.Print(w, "Failed to copy form data to awk")
		}
	}()

	// copy awk stdout to response
	if _, err := io.Copy(w, stdout); err != nil {
		log.Print(w, "Failed to read renumbered list")
		return
	}

	// wait for awk to finish
	if err := cmd.Wait(); err != nil {
		log.Print(w, "Failed to wait for awk to finish")
		return
	}
}

func renumberHandler(w http.ResponseWriter, r *http.Request) {
	// indicate that we're doing full duplex http
	if err := http.NewResponseController(w).EnableFullDuplex(); err != nil {
		log.Print("Failed to enable full duplex for ResponseWriter. Panicking")
		panic("Full duplex is not evailable when it should be")
	}

	if r.Method == "GET" {
		io.WriteString(w, header)
		io.WriteString(w, footer)
	} else if r.Method == "POST" {
		// parse the form as a stream to avoid reading it all at once
		// if anything with the form parsing goes wrong it's the
		// user's fault so we return a 400 status
		r, err := r.MultipartReader()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// we assume the first part is the correct part
		p, err := r.NextPart()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		io.WriteString(w, header)
		renumber(w, p)
		io.WriteString(w, footer)
	} else {
		// we don't allow anything but GET and POST
		w.WriteHeader(http.StatusBadRequest)
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
