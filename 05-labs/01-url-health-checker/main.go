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

// checkURL faz a requisição HTTP em segundo plano (Goroutine) e envia o resultado pelo canal.
func checkURL(url string, wg *sync.WaitGroup, results chan<- Result) {
	// Garante que o WaitGroup será decrementado ao final da execução desta função, mesmo que ocorra um erro.
	defer wg.Done()

	start := time.Now()

	// Configura um cliente com timeout curto de 3s para não travar indefinidamente.
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(url)
	latency := time.Since(start)

	if err != nil {
		// Se der erro (ex: domínio não existe ou timeout), envia o erro pelo canal.
		results <- Result{URL: url, Latency: latency, Err: err}
		return
	}

	// Garante que o corpo da resposta será fechado ao final da função para evitar vazamentos de recursos.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("[ERROR] Failed to close response body for URL %s: %v\n", url, err)
		}
	}(resp.Body)

	// Envia o resultado com sucesso pelo canal.
	results <- Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Latency:    latency,
	}
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

	var wg sync.WaitGroup

	// Criamos um canal com capacidade (buffer) igual ao número de URLs
	// para guardar os resultados de cada Goroutine de forma segura.
	results := make(chan Result, len(urls))

	fmt.Println("[SYSTEM] Checking URL health concurrently...")
	startTime := time.Now()

	// Dispara uma Goroutine para cada URL na lista.
	for _, u := range urls {
		wg.Add(1) // Notifica o WaitGroup: +1 tarefa em andamento
		go checkURL(u, &wg, results)
	}

	// Dispara uma Goroutine auxiliar APENAS para esperar os workers terminarem
	// e fechar o canal de resultados. Se não fechar, o loop de leitura no main travará.
	go func() {
		wg.Wait()
		close(results)
	}()

	// O loop principal consome e exibe os resultados conforme eles vão chegando.
	var successCount, failureCount int

	for res := range results {
		if res.Err != nil {
			failureCount++
			fmt.Printf("[FAIL] %-38s | Latency: %7v | Error: %v\n", res.URL, res.Latency.Round(time.Millisecond), res.Err)
		} else {
			successCount++
			fmt.Printf("[OK]   %-38s | Status: %d | Latency: %7v\n", res.URL, res.StatusCode, res.Latency.Round(time.Millisecond))
		}
	}

	totalDuration := time.Since(startTime)

	fmt.Println("\n================ SUMMARY ================")
	fmt.Printf("Total Elapsed Time : %v\n", totalDuration.Round(time.Millisecond))
	fmt.Printf("URLs Tested        : %d\n", len(urls))
	fmt.Printf("Successful         : %d\n", successCount)
	fmt.Printf("Failed / Errors    : %d\n", failureCount)
	fmt.Println("=========================================")
}
