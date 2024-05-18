# filter-waf

`filter-waf` is a tool written in Go that checks a list of URLs to determine whether they are protected by a Web Application Firewall (WAF). If a website is not protected by a WAF, the tool prints that URL.

## Features

- Accepts a list of URLs for batch processing.
- Detects the presence of a Web Application Firewall (WAF) on each URL.
- Outputs URLs that are not protected by a WAF.

## Installation

1. Make sure you have [Go](https://golang.org/dl/) installed (version 1.17 or later).
2. Install the tool by executing following command:
   ```sh
   go install github.com/isira-adithya/filter-waf@latest
   ```

## Usage
```
Usage of filter-waf:
  -input string
        Input file containing URLs to check or a single URL to check
  -threads int
        Number of threads to use (default 10)
  -verbose
        Verbose mode

Ex: filter-waf -input urls.txt > unprotected-urls.txt
    filter-waf -input https://example.com -verbose
```

## Disclaimer
This tool is intended for educational purposes and ethical testing only. It should only be used on websites that you own or have explicit permission to test. Unauthorized scanning of websites is illegal and unethical. The authors of this tool are not responsible for any misuse or damage caused by this tool.

## Contributing
Contributions are welcome! Please submit a pull request or open an issue to discuss your ideas.