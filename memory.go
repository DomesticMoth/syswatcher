package main

import (
	"os"
	//"fmt"
	"bufio"
	"strings"
	"strconv"
)

type Memory struct {
    MemTotal     uint64
    MemFree      uint64
    MemAvailable uint64
}

func ReadMemoryStats() Memory {
    file, err := os.Open("/proc/meminfo")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    bufio.NewScanner(file)
    scanner := bufio.NewScanner(file)
    res := Memory{}
    for scanner.Scan() {
        key, value := parseLine(scanner.Text())
        switch key {
        case "MemTotal":
            res.MemTotal = value
        case "MemFree":
            res.MemFree = value
        case "MemAvailable":
            res.MemAvailable = value
        }
    }
    return res
}

func parseLine(raw string) (key string, value uint64) {
    //fmt.Println(raw)
    text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
    keyValue := strings.Split(text, ":")
    return keyValue[0], uint64(toInt(keyValue[1]))
}

func toInt(raw string) int {
    if raw == "" {
        return 0
    }
    res, err := strconv.Atoi(raw)
    if err != nil {
        panic(err)
    }
    return res
}
