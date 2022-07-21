package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func ifFileExist(fileName string) bool {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("File does not exist.")
			return false
		}
	}
	log.Println("File does exist. File information:")
	log.Println(fileInfo)
	return true
}

func readFile(fileName string, instanceName string) {
	ifFileExist(fileName)

	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	var element string
	var lines2 []string

	for i, line := range lines {

		if strings.Contains(line, "Prime") {
			lines2 = append(lines2, instanceName)
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				lines[i] = element
				lines2 = append(lines2, lines[i])
			}
		}
		if strings.Contains(line, "events per second") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				lines[i] = element
				lines2 = append(lines2, lines[i])
			}
		}
		if strings.Contains(line, "events/s (eps):") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				lines[i] = element
				lines2 = append(lines2, lines[i])
			}
		}
		if strings.Contains(line, "time elapsed:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				lines[i] = element
				lines2 = append(lines2, lines[i])
			}
		}
		if strings.Contains(line, "total number of events:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				lines[i] = element
				lines2 = append(lines2, lines[i])
			}
		}
		if strings.Contains(line, "min:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				lines[i] = element
				lines2 = append(lines2, lines[i])
			}
		}
		if strings.Contains(line, "avg:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				lines[i] = element
				lines2 = append(lines2, lines[i])
			}
		}
		if strings.Contains(line, "max:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				lines[i] = element
				lines2 = append(lines2, lines[i])
			}
		}
		if strings.Contains(line, "95th percentile:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			lines[i] = submatchall[1]
			lines2 = append(lines2, lines[i])
		}
		if strings.Contains(line, "events (avg/stddev):") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			lines[i] = submatchall[0]
			lines2 = append(lines2, lines[i])
		}
		if strings.Contains(line, "execution time (avg/stddev):") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			lines[i] = submatchall[0]
			lines2 = append(lines2, lines[i])
		}
	}
	output := strings.Join(lines2, "\n")
	lines2 = nil
	err = ioutil.WriteFile("results_temporary.csv", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}

func appendToCSVFile() {
	input, err := ioutil.ReadFile("results_temporary.csv")
	if err != nil {
		log.Fatalln(err)
	}

	result, err := ioutil.ReadFile("results.csv")
	if err != nil {
		result, err = ioutil.ReadFile("template.csv")
	}

	lines := strings.Split(string(input), "\n")
	lines2 := strings.Split(string(result), "\n")

	for i := range lines {
		lines2[i] = lines2[i] + ","
	}
	for i := range lines {
		lines2[i] = lines2[i] + lines[i]
	}

	output := strings.Join(lines2, "\n")
	err = ioutil.WriteFile("results.csv", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}
