// validator.go
// author: Yanzheng SUN (realSunyz)

package main

import (
	"bufio"
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
)

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Subdivision struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type SubdivisionMap map[string][]Subdivision

//go:embed data/*.json
var dataFS embed.FS

func readJSON(filename string, v interface{}) error {
	data, err := dataFS.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func validateGeofeed(geofeedFile string, countries []Country, subdivisions SubdivisionMap) []string {

	countryMap := make(map[string]bool)
	for _, country := range countries {
		countryMap[country.Code] = true
	}

	file, err := os.Open(geofeedFile)
	if err != nil {
		return []string{fmt.Sprintf("Failed to open geofeed file: %v", err)}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	var errors []string

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Skip blank lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Skip comment lines
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		reader := csv.NewReader(strings.NewReader(line))
		reader.FieldsPerRecord = -1
		record, _ := reader.Read()

		// Validate geofeed format (prefix, country_code, subdivision_code, city, postal_code)
		if len(record) < 5 {
			errors = append(errors, fmt.Sprintf("Line %d: invalid "+Magenta+"geofeed "+Reset+"format", lineNumber))
			continue
		}
		prefix, countryCode, subdivisionCode := record[0], record[1], record[2]

		// Validate prefix format
		_, _, prefixErr := net.ParseCIDR(prefix)
		if prefixErr != nil {
			errors = append(errors, fmt.Sprintf("Line %d: invalid "+Magenta+"prefix "+Reset+"format ("+Cyan+"%s"+Reset+")", lineNumber, prefix))
			// continue
		}

		// Skip empty country code
		if countryCode == "" {
			continue
		}

		// Validate country code
		if !countryMap[countryCode] {
			errors = append(errors, fmt.Sprintf("Line %d: invalid "+Magenta+"country code "+Reset+"("+Cyan+"%s"+Reset+")", lineNumber, countryCode))
			continue
		}

		// Skip empty subdivision code
		if subdivisionCode == "" {
			continue
		}

		// Validate subdivision code
		validSubdivision := false
		for _, subdivision := range subdivisions[countryCode] {
			if subdivision.Code == subdivisionCode {
				validSubdivision = true
				break
			}
		}
		if !validSubdivision {
			errors = append(errors, fmt.Sprintf("Line %d: invalid "+Magenta+"subdivision code "+Reset+"("+Cyan+"%s"+Reset+") "+"for country ("+Cyan+"%s"+Reset+")", lineNumber, subdivisionCode, countryCode))
		}
	}

	return errors
}

func main() {
	countryFile := "data/iso3166-1.json"
	subdivisionFile := "data/iso3166-2.json"

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./geofeed-validator <path_to_geofeed.csv>")
		os.Exit(1)
	}

	geofeedFile := os.Args[1]

	var countries []Country
	var subdivisions SubdivisionMap

	if err := readJSON(countryFile, &countries); err != nil {
		fmt.Printf("Failed to load countries: %v\n", err)
		return
	}

	if err := readJSON(subdivisionFile, &subdivisions); err != nil {
		fmt.Printf("Failed to load subdivisions: %v\n", err)
		return
	}

	errors := validateGeofeed(geofeedFile, countries, subdivisions)

	if len(errors) == 0 {
		fmt.Println("Congratulations! Your geofeed file is " + Green + "VALID" + Reset + ".")
	} else {
		fmt.Println("Your geofeed file is " + Red + "INVALID" + Reset + ":")
		for _, err := range errors {
			fmt.Println("- " + err)
		}
	}
}
