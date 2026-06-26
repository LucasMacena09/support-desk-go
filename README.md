# Support Desk Go

Sistema de atendimento ao cliente fullstack, com autenticação, abertura de chamados, histórico de mensagens e um chatbot com respostas automáticas.

Desenvolvido em **Go** (backend + frontend) com **Redis** como banco de dados.

## Funcionalidades

- Cadastro e login de usuários (senha protegida com bcrypt, autenticação via JWT)
- Rotas de chamados e mensagens protegidas por login
- Chatbot que responde automaticamente às mensagens do cliente
- Histórico de chamados e mensagens salvo no Redis
- Tratamento de erros com mensagens claras (ex: campos em branco, usuário não cadastrado)
- Front-end simples em HTML para cadastro, login e chat

## Tecnologias

- Go (`net/http`)
- Redis
- JWT (`golang-jwt/jwt`)
- bcrypt (`golang.org/x/crypto/bcrypt`)

## Como rodar o projeto

### Pré-requisitos

- [Go](https://go.dev/dl/) instalado (versão 1.22+)
- [Redis](https://redis.io/docs/getting-started/) instalado e rodando

### Passo a passo

1. Clone o repositório
   ```bash
   git clone https://github.com/seu-usuario/support-desk-go.git
   cd support-desk-go
   ```

2. Copie o arquivo de variáveis de ambiente
   ```bash
   cp .env.example .env
   ```

3. Edite o `.env` com seus valores (exemplo):
   ```env
   JWT_SECRET=uma-chave-secreta-aleatoria
   REDIS_URL=redis://localhost:6379
   PORT=8080
   ```

   Se estiver usando um Redis hospedado na nuvem (Upstash, Redis Cloud, etc.), basta colar a URL de conexão fornecida pelo provedor, no formato `redis://usuario:senha@host:porta` (ou `rediss://...` para conexões com TLS).

4. Instale as dependências
   ```bash
   go mod tidy
   ```

5. Certifique-se de que o Redis está rodando, depois inicie o servidor (a partir da raiz do projeto):
   ```bash
   go run ./cmd/server
   ```

6. Acesse no navegador:
   ```
   http://localhost:8080
   ```

## Rotas da API

| Método | Rota                          | Protegida | Descrição                          |
|--------|-------------------------------|-----------|-------------------------------------|
| POST   | `/api/register`               | Não       | Cadastra um novo usuário            |
| POST   | `/api/login`                  | Não       | Autentica e retorna um token JWT    |
| POST   | `/api/tickets`                | Sim       | Abre um novo chamado                |
| GET    | `/api/tickets`                | Sim       | Lista os chamados do usuário logado |
| GET    | `/api/tickets/{id}`           | Sim       | Detalha um chamado                  |
| POST   | `/api/tickets/{id}/messages`  | Sim       | Envia uma mensagem (aciona o bot)   |
| GET    | `/api/tickets/{id}/messages`  | Sim       | Lista o histórico de mensagens      |

Rotas protegidas exigem o header:
```
Authorization: Bearer <token>
```

## Estrutura do projeto

```
support-desk-go/
├── cmd/server/         # ponto de entrada da aplicação
├── internal/
│   ├── auth/           # JWT, hash de senha, middleware de autenticação
│   ├── chatbot/        # lógica do chatbot
│   ├── config/         # carregamento de variáveis de ambiente
│   ├── handlers/        # handlers HTTP (rotas)
│   ├── models/         # estruturas de dados
│   └── repository/     # acesso ao Redis
├── web/
│   ├── templates/      # páginas HTML
│   └── static/         # CSS
└── .env.example
```
