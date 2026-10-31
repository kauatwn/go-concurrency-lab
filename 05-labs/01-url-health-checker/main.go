package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Result representa a resposta do teste de saúde de uma URL.
type Result struct {
	URL        string
	StatusCode int
	Latency    time.Duration
	Err        error
}

// Checker encapsula a instância do cliente HTTP para ser reutilizada entre Goroutines.
type Checker struct {
	client *http.Client
}

// NewChecker cria um novo Checker com um cliente HTTP configurado com o timeout especificado.
func NewChecker(timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Check faz a requisição HTTP GET para uma única URL e retorna o resultado da medição.
func (c *Checker) Check(url string) Result {
	start := time.Now()

	resp, err := c.client.Get(url)
	latency := time.Since(start)

	if err != nil {
		// Se ocorrer um erro (ex: falha de DNS ou timeout), retorna o erro capturado.
		return Result{URL: url, Latency: latency, Err: err}
	}

	// Garante que o corpo da resposta será fechado ao final da função para evitar vazamentos de recursos.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("[ERROR] Failed to close response body for URL %s: %v\n", url, err)
		}
	}(resp.Body)

	// Consome e descarta o corpo da resposta HTTP para garantir a reutilização da conexão TCP (Keep-Alive).
	_, _ = io.Copy(io.Discard, resp.Body)

	return Result{URL: url, StatusCode: resp.StatusCode, Latency: latency}
}

// RunChecks orquestra a execução concorrente dos testes de URLs, retornando um canal somente leitura (<-chan Result).
// Isso permite que o chamador consuma e exiba os resultados em tempo real (streaming) à medida que são concluídos.
func RunChecks(checker *Checker, urls []string) <-chan Result {
	// Canal com capacidade (buffer) dimensionada ao número de URLs
	// para receber os resultados com segurança.
	results := make(chan Result, len(urls))
	var wg sync.WaitGroup

	// Dispara uma Goroutine para cada URL na lista.
	for _, url := range urls {
		wg.Add(1) // Notifica o WaitGroup: +1 tarefa em andamento

		go func(targetURL string) {
			defer wg.Done() // Garante a baixa da tarefa ao finalizar a função
			results <- checker.Check(targetURL)
		}(url)
	}

	// Dispara uma Goroutine auxiliar APENAS para esperar os workers terminarem
	// e fechar o canal de resultados. Se não fechar, o loop de leitura no main travará.
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// PrintResult exibe o log formatado no terminal para um resultado individual.
func PrintResult(res Result) {
	if res.Err != nil {
		fmt.Printf("[FAIL] %-38s | Latency: %7v | Error: %v\n", res.URL, res.Latency.Round(time.Millisecond), res.Err)
		return
	}

	fmt.Printf("[OK]   %-38s | Status: %d | Latency: %7v\n", res.URL, res.StatusCode, res.Latency.Round(time.Millisecond))
}

// PrintSummary exibe o relatório consolidado final com as estatísticas de execução.
func PrintSummary(totalDuration time.Duration, totalTested, successCount, failureCount int) {
	fmt.Println("\n================ SUMMARY ================")
	fmt.Printf("Total Elapsed Time : %v\n", totalDuration.Round(time.Millisecond))
	fmt.Printf("URLs Tested        : %d\n", totalTested)
	fmt.Printf("Successful         : %d\n", successCount)
	fmt.Printf("Failed / Errors    : %d\n", failureCount)
	fmt.Println("=========================================")
}

func main() {
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
		"https://httpbin.org/status/200",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/status/500",
		"https://domain-does-not-exist-xyz.com",
	}

	// Configura um cliente com timeout curto de 3s para não travar indefinidamente.
	checker := NewChecker(3 * time.Second)

	fmt.Println("[SYSTEM] Checking URL health concurrently...")
	startTime := time.Now()

	// Dispara a checagem concorrente e obtém o canal de streaming
	resultsChan := RunChecks(checker, urls)

	var successCount, failureCount int

	// O loop principal consome e exibe os resultados em tempo real conforme chegam no canal.
	for res := range resultsChan {
		if res.Err != nil {
			failureCount++
		} else {
			successCount++
		}
		PrintResult(res)
	}

	totalDuration := time.Since(startTime)
	PrintSummary(totalDuration, len(urls), successCount, failureCount)
}
