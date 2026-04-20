# Golang API Rest

Este projeto é uma API REST escrita em Go, estruturada para facilitar o aprendizado de quem vem de linguagens como PHP.

## 🚀 Como Rodar o Projeto

### 1. Usando Docker (Recomendado)

Como o projeto possui um arquivo `docker-compose.yml`, você pode rodar tudo em containers:

*   **Subir o projeto (e buildar):**
    ```bash
    docker-compose up --build
    ```
*   **Subir em background:**
    ```bash
    docker-compose up -d
    ```
*   **Parar o projeto:**
    ```bash
    docker-compose down
    ```

### 2. Rodando Localmente (Sem Docker)

Se você tiver o Go instalado na sua máquina (v1.22+):

*   **Instalar dependências e organizar arquivos:**
    ```bash
    go mod tidy
    ```
*   **Rodar em modo desenvolvimento (tipo `php artisan serve`):**
    ```bash
    go run cmd/api/main.go
    ```
*   **Compilar o projeto (gerar executável):**
    ```bash
    go build -o api cmd/api/main.go
    ```
*   **Executar o arquivo compilado:**
    ```bash
    ./api
    ```

---

## 🧠 Dicas para ex-PHPs (Ponteiros e Tipos)

No Go, temos conceitos que o PHP gerencia automaticamente. Veja o resumo:

| Símbolo | Nome | Função | Analogia PHP |
| :--- | :--- | :--- | :--- |
| `&` | Endereço | Pega o local da memória onde o valor está. | Seria como pegar o ID único de um objeto. |
| `*` | Ponteiro | Indica que a variável guarda um endereço, não o valor real. | Padrão de objetos em classes PHP. |
| `nil` | Nulo | O valor zero para ponteiros, interfaces e slices. | `null` |
| `:=` | Declaração | Cria e atribui uma variável inferindo o tipo. | `$variavel = ...` |

### Exemplo Rápido:
```go
user := User{Name: "Eduardo"} // Criou o valor
ref  := &user                 // Pegou o endereço (&)
```

## 📁 Estrutura de Pastas

*   `cmd/api/`: Ponto de entrada da aplicação (`main.go`).
*   `internal/`: Código privado da aplicação (Controllers, Services, Models, Repositories).
*   `go.mod`: Gerenciador de dependências (o `composer.json`).
*   `go.sum`: Checksum das dependências (o `composer.lock`).
