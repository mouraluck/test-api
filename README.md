# test-api

API simples em Go com endpoint de healthcheck:

- `GET /ping` -> `pong`

## Rodar local

```bash
go run .
```

Servidor sobe na porta `8080` por padrão, ou usa a variável `PORT`.

## Testes

```bash
go test ./...
```

## Exemplo rápido

```bash
curl http://localhost:8080/ping
```

## Docker

```bash
docker build -t test-api:local .
docker run --rm -p 8080:8080 test-api:local
```

## Deploy no Render

1. Suba este repositório para o GitHub.
2. No Render, crie um novo serviço Web usando o repositório.
3. O arquivo `render.yaml` já configura o serviço com runtime Docker.
4. Após o deploy, teste:

```bash
curl https://SEU-SERVICO.onrender.com/ping
```
