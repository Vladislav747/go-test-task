package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"net/http"
	"time"
	//"io"
	"sync"
)

func main() {
	fmt.Println("Programm started");
	start := time.Now()
	ch := make(chan string)

	file, err := os.Open("links.txt")
	//var buf string
	var wg sync.WaitGroup
	

	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}

	createNewFile()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
 
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
 
	file.Close()
 
	for _, eachline := range txtlines {
		wg.Add(1)
		go checkResource(eachline, ch, &wg);
		valFromCh := string(<-ch)
		fmt.Println(valFromCh)
	
		
	}
	wg.Wait()

	// for range txtlines {
		
		
	// 	//buf += valFromCh + "\n"
	// }

	//Получение из канала ch
//	writeToFile("results.txt", buf)

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func checkResource(url string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	/* Таймаут 5 секунд */
	client := http.Client{
		Timeout: 5 * time.Second,
	} 
	resp, err := client.Get(url) 
	if err != nil { 
		fmt.Printf("Error host - %s", url); 
		fmt.Println(err)
		ch <- fmt.Sprint(err) // отправка в канал ch
		return
	} 
	defer resp.Body.Close()  // исключение утечки ресурсов
	//fmt.Printf("Статус хоста %s - %d", url, resp.StatusCode);


	if resp.StatusCode == http.StatusOK {
		// bodyBytes, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// bodyString := string(bodyBytes)
		//fmt.Println(bodyString)

		ch <- fmt.Sprintf("Статус хоста %s - %d", url, resp.StatusCode)
	}
}

func createNewFile(){
	_, err := os.Create("results.txt")

	if err != nil{
        fmt.Println("Unable to create file:", err) 
        os.Exit(1)
		os.Remove("1.txt")
		res := deleteFile("results.txt")
		if(res){
			createNewFile();
		}
    }
}

func deleteFile(nameFile string) bool {
	err := os.Remove(nameFile)

	if err != nil {
		log.Fatalf("Failed deleting file: %s", err)
		return false
	}
	return true
}

func writeToFile(nameFile string, data string) {
	file, err := os.Open(nameFile)

	if err != nil{
        fmt.Println(err) 
        os.Exit(1) 
    }

	defer file.Close() 
    file.Write([]byte(data))
     
    fmt.Println("Done.")
}