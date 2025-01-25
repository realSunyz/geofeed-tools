package validate

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/realsunyz/geofeed-tools/plugin/color"
	"github.com/realsunyz/geofeed-tools/plugin/isocode"
)

const (
	countryFile     = "iso3166-1.json"
	subdivisionFile = "iso3166-2.json"
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

var countries []Country
var subdivisions SubdivisionMap

func readJSON(filename string, v interface{}) error {
	data, err := isocode.DataFS.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func readCSV(filePath string) (*bufio.Scanner, *os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(file)

	return scanner, file, nil
}

func genSubdivisionMap(subdivisions SubdivisionMap) map[string]map[string]bool {
	subdivisionMap := make(map[string]map[string]bool)
	for country, subs := range subdivisions {
		subMap := make(map[string]bool)
		for _, sub := range subs {
			subMap[sub.Code] = true
		}
		subdivisionMap[country] = subMap
	}
	return subdivisionMap
}

func checkSyntax(scanner *bufio.Scanner, countryMap map[string]bool, subdivisionMap map[string]map[string]bool) []string {
	lineNumber := 0
	var feedErrors []string

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
			feedErrors = append(feedErrors, fmt.Sprintf("Line %d: invalid "+color.Magenta+"geofeed "+color.Reset+"format", lineNumber))
			continue
		}

		prefix, countryCode, subdivisionCode := record[0], record[1], record[2]

		// Validate IP prefix format
		_, _, prefixErr := net.ParseCIDR(prefix)
		if prefixErr != nil {
			feedErrors = append(feedErrors, fmt.Sprintf("Line %d: invalid "+color.Magenta+"prefix "+color.Reset+"format ("+color.Cyan+"%s"+color.Reset+")", lineNumber, prefix))
			// continue
		}

		// Skip empty country code
		if countryCode == "" {
			continue
		}

		// Validate country code
		if !countryMap[countryCode] {
			feedErrors = append(feedErrors, fmt.Sprintf("Line %d: invalid "+color.Magenta+"country code "+color.Reset+"("+color.Cyan+"%s"+color.Reset+")", lineNumber, countryCode))
			continue
		}

		// Skip empty subdivision code
		if subdivisionCode == "" {
			continue
		}

		// Validate subdivision code
		if subMap, exists := subdivisionMap[countryCode]; exists {
			if _, valid := subMap[subdivisionCode]; !valid {
				feedErrors = append(feedErrors, fmt.Sprintf("Line %d: invalid "+color.Magenta+"subdivision code "+color.Reset+"("+color.Cyan+"%s"+color.Reset+") "+"for country ("+color.Cyan+"%s"+color.Reset+")", lineNumber, subdivisionCode, countryCode))
			}
		} else {
			feedErrors = append(feedErrors, fmt.Sprintf("Line %d: no subdivisions found for country ("+color.Cyan+"%s"+color.Reset+")", lineNumber, countryCode))
		}
	}

	return feedErrors
}

func Execute(filePath string) {

	if err := readJSON(countryFile, &countries); err != nil {
		fmt.Printf("Failed to load countries: %v\n", err)
		return
	}
	if err := readJSON(subdivisionFile, &subdivisions); err != nil {
		fmt.Printf("Failed to load subdivisions: %v\n", err)
		return
	}

	countryMap := make(map[string]bool)
	for _, country := range countries {
		countryMap[country.Code] = true
	}

	subdivisionMap := genSubdivisionMap(subdivisions)

	scanner, file, err := readCSV(filePath)
	if err != nil {
		fmt.Printf("Failed to open geofeed file: %v\n", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Failed to close geofeed file: %v\n", err)
		}
	}()

	feedErrors := checkSyntax(scanner, countryMap, subdivisionMap)

	if len(feedErrors) == 0 {
		fmt.Println("Congratulations! Your geofeed file is " + color.Green + "VALID" + color.Reset + ".")
	} else {
		fmt.Println("Your geofeed file is " + color.Red + "INVALID" + color.Reset + ":")
		for _, errMsg := range feedErrors {
			fmt.Println("- " + errMsg)
		}
	}
}
