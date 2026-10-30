package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker simula uma tarefa demorada executada em segundo plano (como I/O de rede ou processamento).
// Recebe um ID e um ponteiro para o sync.WaitGroup, que é usado para sincronizar a conclusão da goroutine
// (gerenciar a barreira de sincronização).
func Worker(id int, wg *sync.WaitGroup) {

	// O 'defer' garante que o Done() será executado ao sair da função Worker, mesmo que ocorra um erro ou uma saída
	// antecipada.
	defer wg.Done() // Importante para evitar deadlocks no WaitGroup.

	fmt.Printf("[Worker %d] Starting work...\n", id)

	// Simula algum trabalho demorado (como I/O de rede ou processamento).
	time.Sleep(1 * time.Second) // Simula um atraso de 1 segundo.

	fmt.Printf("[Worker %d] Work finished!\n", id)
}

func main() {
	// sync.WaitGroup é o contador de sincronização do Go.
	// Conceitualmente, é semelhante a esperar um Task.WhenAll no C#.
	var wg sync.WaitGroup
	workers := 5

	fmt.Println("[MAIN] Spawning goroutines...")

	for i := 1; i <= workers; i++ {
		// Incrementa o contador antes de disparar a goroutine
		wg.Add(1)

		// A palavra-chave 'go' inicia a função de forma assíncrona em uma Goroutine leve, disparando a execução do
		// Worker de forma concorrente.
		go Worker(i, &wg)

		// Obs.: Passar o WaitGroup (&wg) como ponteiro. Se passar por valor (sem o &), a goroutine receberá uma cópia
		// e o contador original nunca chegará a zero, resultando em um deadlock.
	}

	fmt.Println("[MAIN] Waiting for workers to complete...")

	// Bloqueia a execução até que o contador do WaitGroup chegue a zero, aguardando que todas as goroutines concluam
	// suas execuções e chamem Done().
	wg.Wait()

	fmt.Println("[MAIN] All workers finished. Exiting process safely.")
}
