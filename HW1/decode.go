package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func populateDic() map[string]string {

	file, err := os.Open("dictionary.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	dict := make(map[string]string)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		value := scanner.Text()
		hash := sha256.New()
		hash.Write([]byte(value))
		key := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
		dict[key] = value
	}

	file.Close()

	return dict
}

func populateNames() map[string]string {

	file, err := os.Open("names.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	names := make(map[string]string)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		name := scanner.Text()
		namespace := "d9b2d63d-a233-4123-847a-76838bf2413a"
		tempUUID := uuid.FromStringOrNil(namespace)
		key := uuid.NewV5(tempUUID, name).String()
		names[key] = name
	}

	file.Close()

	return names
}

func timeStampHash(timestamp string) string {
	loc, _ := time.LoadLocation("America/Regina")

	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
	}
	tm := time.Unix(i, 0).In(loc)
	arguments := strings.Split(tm.String(), " ")

	return arguments[0] + "T" + arguments[1] + arguments[2]
}

func main() {
	dictMap := populateDic()
	namesMap := populateNames()

	file, err := os.Open("database_dump.csv")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		arguments := strings.Split(line, ",")

		if arguments[0] == "username" {
			fmt.Println("username,password,last_access")
		} else {
			fmt.Println(namesMap[arguments[0]] + "," + dictMap[arguments[1]] + "," + timeStampHash(arguments[2]))
		}

	}

	file.Close()

}
