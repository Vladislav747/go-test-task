package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"net/http"
	"time"
	//"io"
)

func main() {
	fmt.Println("Programm started");
	file, err := os.Open("links.txt")

	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
 
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
 
	file.Close()
 
	for _, eachline := range txtlines {
		//fmt.Println(eachline)
		checkResource(eachline);
	}
}

func checkResource(s string) {
	/* Таймаут 5 секунд */
	client := http.Client{
		Timeout: 5 * time.Second,
	} 
	resp, err := client.Get(s) 
	if err != nil { 
		fmt.Println(err) 
	return
	} 
	defer resp.Body.Close() 
	fmt.Println(resp.StatusCode);
	// if resp.StatusCode == http.StatusOK {
	// 	bodyBytes, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	bodyString := string(bodyBytes)
	// 	//fmt.Println(bodyString)
	// }
}