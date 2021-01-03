package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

/*
	Input: cat a list of URLs via stdin

	Aim: GET the page and check the response. Are we raw source with an extension we care about? (PHP to start with)

	Reason: First paying bounty was a server PHP source leak. Server wasn't interpreting and returned the file (LFI type result with no LFI needed).
*/

var out io.Writer = os.Stdout

func main() {
	/*
		sc := bufio.NewScanner(os.Stdin)
		urls := []string{}

		for sc.Scan() {
			domain := strings.ToLower(sc.Text())

			if domain != "" && len(domain) > 0 {
				urls = append(urls, domain)
			}
		}
	*/

	var outputFileFlag string
	flag.StringVar(&outputFileFlag, "o", "", "Output file for identified leakd source")
	quietModeFlag := flag.Bool("q", false, "Only output the URL's with leaked source")
	flag.Parse()

	quietMode := *quietModeFlag
	saveOutput := outputFileFlag != ""
	outputToSave := []string{}

	if !quietMode {
		banner()
		fmt.Println("")
	}

	writer := bufio.NewWriter(out)
	urls := make(chan string, 1)
	var wg sync.WaitGroup

	ch := readStdin()
	go func() {
		//translate stdin channel to domains channel
		for u := range ch {
			urls <- u
		}
		close(urls)
	}()

	// flush to writer periodically
	t := time.NewTicker(time.Millisecond * 500)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				writer.Flush()
			}
		}
	}()

	for u := range urls {
		wg.Add(1)
		go func(site string) {
			defer wg.Done()
			finalUrls := []string{}

			// If the identified URL has neither http or https infront of it. Create both and scan them.
			if !strings.Contains(u, "http://") && !strings.Contains(u, "https://") {
				finalUrls = append(finalUrls, "http://"+u)
				finalUrls = append(finalUrls, "https://"+u)
			} else if strings.Contains(u, "http://") {
				finalUrls = append(finalUrls, "https://"+u)
			} else if strings.Contains(u, "https://") {
				finalUrls = append(finalUrls, "http://"+u)
			} else {
				// else, just scan the submitted one as it has either protocol
				finalUrls = append(finalUrls, u)
			}

			// now loop the slice of finalUrls (either submitted OR 2 urls with http/https appended to them)
			for _, uu := range finalUrls {
				leaks, ext := makeRequest(uu, quietMode)
				if leaks {
					// if we had a leak, let the user know
					fmt.Printf("[%s] %s\n", ext, uu)

					if saveOutput {
						outputToSave = append(outputToSave, uu)
					}
				}
			}
		}(u)
	}

	wg.Wait()

	// just in case anything is still in buffer
	writer.Flush()

	if saveOutput {
		file, err := os.OpenFile(outputFileFlag, os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil && !quietMode {
			log.Fatalf("failed creating file: %s", err)
		}

		datawriter := bufio.NewWriter(file)

		for _, data := range outputToSave {
			_, _ = datawriter.WriteString(data + "\n")
		}

		datawriter.Flush()
		file.Close()
	}
}

func banner() {
	fmt.Println("---------------------------------------------------")
	fmt.Println("LeakyTap -> Crawl3r")
	fmt.Println("List URL's which appear to be leaking source instead of having the server interpret it")
	fmt.Printf("Currently looks for:\n\tphp\n\n")
	fmt.Println("Run again with -q for cleaner output")
	fmt.Println("---------------------------------------------------")
}

func readStdin() <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			url := strings.ToLower(sc.Text())
			if url != "" {
				lines <- url
			}
		}
	}()
	return lines
}

func makeRequest(url string, quietMode bool) (bool, string) {
	targetExtension := getEndpointFileExtension(url)
	if targetExtension == "" {
		if !quietMode {
			fmt.Println("No extension identified for", url, " - skipping")
		}
		return false, ""
	}

	resp, err := http.Get(url)
	if err != nil {
		if !quietMode {
			fmt.Println("[error] performing the request to:", url)
		}
		return false, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if !quietMode {
				fmt.Println("[error] reading response bytes from:", url)
			}
			return false, ""
		}
		bodyString := string(bodyBytes)
		return parseBodyForSource(targetExtension, bodyString, quietMode), targetExtension
	} else {
		return false, ""
	}
}

func parseBodyForSource(language string, body string, quietMode bool) bool {
	lines := strings.Split(body, "\n")
	if len(lines) == 0 {
		if !quietMode {
			fmt.Println("Empty response body during parse")
		}
		return false
	}

	// logic here based on the identified language. Should be easy enough to extend
	if language == "php" {
		return strings.Contains(lines[0], "<?php") || strings.Contains(lines[len(lines)-1], "?>")
	}

	return false
}

func getEndpointFileExtension(url string) string {
	splitByParam := strings.Split(url, "?")                       // from www.url.com/lol/lol.php?lol=lol -> "www.url.com/lol/lol.php","lol=lol"
	endpointExtensionSplit := strings.Split(splitByParam[0], ".") // from www.url.com/lol/lol.php -> "www.url.com/lol/lol","php"

	if len(endpointExtensionSplit) > 1 {
		// working extensions:
		ep := endpointExtensionSplit[len(endpointExtensionSplit)-1]

		// only return an extension if it's legal. This will skip the request if not.
		if ep == "php" {
			return ep
		}

		return ""
	} else {
		return ""
	}
}
