package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-test-task/title"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	FileName        = "links.txt"
	ResultsFileName = "results.txt"
)

var client = http.DefaultClient

func main() {
	log.Println("Program started")
	//Положить горутину в отд ф-цию чтобы сборщику мусора было легче
	func() {
		startTime := checkFile()
		fmt.Printf("%.2fs elapsed\n", time.Since(startTime).Seconds())
	}()

	// Все то касается изменение файла
	doneChan := make(chan bool)

	go waitForFile(doneChan)

	//Ждем изменений в файле и перезапускаем проверку файла
	if <-doneChan {
		fmt.Println("check file")
		startTime := checkFile()
		fmt.Printf("%.2fs elapsed\n", time.Since(startTime).Seconds())

	}
}

func checkFile() time.Time {
	// Получение структуры файла с ссылками
	urlLinksFile, err := os.Open(FileName)
	if err != nil {
		log.Fatal("error on read file " + urlLinksFile.Name())
	}

	resultFile := createNewFile()

	var wg sync.WaitGroup

	startTime := time.Now()

	//Иницализирую канал
	chTxtlines := make(chan string, 5)
	resultUrlResp := make(chan string)

	// Создание буфера, который читает строки файла с ссылками
	scanner := bufio.NewScanner(urlLinksFile)
	scanner.Split(bufio.ScanLines)

	go writeToFile(resultFile, resultUrlResp)

	for scanner.Scan() {
		wg.Add(1)
		go checkResource(resultUrlResp, chTxtlines, &wg)
		chTxtlines <- scanner.Text()

	}

	//Закрыть файл
	urlLinksFile.Close()

	wg.Wait()

	//Закрыть канал
	close(resultUrlResp)
	return startTime
}

func waitForFile(doneChan chan bool) {

	fmt.Println("Wait for changes links.txt")

	err := watchFile("links.txt")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("File has been changed")
	doneChan <- true
}

/* Проверка ресурса на доступность */
func checkResource(resultCh chan string, chWithUrlLine <-chan string, wg *sync.WaitGroup) error {
	defer wg.Done()

	//Функции дается 5 секунд на проверку иначе срабатывает cancel context
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
			log.Printf("[checkResource]Error host - %s : %s\n", url, err.Error())

			resultCh <- fmt.Sprintf("%s | %s", url, "500") // отправка в канал ch

			return err
		}

		defer resp.Body.Close() // исключение утечки памяти

		log.Printf("[checkResource] [%d] Url - %s \n", resp.StatusCode, url)

		if resp.StatusCode == http.StatusOK {

			//Если ответ корректный то проходимся по структуре сайта
			if title, ok := title.GetHtmlTitle(resp.Body); ok {
				resultCh <- fmt.Sprintf("%s | %d | %s", url, resp.StatusCode, title)
			} else {
				resultCh <- fmt.Sprintf("%s | %s | %s", url, resp.StatusCode, "Can't get title")
			}

		} else {
			resultCh <- fmt.Sprintf("%s | %d", url, resp.StatusCode)
		}

	} else {
		log.Printf("[checkResource] канал закрыт")

	}
	return nil
}

func createNewFile() *os.File {
	resultFile, err := os.Create("results.txt")

	if err != nil {
		fmt.Println("Unable to create file:", err)
		res := deleteFile("results.txt")
		if res {
			createNewFile()
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

// Запись в файл
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

/*
*
Проверка файла пока в файле не появятся изменения
*/
func watchFile(filePath string) error {
	//Изначальное время модификации файла
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	//Бесконечная проверка с частотой раз в 1 сек - time.Sleep(1*time.Second)
	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		//Текущее время модификации файла
		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
