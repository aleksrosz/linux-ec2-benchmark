package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	instanceName                  string
	primeNumberLimit              string
	cpuSpeed                      string
	throughputEventsPerSecond     string
	throughputTimeElapsed         string
	throughputTotalNumberOfEvents string
	latencyMin                    string
	latencyAvg                    string
	latencyMax                    string
	latency95percentile           string
	threadsFairnessEvents         string
	threadsFairnessExecutionTime  string
)

func readFile(fileName string) (result sysbenchResult) {
	input, err := os.ReadFile("./results/" + fileName)
	if err != nil {
		log.Fatalln(err)
	}
	splitInstanceName := strings.Split(fileName, "_")
	instanceName = strings.TrimSuffix(splitInstanceName[1], ".txt")

	lines := strings.Split(string(input), "\n")
	var element string

	for _, line := range lines {

		if strings.Contains(line, "Prime") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				primeNumberLimit = element
			}
		}
		if strings.Contains(line, "events per second") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				cpuSpeed = element
			}
		}
		if strings.Contains(line, "events/s (eps):") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				throughputEventsPerSecond = element
			}
		}
		if strings.Contains(line, "time elapsed:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				throughputTimeElapsed = element
			}
		}
		if strings.Contains(line, "total number of events:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				throughputTotalNumberOfEvents = element
			}
		}
		if strings.Contains(line, "min:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				latencyMin = element
			}
		}
		if strings.Contains(line, "avg:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				latencyAvg = element
			}
		}
		if strings.Contains(line, "max:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			for _, element = range submatchall {
				latencyMax = element
			}
		}
		if strings.Contains(line, "95th percentile:") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			latency95percentile = submatchall[1]
		}
		if strings.Contains(line, "events (avg/stddev):") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			threadsFairnessEvents = submatchall[0]
		}
		if strings.Contains(line, "execution time (avg/stddev):") {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			submatchall := re.FindAllString(line, -1)
			threadsFairnessExecutionTime = submatchall[0]
		}
	}

	// convert values to sysbenchResults struct
	primeNumberLimitInt, err := strconv.Atoi(primeNumberLimit)
	if err != nil {
		log.Fatalln(err)
	}
	cpuSpeedFloat64, err := strconv.ParseFloat(cpuSpeed, 64)
	if err != nil {
		log.Fatalln(err)
	}
	throughputEventsPerSecondFloat64, err := strconv.ParseFloat(throughputEventsPerSecond, 64)
	if err != nil {
		log.Fatalln(err)
	}
	throughputTimeElapsedFloat64, err := strconv.ParseFloat(throughputTimeElapsed, 64)
	if err != nil {
		log.Fatalln(err)
	}
	throughputTotalNumberOfEventsInt, err := strconv.Atoi(throughputTotalNumberOfEvents)
	if err != nil {
		log.Fatalln(err)
	}
	latencyMinFloat64, err := strconv.ParseFloat(latencyMin, 64)
	if err != nil {
		log.Fatalln(err)
	}
	latencyAvgFloat64, err := strconv.ParseFloat(latencyAvg, 64)
	if err != nil {
		log.Fatalln(err)
	}
	latencyMaxFloat64, err := strconv.ParseFloat(latencyMax, 64)
	if err != nil {
		log.Fatalln(err)
	}
	latency95percentileFloat64, err := strconv.ParseFloat(latency95percentile, 64)
	if err != nil {
		log.Fatalln(err)
	}
	threadsFairnessEventsFloat64, err := strconv.ParseFloat(threadsFairnessEvents, 64)
	if err != nil {
		log.Fatalln(err)
	}
	threadsFairnessExecutionTimeFloat64, err := strconv.ParseFloat(threadsFairnessExecutionTime, 64)
	if err != nil {
		log.Fatalln(err)
	}
	result = sysbenchResult{
		instanceName:                  instanceName,
		primeNumberLimit:              primeNumberLimitInt,
		cpuSpeed:                      cpuSpeedFloat64,
		throughputEventsPerSecond:     throughputEventsPerSecondFloat64,
		throughputTimeElapsed:         throughputTimeElapsedFloat64,
		throughputTotalNumberOfEvents: throughputTotalNumberOfEventsInt,
		latencyMin:                    latencyMinFloat64,
		latencyAvg:                    latencyAvgFloat64,
		latencyMax:                    latencyMaxFloat64,
		latency95percentile:           latency95percentileFloat64,
		threadsFairnessEvents:         threadsFairnessEventsFloat64,
		threadsFairnessExecutionTime:  threadsFairnessExecutionTimeFloat64,
	}
	return result
}

func appendToCSVFile() {
	input, err := os.ReadFile("results_temporary.csv")
	if err != nil {
		log.Fatalln(err)
	}

	result, err := os.ReadFile("results.csv")
	if err != nil {
		result, err = os.ReadFile("template.csv")
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
	err = os.WriteFile("results.csv", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
