package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"net/http"
	"time"
	"context"
	"sync"
)

//Интерфейс проверки источника
type LookUpResolver struct {
	txtlines          [] string
	//канал проверки строчек файла
	chTxtlinesQueue      chan string
	resolveResultsMap       map[string]string
	mutexDevice      sync.RWMutex
	mutexResolved     sync.RWMutex
}


const (
	FileName = "links.txt"
	ResultsFileName = "results.txt"
)

var client = http.DefaultClient

type ResultMap map[string]string

func main() {
	log.Println("Programm started")

	// Получение структуры файла с ссылками
	urlLinksFile, err := os.Open(FileName)
	if err != nil {
		log.Fatal("error on read file " + urlLinksFile.Name())
	}


	// resolver := newResolver();

	// //Запуск метода LookUpresolver
	// resolver.Run()

	resultFile := createNewFile()

	var wg sync.WaitGroup

	startTime := time.Now()


	// var txtlines []string
	//Иницализирую канал
	chTxtlines := make(chan string, 5)
	resultUrlResp := make(chan string)

	fmt.Printf("type of `c` is %T\n", chTxtlines)
	fmt.Printf("value of `c` is %v\n", chTxtlines)
	
	// Создание буфера, который читает строки файла с ссылками
	scanner := bufio.NewScanner(urlLinksFile)
	scanner.Split(bufio.ScanLines)

	go writeToFile(resultFile, resultUrlResp)

	for scanner.Scan() {
		wg.Add(1)
		go checkResource(resultUrlResp, chTxtlines, &wg);
		chTxtlines <- scanner.Text()
 
	}
 
	//Закрыть файл
	urlLinksFile.Close()
 
	wg.Wait()

	//Закрыть канал
	close(resultUrlResp)
	
	fmt.Printf("%.2fs elapsed\n", time.Since(startTime).Seconds())
	return
}

//Проверка ресурса на доступность
func checkResource(resultCh chan string, chWithUrlLine <-chan string, wg *sync.WaitGroup)error {
	defer wg.Done()

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	
	url, ok := <-chWithUrlLine
	if ok {
		// https://gist.github.com/2minchul/6d344a0f1f85ead1530803df2e4f9894 - объяснение запроса с контекстом
		req, err := http.NewRequestWithContext(context, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("[checkResource]Error host - %s : %s\n", url, err.Error())
			
			resultCh <- fmt.Sprintf("%s | %s", url, err.Error()) // отправка в канал ch

			return err
		}

		defer resp.Body.Close() // исключение утечки памяти

		log.Printf("[checkResource] [%d] Url - %s \n", resp.StatusCode, url)

		resultCh <- fmt.Sprintf("%s | %d", url, resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			// bodyBytes, err := io.ReadAll(resp.Body)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// bodyString := string(bodyBytes)
			//fmt.Println(bodyString)

			// chW[url] = fmt.Sprintf("Статус хоста %s - %d", url, resp.StatusCode)
		}
	} else {
		log.Printf("[checkResource] канал закрыт")
		
	}
	return nil
}

func createNewFile()*os.File{
	resultFile, err := os.Create("results.txt")

	if err != nil{
        fmt.Println("Unable to create file:", err) 
		res := deleteFile("results.txt")
		if(res){
			createNewFile();
		}
    }

	return resultFile
}

func deleteFile(nameFile string) bool {
	err := os.Remove(nameFile)

	if err != nil {
		log.Fatalf("Failed deleting file: %s", err)
		return false
	}

	log.Printf("Deleted file - %s", nameFile)
	return true
}

func writeToFile(fileForWrite *os.File, resultUrlResp chan string) error {
	defer fileForWrite.Close() 

	for resultUrl := range resultUrlResp {
		if _, err := fileForWrite.Write([]byte(resultUrl + "\n")); err != nil {
			log.Fatalf("[writeToFile]Failed writing to file: %s", err)
		}
	}
	
	log.Printf("Data is written to file %s. \n", fileForWrite.Name())
	return nil
}


func newResolver() *LookUpResolver {
	return &LookUpResolver{
		chTxtlinesQueue:      make(chan string, 10000),
		resolveResultsMap:       make(map[string]string),
	}
}
