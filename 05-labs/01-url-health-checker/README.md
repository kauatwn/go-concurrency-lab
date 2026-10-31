# Lab 01 - Concurrent URL Health Checker

Este lab apresenta a aplicação prática da concorrência em Go, utilizando Goroutines, `sync.WaitGroup` e canais (`chan`) para monitorar a disponibilidade HTTP de múltiplos serviços web.

- **Propósito:** Demonstrar como a execução concorrente de operações de I/O reduz o tempo total de monitoramento em comparação com a execução sequencial.
- **Aplicação Prática:** Construir uma ferramenta CLI de monitoramento de URLs, coordenando workers concorrentes e coletando seus resultados por meio de canais.
- **Fundamentos Aplicados:** Para compreender as primitivas de concorrência e sincronização utilizadas neste laboratório, consulte o módulo [Goroutines & sync.WaitGroup](../../01-primitives-and-synchronization/01-goroutines-waitgroups/README.md).

## Descrição do Problema

Verificar o status HTTP de múltiplos serviços web de forma sequencial força a aplicação a executar uma requisição por vez. O tempo total de execução passa a ser a **soma exata da latência de todas as requisições**. Se 7 serviços levarem 500ms cada, uma checagem sequencial levará aproximadamente 3,5 segundos.

Ao aplicar **Goroutines** e **`sync.WaitGroup`**, as verificações ocorrem de forma concorrente em paralelo. O tempo total de execução cai drasticamente, passando a corresponder apenas à latência da **requisição individual mais lenta** (~1,5 segundo).

> [!TIP]
> **Modelo Mental (A Analogia das Chamadas de Telefone):**
>
> Imagine que você precisa ligar para 7 fornecedores para confirmar a entrega de um pedido:
>
> - **Sequencial:** Você liga para o primeiro, espera atender e desligar; depois liga para o segundo... O tempo total é a soma de todas as chamadas.
> - **Concorrente (Goroutines):** Você chama 7 ajudantes na sala e cada um faz 1 ligação ao mesmo tempo. O tempo total será apenas o tempo da ligação que mais demorou.

## Arquitetura & Fluxo de Dados

```text
                     +---> [Worker 1: Google]  ---+
                     |                            |
[Dispatcher (main)] -+---> [Worker 2: GitHub]  ---+---> (results chan) ---> [Collector (main)]
                     |                            |
                     +---> [Worker N: Status]  ---+
                                |
                         (sync.WaitGroup)
```

1. **Dispatcher (`main`)**: Itera sobre a lista de URLs, incrementa o contador com `wg.Add(1)` e dispara uma Goroutine dedicada por endpoint.
2. **Workers (`checkURL`)**: Executam a requisição HTTP GET com timeout de 3 segundos, enviam a struct `Result` para o canal `results` e notificam a conclusão via `defer wg.Done()`.
3. **Goroutine Monitora**: Executa `wg.Wait()` em segundo plano para aguardar a finalização de todas as tarefas e encerra com segurança o canal de resultados (`close(results)`).
4. **Collector (`main`)**: Consome o canal via `range results`, exibindo o log individual de cada requisição e gerando o relatório consolidado ao final.

> [!IMPORTANT]
> **Ponto-chave (Fechamento do Canal com Routine Monitora):**
>
> O loop de leitura `for res := range results` na função `main` continuará bloqueado esperando por novos dados até que o canal seja formalmente fechado com `close(results)`. Por isso, utilizamos uma goroutine de monitoramento em segundo plano com `wg.Wait()`. Se executássemos o `wg.Wait()` na própria `main` antes do loop, o canal encheria e causaria um *deadlock*.

## Conceitos Aplicados

- **Goroutines (`go checkURL`)**: Disparo de tarefas I/O leves executadas em segundo plano pelo Go Scheduler.
- **Barreira de Sincronização (`sync.WaitGroup`)**: Coordenação de estado para impedir que o programa principal encerre antes que todos os trabalhadores finalizem suas requisições.
- **Comunicação por Canais com Buffer (`chan Result`)**: Transferência thread-safe de dados entre goroutines sem acoplamento direto ou condições de corrida em variáveis compartilhadas.
