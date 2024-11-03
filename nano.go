package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var regex = make(map[string]string)

func matcher(url string) {
	response := requester(url)
	if response != "" {
		Hach, _ := CreatHashSum(response)
		if !contains(HashList, Hach) {
			HashList = append(HashList, Hach)
			for k, p := range regex {
				rgx := regexp.MustCompile(p)
				if rgx.MatchString(response) {
					match := rgx.FindStringSubmatch(response)
					if len(match) > 0 {
						fmt.Printf("%s  \033[32m  %s : %s \033[00m\n", url, k, match[0])
					}
				}
			}
		}
	}
}

func CreatHashSum(input string) (string, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(input))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func isUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func requester(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

// loadRegexFromFile loads regex patterns from the specified file
func loadRegexFromFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening regex file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			regex[parts[0]] = parts[1]
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading regex file:", err)
	}
}
