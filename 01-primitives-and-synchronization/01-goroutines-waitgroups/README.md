# Goroutines & sync.WaitGroup

Este módulo apresenta os fundamentos da concorrência em Go, explorando a palavra-chave `go`, o **escalonador M:N** (_modelo G-M-P_) e a primitiva de sincronização `sync.WaitGroup`.

O objetivo é compreender o ciclo de vida das goroutines e os mecanismos utilizados pelo _Go Runtime_ para gerenciar a execução concorrente, sem a interferência de regras de negócio ou bibliotecas externas.

- **Propósito:** Demonstrar a mecânica de execução de goroutines e sincronização de tarefas utilizando recursos nativos da linguagem.
- **Aplicação Prática:** Para ver essas primitivas resolvendo um problema real de I/O e latência de rede, consulte o [Lab 01 - Concurrent URL Health Checker](../../05-labs/01-url-health-checker/README.md).

## Concorrência vs. Paralelismo

A diferença entre ambos reside na **arquitetura do código** (software) versus a **execução física** (hardware):

- **Concorrência (Estrutura de Software):** É a capacidade de **gerenciar e intercalar** múltiplas tarefas em um determinado intervalo de tempo. O _runtime_ do Go alterna a execução entre as goroutines (_context-switching_), permitindo que várias tarefas progridam mesmo em um computador com apenas 1 núcleo de CPU.
- **Paralelismo (Execução em Hardware):** É a capacidade de **executar simultaneamente** duas ou mais instruções no exato mesmo milissegundo. Requer obrigatoriamente suporte do hardware (múltiplos núcleos de CPU) e `GOMAXPROCS > 1`.

> [!TIP]
> **Modelo Mental (A Analogia da Cozinha):**
>
> - **Concorrência:** Um único chef gerenciando vários pratos—coloca o pão no forno, corta os legumes enquanto ele assa e tira o pão no tempo certo. Várias tarefas progridem no mesmo período, mas o chef só faz uma ação física por vez.
> - **Paralelismo:** Dois chefs trabalhando juntos ao mesmo tempo, cada um cortando legumes em tábuas separadas no exato mesmo segundo.

## Escalonador M:N do Go (Modelo G-M-P)

O Go gerencia a concorrência por meio de um **Escalonador M:N** interno, mapeando de forma transparente **M** _Goroutines_ sobre **N** _Threads_ do Sistema Operacional.

| Entidade | Nome          | Descrição & Papel                                                                                               |
|:---------|:--------------|:----------------------------------------------------------------------------------------------------------------|
| **G**    | **Goroutine** | _Thread_ leve em nível de usuário (~2 KB iniciais de pilha, crescimento dinâmico).                              |
| **M**    | **Machine**   | _Thread_ real de kernel do SO criada e gerenciada pelo sistema operacional.                                     |
| **P**    | **Processor** | Contexto/recurso lógico de execução necessário para rodar código Go. Por padrão, equivale a `runtime.NumCPU()`. |

> [!IMPORTANT]
> **Ponto-chave (Concorrência vs. Paralelismo no Go):**
>
> A palavra-chave `go` garante apenas **concorrência**. Se essas goroutines serão executadas em **paralelo**, isso depende da quantidade de núcleos de CPU físicos da máquina e da configuração de `GOMAXPROCS`.

## Sincronização de Estado (`sync.WaitGroup`)

O `sync.WaitGroup` funciona como uma barreira de sincronização. Ele mantém um contador interno para evitar que a função `main` finalize o programa antes que todas as goroutines em segundo plano concluam suas tarefas.

| Método      | Operação                     | Regra de Uso                                                                        |
|:------------|:-----------------------------|:------------------------------------------------------------------------------------|
| `wg.Add(n)` | Incrementa o contador        | Deve ser chamado na thread principal **antes** do disparo da goroutine (`go`).      |
| `wg.Done()` | Decrementa o contador em `1` | Chamado dentro da goroutine ao encerrar a execução (recomenda-se o uso de `defer`). |
| `wg.Wait()` | Bloqueia a execução          | Pausa a thread principal até que o contador zere.                                   |
