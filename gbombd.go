package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

/* defaults and configs */
const (
	default_fname = "bomb.gz"
	default_portn = 80
	default_sizen = 10240
)

var (
	bomb     *os.File
	verbose  bool
	preserve bool
	fname    string
	portn    int
	sizen    int
    // we want the scanners/bots to recognize
    // this common files
	handle   = []string{"/", "/robots.txt", "/admin", "/.htpasswd"}
	useragents = []string{ // user agent list!
		"(hydra)", ".nasl", "absinthe", "advanced email extractor",
		"arachni/", "autogetcontent", "bilbo", "BFAC",
		"brutus", "brutus/aet", "bsqlbf", "cgichk",
		"cisco-torch", "commix", "core-project/1.0", "crimscanner/",
		"datacha0s", "dirbuster", "domino hunter", "dotdotpwn",
		"email extractor", "fhscan core 1.", "floodgate", "get-minimal",
		"gootkit auto-rooter scanner", "grabber", "grendel-scan", "havij",
		"inspath", "internet ninja", "jaascois", "zmeu", "masscan", "metis", "morfeus fucking scanner",
		"mysqloit", "n-stealth", "nessus", "netsparker",
		"Nikto", "nmap nse", "nmap scripting engine", "nmap-nse",
		"nsauditor", "openvas", "pangolin", "paros",
		"pmafind", "prog.customcrawler", "qualys was", "s.t.a.l.k.e.r.",
		"security scan", "springenwerk", "sql power injector", "sqlmap",
		"sqlninja", "teh forest lobster", "this is an exploit", "toata dragostea",
		"toata dragostea mea pentru diavola", "uil2pn", "user-agent:", "vega/",
		"voideye", "w3af.sf.net", "w3af.sourceforge.net", "w3af.org",
		"webbandit", "webinspect", "webshag", "webtrends security analyzer",
		"webvulnscan", "whatweb", "whcc/", "wordpress hash grabber",
		"xmlrpc exploit", "WPScan", "curl", "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; FDM; .NET CLR 2.0.50727; InfoPath.2; .NET CLR 1.1.4322)",
	}
)

/***/

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Be verbose")
	flag.BoolVar(&preserve, "preserve", false, "Preserve the bomb file for future use")
	flag.StringVar(&fname, "filename", default_fname, "Filename to create and use")
	flag.IntVar(&sizen, "size", default_sizen, "Size of file to create")
	flag.IntVar(&portn, "port", default_portn, "HTTPd port to use")
	flag.Parse()
}

func remove_bomb(f string) error {
	if preserve {
		return nil
	}
	if verbose {
		fmt.Printf("[+] Removing bomb file: %s\n", f)
	}
	err := os.Remove(f)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func blast_into_oblivion(w http.ResponseWriter, r *http.Request) {
	var ua string = r.UserAgent()

	for j := 0; j < len(useragents); j++ {
		if ua == "" || strings.Contains(ua, useragents[j]) {
			fi, _ := bomb.Stat()
			fsize := fi.Size()
			if verbose {
                fmt.Printf("[+] Serving %d bytes gzipped to %s (UA:\"%s\")\n",
                    fsize, r.RemoteAddr, useragents[j])
			}

			w.Header().Add("Content-Encoding", "gzip")
			w.Header().Add("Content-type", "text/html; charset=UTF-8")
			w.Header().Add("Content-Length", fmt.Sprintf("%d", fsize))
			w.WriteHeader(http.StatusOK)

			bomb.Seek(0, os.SEEK_SET)
			var o io.Reader = bomb
			_, err := io.CopyN(w, o, int64(fsize))
			if err != nil {
				fmt.Fprintf(os.Stderr, "[-] Error writing bomb (\"%s\")\n", err)
				return
			}

			return
		}
	}
	w.WriteHeader(http.StatusForbidden)
}

func main() {
	var err error
	c_sig := make(chan os.Signal, 1)
	signal.Notify(c_sig, os.Interrupt, os.Kill)
	defer bomb.Close()
	defer remove_bomb(fname)

	// let's create the bomb
	bomb, err = os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal(err)
	}

	s, _ := bomb.Stat()
	if s.Size() <= 0 {
		fmt.Printf("[+] Creating file: %s. Wait...\n", fname)

		// use /usr/bin/dd and /us/bin/gzip to create content faster and realiable
		// piping output to file
		fill_cmd := fmt.Sprintf("(echo '<html>' && dd if=/dev/zero bs=1M count=%d) | gzip -5 -c --", sizen)
		cmd := exec.Command("sh", "-c", fill_cmd)
		cmd.Stdout = bomb
		cmd_err, _ := cmd.StderrPipe()
		err = cmd.Start()
		if err != nil {
			error_out, _ := ioutil.ReadAll(cmd_err)
			println(error_out)
			log.Fatal(err)
		}
		err = cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
		println("[+] Done!")
	} else {
		println("[i] Already exists, using that file.")
	}

	for i := 0; i < len(handle); i++ {
		http.HandleFunc(handle[i], blast_into_oblivion)
	}

	go func() {
        if verbose {
            fmt.Printf("[+] Serving data on port %d/TCP\n", portn)
        }
		if err := http.ListenAndServe(fmt.Sprintf(":%d", portn), nil); err != nil {
			fmt.Fprintf(os.Stderr, "[-] Error creating httpd server (\"%s\")\n", err)
			os.Exit(1)
		}
	}()
	sig := <-c_sig
	fmt.Printf("\r[i] Catched signal: %s. Exiting...\n", sig.String())

}
