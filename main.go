// Small tool to test supported TLS versions by a URL
//   - This will test for TLS versions 1.0, 1.1, 1.2 and 1.3
//   - Supports proxy through environment variable "http_proxy"
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// TestTLSVersion attempts to connect to the given URL using the provided TLS configuration.
func TestTLSVersion(url string, tlsConfig *tls.Config) error {
	// Create an HTTP client with the custom TLS configuration

	// if http_proxy is found in environment, client will use that proxy
	_, useProxy := os.LookupEnv("http_proxy")
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	if useProxy {
		transport.Proxy = http.ProxyFromEnvironment
	}

	// client times out after 5 seconds
	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and discard the response body to avoid leaking resources
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	var url string
	var skipVerification bool

	flag.StringVar(&url, "u", "", "URL to test")
	flag.BoolVar(&skipVerification, "s", false, "Skip certificate verification (optional)")
	flag.Parse()

	if url == "" {
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s [-s] -u <URL>\n\tSmall tool to test supported TLS versions by a URL\nFlags:\n\t-s\t: Skip certificate verification (optional)\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	if !strings.Contains(url, "http") {
		fmt.Fprintf(
			os.Stderr,
			"Please include http/https in the URL. Only http/https protocols are supported.\n",
		)
		os.Exit(1)
	}

	tlsVersions := map[string]uint16{
		"TLS 1.0": tls.VersionTLS10,
		"TLS 1.1": tls.VersionTLS11,
		"TLS 1.2": tls.VersionTLS12,
		"TLS 1.3": tls.VersionTLS13,
	}

	fmt.Printf("Testing supported TLS versions for %s\n", url)

	for name, version := range tlsVersions {
		tlsConfig := &tls.Config{
			MinVersion: version,
			MaxVersion: version,
		}

		if skipVerification {
			tlsConfig.InsecureSkipVerify = skipVerification
		}

		err := TestTLSVersion(url, tlsConfig)
		if err != nil {
			fmt.Printf("  %s is not supported: %v\n", name, err)
		} else {
			fmt.Printf("  %s is supported\n", name)
		}
	}
}
