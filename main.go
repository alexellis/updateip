package main

import (
	b64 "encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/alexellis/updateip/cmd"
	"gopkg.in/yaml.v2"
)

func main() {
	var configFile string
	var version bool
	flag.StringVar(&configFile, "config", "./config.yaml", "Config file for domains")
	flag.BoolVar(&version, "version", false, "Print version")
	flag.Parse()

	if version {
		cmd.PrintupdateipASCIIArt()
		return
	}

	if err := runE(configFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func runE(configFile string) error {
	config, err := readConfig(configFile)
	if err != nil {
		return err
	}

	for _, domain := range config.Domains {
		log.Printf("%v %v\n", domain, domain.PlainPassword())

		if domain.Provider != "namecheap" {
			return fmt.Errorf("provider %s not yet supported", domain.Provider)
		}

		// namecheap returns XML and a 200 code even if
		// there is an error with the update.
		res, err := updateNamecheap(domain)
		if err != nil {
			return err
		}

		log.Printf("%s result: %s", domain, res)

	}

	return nil
}

func updateNamecheap(domain Domain) (string, error) {

	hostIndex := strings.Index(domain.Domain, ".")
	if hostIndex == -1 {
		return "", fmt.Errorf("invalid sub-domain %s", domain.Domain)
	}

	host := domain.Domain[:hostIndex]
	domainName := domain.Domain[hostIndex+1:]

	ip := ""
	if domain.IP == "external" {
		ipv, err := GetExternalIP()
		if err != nil {
			return "", err
		}
		ip = ipv
	}

	uri := fmt.Sprintf("https://dynamicdns.park-your-domain.com/update?host=%s&domain=%s&password=%s&ip=%s",
		host, domainName, domain.PlainPassword(), ip)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting %s: %s", uri, err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting %s: %s, error: %s", uri, res.Status, string(body))
	}

	return string(body), nil
}

func readConfig(configFile string) (Config, error) {
	var config Config

	fileData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	if err := yaml.Unmarshal(fileData, &config); err != nil {
		return config, err
	}

	return config, nil
}

type Config struct {
	Domains []Domain `yaml:"domains"`
}

type Domain struct {
	Domain   string `yaml:"domain"`
	IP       string `yaml:"ip"`
	Password string `yaml:"password"`
	Provider string `yaml:"provider"`
}

func (d Domain) PlainPassword() string {
	sDec, _ := b64.StdEncoding.DecodeString(d.Password)
	return string(sDec)
}

// GetExternalIP uses https://checkip.amazonaws.com to determine the
// external IP address of the host, or returns an error.
func GetExternalIP() (string, error) {
	lookupURL := "https://checkip.amazonaws.com"

	req, err := http.NewRequest(http.MethodGet, lookupURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-agent", "updateip (github.com/alexellis)")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error with request to %s, %w", lookupURL, err)
	}

	var body []byte
	if res.Body != nil {
		defer res.Body.Close()

		body, _ = ioutil.ReadAll(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d, from %s, error: %s", res.StatusCode, lookupURL, string(body))
	}

	s := strings.TrimSpace(string(body))
	if v := net.ParseIP(s); v == nil {
		return "", fmt.Errorf("%s was not a valid IP", s)
	}

	return s, nil
}
