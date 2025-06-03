package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type Result struct {
	statusCode int
	duration   time.Duration
}

func worker(id int, wg *sync.WaitGroup, url string, jobs <-chan int, results chan<- Result) {
	defer wg.Done()
	for range jobs {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)

		if err != nil {
			results <- Result{statusCode: 0, duration: duration}
			continue
		}

		results <- Result{statusCode: resp.StatusCode, duration: duration}
		resp.Body.Close()
	}
}

func runStressTest(url string, requests int, concurrency int) {
	fmt.Printf("Iniciando teste de carga em %s com %d requests e %d concorrentes\n", url, requests, concurrency)

	startTime := time.Now()

	jobs := make(chan int, requests)
	results := make(chan Result, requests)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(i, &wg, url, jobs, results)
	}

	for i := 0; i < requests; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	close(results)

	// Relatório
	total := 0
	success := 0
	statusCount := make(map[int]int)
	totalDuration := time.Since(startTime)

	for res := range results {
		total++
		statusCount[res.statusCode]++
		if res.statusCode == 200 {
			success++
		}
	}

	fmt.Println("===== RELATÓRIO =====")
	fmt.Printf("Tempo total: %v\n", totalDuration)
	fmt.Printf("Total de requests: %d\n", total)
	fmt.Printf("Status 200 OK: %d\n", success)
	fmt.Println("Outros códigos de status:")
	for code, count := range statusCount {
		if code != 200 {
			fmt.Printf("  %d: %d\n", code, count)
		}
	}
}

func main() {
	var url string
	var requests int
	var concurrency int

	var rootCmd = &cobra.Command{
		Use:   "stress-test",
		Short: "Ferramenta de teste de carga para APIs",
		Long:  `Ferramenta para realizar testes de carga em APIs, permitindo configurar o número de requests e concorrência.`,
		Example: `stress-test --url=https://httpbin.org/get --requests=100 --concurrency=10
stress-test -u https://example.com -r 500 -c 20`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if url == "" {
				return fmt.Errorf("URL é obrigatória. Use --url ou -u para especificar a URL")
			}
			if requests <= 0 {
				return fmt.Errorf("número de requests deve ser maior que 0")
			}
			if concurrency <= 0 {
				return fmt.Errorf("número de concorrentes deve ser maior que 0")
			}

			runStressTest(url, requests, concurrency)
			return nil
		},
	}

	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 1, "Número total de requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Número de chamadas simultâneas")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Erro: %v\n", err)
	}
}
