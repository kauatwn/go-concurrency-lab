# 01 - Primitivas & Sincronização

Este módulo cobre as primitivas fundamentais de concorrência em Go, a mecânica de execução do runtime e as técnicas de sincronização de memória.

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
> **Ponto-chave:** A palavra-chave `go` garante **concorrência**. Se essas goroutines serão executadas em **paralelo**, isso depende da quantidade de núcleos de CPU físicos da máquina e da configuração de `GOMAXPROCS`.
