# Golang API Rest

API REST em Go com autenticação JWT, MySQL e envio de email — prova de conceito com foco em código limpo e idiomático.

## 🚀 Como Rodar

### 1. Pré-requisitos

- **Docker** e **Docker Compose** (recomendado) ou **Go 1.26+**

### 2. Configuração

```bash
cp .env.example .env
```

Edite o `JWT_SECRET` no `.env` para uma chave secreta. Os defaults funcionam com o Docker Compose.

### 3. Subir com Docker

**Modo desenvolvimento (hot reload)** — recomendado para codar:

```bash
docker compose --profile dev up -d

# Acompanhar os logs
docker logs golang_api_app_dev -f
```

Qualquer alteração no código é recarregada automaticamente. Não precisa rebuildar.

**Modo produção/build estático** — simula o que vai pro ECS:

```bash
docker compose --profile prod up -d
```

> Use `docker compose --profile dev down` ou `--profile prod down` para parar cada modo.

### 4. Rodar localmente (sem Docker)

```bash
# Pré-requisito: ter MySQL rodando na máquina

# Instalar dependências
go mod tidy

# Rodar
go run cmd/api/main.go
```

---

## 📖 Documentação Swagger

Com a API rodando, acesse:

➡️ **http://localhost:8080/swagger/index.html**

Na interface Swagger você pode:

- Ver todos os endpoints disponíveis
- Testar chamadas diretamente pela UI
- Clicar em **Authorize** e colar o token JWT para testar rotas protegidas

### Regenerar a documentação

```bash
~/go/bin/swag init -g cmd/api/main.go
```

> Se `swag` não for encontrado, instale com: `go install github.com/swaggo/swag/cmd/swag@latest` e adicione `export PATH=$PATH:~/go/bin` ao seu `~/.zshrc`.

---

## 🧪 Endpoints

### Autenticação (públicos)

| Método | Rota | Descrição |
|--------|------|-----------|
| `POST` | `/api/v1/auth/register` | Criar conta (retorna dados + JWT) |
| `POST` | `/api/v1/auth/login` | Login (retorna access + refresh token) |
| `POST` | `/api/v1/auth/refresh` | Renovar tokens com refresh token |
| `GET` | `/api/v1/auth/verify-email?token=...` | Verificar email |

### Usuários (protegidos — exigem `Authorization: Bearer <token>`)

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/api/v1/users` | Listar todos |
| `GET` | `/api/v1/users/{id}` | Buscar por ID |
| `PUT` | `/api/v1/users/{id}` | Atualizar dados |
| `DELETE` | `/api/v1/users/{id}` | Deletar (requer role **admin**) |

---

## 🔐 Fluxo de Autenticação

```
1. POST /register → cria usuário, gera token de verificação, envia email
2. GET /verify-email?token=xxx → confirma o email
3. POST /login → retorna access_token (15min) + refresh_token (24h)
4. GET /users -H "Authorization: Bearer <access_token>" → dados protegidos
5. POST /refresh → renova os tokens quando o access_token expirar
```

### Headers para rotas protegidas

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

---

## 📁 Estrutura de Pastas

```
cmd/api/main.go              # Entry point + DI
internal/
  config/config.go           # Config via .env
  database/mysql.go          # Conexão MySQL + migrações
  models/                    # Structs de dados (User, DTOs)
  repositories/              # Camada de dados (MySQLRepository)
  services/                  # Regras de negócio (auth, user, email)
  controllers/               # Handlers HTTP
  middleware/                # Logging, JWT auth, roles
  utils/                     # Hash, JWT, validação
docs/                        # Swagger docs (gerado automaticamente)
```

### Padrão arquitetural

O projeto segue **separação por camadas** (tipo MVC hexagonal):

```
Controller → Service → Repository → MySQL
                ↓
          EmailService (SMTP)
```

Cada camada depende apenas da camada abaixo via **interfaces**, facilitando testes e manutenção.

---

## 🛠️ Comandos Úteis

```bash
# Compilar
go build ./...

# Verificar código
go vet ./...

# Rodar local
go run cmd/api/main.go

# Regenerar Swagger
swag init -g cmd/api/main.go

# Docker
docker compose up -d          # Subir
docker compose down           # Parar
docker compose build --no-cache  # Rebuildar
docker logs golang_api_app    # Logs da API
```

---

## ⚙️ Variáveis de Ambiente

| Variável | Default | Descrição |
|----------|---------|-----------|
| `APP_ENV` | `development` | `development` = emails logados no console |
| `APP_PORT` | `8080` | Porta do servidor |
| `DB_HOST` | `localhost` | Host MySQL |
| `DB_PORT` | `3306` | Porta MySQL |
| `DB_USER` | `root` | Usuário MySQL |
| `DB_PASSWORD` | *(vazio)* | Senha MySQL |
| `DB_NAME` | `golang_api` | Nome do banco |
| `JWT_SECRET` | *(obrigatório)* | Chave para assinar tokens |
| `SMTP_*` | — | Config de email (só usado em produção) |

---

## 🧠 Conceitos Go (para ex-PHPs)

| Símbolo | Nome | Função | Analogia PHP |
| :--- | :--- | :--- | :--- |
| `&` | Endereço | Pega o local da memória onde o valor está | ID único de um objeto |
| `*` | Ponteiro | Variável guarda um endereço, não o valor | Objetos em classes PHP |
| `nil` | Nulo | Valor zero para ponteiros e interfaces | `null` |
| `:=` | Declaração curta | Cria e atribui inferindo o tipo | `$var = ...` |
