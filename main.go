package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Result struct {
	StatusCode int
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requisições")
	concurrency := flag.Int("concurrency", 0, "Número de chamadas simultâneas")

	flag.Parse()

	if err := validate(*url, *requests, *concurrency); err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}

	fmt.Println("Iniciando teste de carga...")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Requests: %d\n", *requests)
	fmt.Printf("Concorrência: %d\n", *concurrency)
	fmt.Println()

	start := time.Now()

	results := make(chan Result, *requests)
	jobs := make(chan struct{}, *requests)

	var wg sync.WaitGroup
	var completed int64

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range jobs {
				statusCode := executeRequest(client, *url)

				results <- Result{
					StatusCode: statusCode,
				}

				atomic.AddInt64(&completed, 1)
			}
		}()
	}

	for i := 0; i < *requests; i++ {
		jobs <- struct{}{}
	}

	close(jobs)
	wg.Wait()
	close(results)

	elapsed := time.Since(start)

	statusCount := make(map[int]int)

	for result := range results {
		statusCount[result.StatusCode]++
	}

	printReport(*requests, completed, elapsed, statusCount)
}

func validate(url string, requests, concurrency int) error {
	if url == "" {
		return fmt.Errorf("o parâmetro --url é obrigatório")
	}

	if requests <= 0 {
		return fmt.Errorf("o parâmetro --requests deve ser maior que zero")
	}

	if concurrency <= 0 {
		return fmt.Errorf("o parâmetro --concurrency deve ser maior que zero")
	}

	return nil
}

func executeRequest(client *http.Client, url string) int {
	response, err := client.Get(url)

	if err != nil {
		return 0
	}

	defer response.Body.Close()

	return response.StatusCode
}

func printReport(totalRequests int, completed int64, elapsed time.Duration, statusCount map[int]int) {
	fmt.Println("========================================")
	fmt.Println("           RELATÓRIO DO TESTE")
	fmt.Println("========================================")
	fmt.Printf("Tempo total: %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("Requests realizados: %d\n", completed)
	fmt.Printf("Requests solicitados: %d\n", totalRequests)
	fmt.Printf("Status HTTP 200: %d\n", statusCount[http.StatusOK])
	fmt.Println()
	fmt.Println("Distribuição dos status HTTP:")

	statusCodes := make([]int, 0, len(statusCount))
	for statusCode := range statusCount {
		statusCodes = append(statusCodes, statusCode)
	}

	sort.Ints(statusCodes)

	for _, statusCode := range statusCodes {
		if statusCode == 0 {
			fmt.Printf("Falha de conexão: %d\n", statusCount[statusCode])
			continue
		}

		fmt.Printf("HTTP %d: %d\n", statusCode, statusCount[statusCode])
	}

	fmt.Println("========================================")
}
