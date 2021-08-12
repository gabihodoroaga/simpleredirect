package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"os"
)

var (
	redirects map[string]redirectEntry = make(map[string]redirectEntry)
	listen    string
)

func main() {
	
	err := parseFlags(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:])
	if err != nil {
		log.Fatal("readSettings error:", err)
	}
	http.HandleFunc("/", redirect) 
	fmt.Printf("starting server on %s\n", listen)	
	err = http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {

	if r.RequestURI == "/hc" {
		w.WriteHeader(http.StatusOK)
		return
	}

	parts := strings.Split(r.Host, ".")
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]

	v, found := redirects[domain]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	url := v.url
	if v.scheme != "" {
		url = v.scheme + "://" + url
	} else {
		url = r.URL.Scheme + "://" + url
	}
	http.Redirect(w, r, url, v.code)
}



func parseFlags(commandLine *flag.FlagSet, args []string) error {
	// setup flags
	var redirectFlag arrayFlags
	var listenFlag string
	commandLine.Var(&redirectFlag, "redirect", "Redirect map source:destination:[code]. Example: test.com:none.com:301 or test com:https://none.com. [code] is optional and defaults to 301")
	commandLine.StringVar(&listenFlag, "listen", ":8080", "Listen host and port. Example: :80, 0.0.0.0:8080, localhost:8080. ")
	commandLine.Parse(args)

	// validate listen
	re := regexp.MustCompile(`^[a-zA-Z0-9_\.-]*:\d{2,5}$`)
	if !re.MatchString(listenFlag) {
		return fmt.Errorf("invalid value for listen flag. Expected [host]:port, got %s", listenFlag)
	}
	listen = listenFlag

	// validate redirects
	regex := regexp.MustCompile(`^(\w*\.\w*):((http|https):\/\/)?((\w+\.)?(\w+\.)?\w+\.\w+(/.*?)?)(:([0-9]{3}))?$`)
	for _, v := range redirectFlag {
		match := regex.FindStringSubmatch(v)
		/*  This is the result of the regex
		Group 0, value: domain.com:http://vm1.dev.test.com/test:301
		Group 1, value: domain.com
		Group 2, value: http://
		Group 3, value: http
		Group 4, value: vm1.dev.test.com/test
		Group 5, value: vm1.
		Group 6, value: dev.
		Group 7, value: /test
		Group 8, value: :301
		Group 9, value: 301
		*/
		if match != nil {
			code := 301
			if match[9] != "" {
				parsedCode, err := strconv.Atoi(match[9])
				if err != nil {
					return fmt.Errorf("invalid value for code. Cannot convert %s to int", match[9])
				}
				code = parsedCode
			}
			from := match[1]
			to := redirectEntry{
				scheme: match[3],
				url:    match[4],
				code:   code,
			}
			redirects[from] = to
			fmt.Printf("redirect %s => %s with code %d and scheme %s\n", from, to.url, to.code, to.scheme)
		} else {
			return fmt.Errorf("invalid value for redirect. Must match this regex %s", regex.String())
		}
	}



	return nil
}

type redirectEntry struct {
	scheme string
	url    string
	code   int
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
