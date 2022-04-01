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
)

func main() {
	fmt.Println("Programm started")
	start := time.Now()
	ch := make(chan string)

	// resolver := newResolver();

	// //Запуск метода LookUpresolver
	// resolver.Run()



	file, err := os.Open(FileName)
	//var buf string
	var wg sync.WaitGroup
	
	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}

	createNewFile()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	//Иницализирую канал
	chTxtlines := make(chan string, 5)

	fmt.Printf("type of `c` is %T\n", chTxtlines)
	fmt.Printf("value of `c` is %v\n", chTxtlines)

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		// val, _ := <-chTxtlines
		// fmt.Printf("value of `c` is %v\n", val)
		go checkResource(ch, chTxtlines, &wg);
		// valFromCh := string(<-ch)
		// fmt.Println(valFromCh)
		// writeToFile("results.txt", valFromCh)
	}
 
	for scanner.Scan() {
		//txtlines = append(txtlines, scanner.Text())
		chTxtlines <- scanner.Text()
		
 
	}
 
	file.Close()
 
		for _, eachline := range txtlines {
			ch <- eachline
		}	
		
	wg.Wait()

	// for range txtlines {
		
	// 	//buf += valFromCh + "\n"
	// }

	//Получение из канала ch
//	writeToFile("results.txt", buf)

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

//Проверка ресурса на доступность
func checkResource(chW chan<- string, chR <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	/* Таймаут 5 секунд */
	client := http.Client{
		Timeout: 5 * time.Second,
	} 
	url, ok := <-chR
	if ok {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Error host - %s", url)
			fmt.Println(err)
			chW <- fmt.Sprint(err) // отправка в канал ch
			return
		}
		defer resp.Body.Close() // исключение утечки ресурсов
		//fmt.Printf("Статус хоста %s - %d", url, resp.StatusCode);

		if resp.StatusCode == http.StatusOK {
			// bodyBytes, err := io.ReadAll(resp.Body)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// bodyString := string(bodyBytes)
			//fmt.Println(bodyString)

			chW <- fmt.Sprintf("Статус хоста %s - %d", url, resp.StatusCode)
		}
	} else {
		fmt.Println("канал закрыт")

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
	defer file.Close() 

	if err != nil{
        fmt.Println(err) 
        os.Exit(1) 
    }

	
	_, errNew := file.Write([]byte(data + "\n"))
	if errNew != nil {
		log.Fatalf("Failed writing to file: %s", errNew)
	}

	fmt.Printf("Data is written to file %s. \n", nameFile)
}


func newResolver() *LookUpResolver {
	return &LookUpResolver{
		chTxtlinesQueue:      make(chan string, 10000),
		resolveResultsMap:       make(map[string]string),
	}
}
