package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// worker legge numeri da jobs, fa un lavoro finto (sleep variabile) e manda il risultato su results.
// Si ferma se il contesto viene cancellato o se il canale jobs viene chiuso.
func worker(ctx context.Context, id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			// Interrompe subito e rientra
			return

		case n, ok := <-jobs:
			if !ok {
				// Nessun altro job
				return
			}

			// Simula lavoro (100–400ms)
			workTime := time.Duration(100+rand.Intn(300)) * time.Millisecond
			select {
			case <-ctx.Done():
				return
			case <-time.After(workTime):
				// lavoro "completato"
			}

			// Prova ad inviare il risultato, ma resta reattivo alla cancellazione
			select {
			case results <- fmt.Sprintf("worker %d: f(%d) = %d (in %v)", id, n, n*n, workTime):
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	// Cancella tutto se non finiamo entro 800ms.
	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()

	jobs := make(chan int)
	results := make(chan string)

	const numWorkers = 3
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Avvia i worker
	for i := 1; i <= numWorkers; i++ {
		go worker(ctx, i, jobs, results, &wg)
	}

	// Producer: invia 10 job e chiude il canale
	go func() {
		for n := 1; n <= 10; n++ {
			select {
			case jobs <- n:
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	// Quando tutti i worker terminano, chiudiamo results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collector: legge risultati finché disponibili o finché scatta il timeout
	for {
		select {
		case msg, ok := <-results:
			if !ok {
				fmt.Println("✅ Tutto terminato (results chiuso).")
				return
			}
			fmt.Println(msg)

		case <-ctx.Done():
			fmt.Println("⏰ Timeout raggiunto, cancello il lavoro in corso...")
			// I worker usciranno perché leggono ctx.Done(); aspettiamo la chiusura di results.
			// (qui potresti anche fare un return se non ti interessa attendere)
		}
	}
}
