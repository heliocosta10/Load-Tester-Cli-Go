# Load Tester em Go

Este projeto é um **CLI (Command Line Interface) em Go** que desenvolvi para realizar testes de carga em serviços web.

A aplicação recebe uma URL, a quantidade total de requisições e o nível de concorrência. Depois disso, executa as requisições HTTP e apresenta um relatório com o tempo total, quantidade de requisições realizadas e a distribuição dos códigos de status HTTP.

## Objetivo

O objetivo deste desafio foi criar uma ferramenta simples para testar a capacidade de resposta de um serviço web utilizando requisições simultâneas.

Os parâmetros utilizados são:

- `--url`: URL do serviço que será testado.
- `--requests`: quantidade total de requisições.
- `--concurrency`: quantidade de chamadas simultâneas.

## Tecnologias

- Go
- HTTP Client
- Goroutines
- Channels
- Docker

## Estrutura do projeto

```text
load-tester-go/
├── main.go
├── main_test.go
├── go.mod
├── Dockerfile
├── .dockerignore
├── .gitignore
└── README.md
```

## Como executar localmente

Primeiro, clone o repositório e entre na pasta:

```bash
git clone https://github.com/heliocosta10/load-tester-go.git
cd load-tester-go
```

Depois, posso executar os testes:

```bash
go test ./...
```

Para executar diretamente:

```bash
go run . --url=http://google.com --requests=100 --concurrency=10
```

## Como funciona a concorrência

Para controlar a quantidade de chamadas simultâneas, utilizei goroutines e um canal de tarefas.

O número informado em `--requests` determina exatamente quantas tarefas serão criadas.

O número informado em `--concurrency` determina quantos workers ficam processando essas tarefas ao mesmo tempo.

Por exemplo:

```bash
--requests=1000 --concurrency=10
```

Nesse caso:

- serão criadas exatamente 1000 tarefas;
- serão utilizados 10 workers;
- cada worker processará as tarefas disponíveis;
- quando todas as tarefas terminarem, o programa gera o relatório.

Dessa forma, o número total de requisições não depende da quantidade de workers.

## Relatório

Ao finalizar o teste, a aplicação apresenta:

- tempo total gasto;
- quantidade de requests realizados;
- quantidade de requests solicitados;
- quantidade de respostas HTTP 200;
- quantidade de cada outro código HTTP encontrado;
- quantidade de falhas de conexão.

Exemplo:

```text
========================================
           RELATÓRIO DO TESTE
========================================
Tempo total: 2.341s
Requests realizados: 100
Requests solicitados: 100
Status HTTP 200: 96

Distribuição dos status HTTP:
HTTP 200: 96
HTTP 404: 2
HTTP 500: 1
Falha de conexão: 1
========================================
```

## Executando com Docker

Primeiro, construir a imagem:

```bash
docker build -t load-tester-go .
```

Depois, executar o teste:

```bash
docker run --rm load-tester-go --url=http://google.com --requests=1000 --concurrency=10
```

O formato exigido pelo desafio também pode ser usado desta forma:

```bash
docker run load-tester-go --url=http://google.com --requests=1000 --concurrency=10
```

## Testando uma aplicação local

Se o serviço que quero testar estiver rodando na minha máquina na porta 8080, quando estiver usando Docker no Windows, posso usar:

```bash
docker run --rm load-tester-go --url=http://host.docker.internal:8080 --requests=100 --concurrency=10
```

O `host.docker.internal` permite que o container acesse um serviço que está rodando no host.

## Exemplos

### 100 requisições com 10 chamadas simultâneas

```bash
docker run --rm load-tester-go --url=http://google.com --requests=100 --concurrency=10
```

### 1000 requisições com 20 chamadas simultâneas

```bash
docker run --rm load-tester-go --url=http://google.com --requests=1000 --concurrency=20
```

### 5000 requisições com 50 chamadas simultâneas

```bash
docker run --rm load-tester-go --url=http://google.com --requests=5000 --concurrency=50
```

## Validação dos parâmetros

A aplicação impede a execução quando os parâmetros obrigatórios estão incorretos.

Exemplos:

```bash
go run . --requests=100 --concurrency=10
```

Retorno:

```text
Erro: o parâmetro --url é obrigatório
```

Também não são aceitos valores menores ou iguais a zero para `--requests` e `--concurrency`.

## Testes automatizados

Para executar os testes:

```bash
go test ./...
```

Os testes verificam principalmente a validação dos parâmetros:

- aceita parâmetros válidos;
- exige a URL;
- exige `--requests` maior que zero;
- exige `--concurrency` maior que zero.

## Docker: build e execução

O Dockerfile utiliza duas etapas.

Na primeira etapa, a aplicação é compilada usando a imagem oficial do Go.

Na segunda etapa, utilizo uma imagem Alpine menor apenas para executar o binário compilado.

Para criar a imagem:

```bash
docker build -t load-tester-go .
```

Para verificar a imagem:

```bash
docker images
```

Para executar:

```bash
docker run --rm load-tester-go --url=http://google.com --requests=1000 --concurrency=10
```



## Observação

O código foi desenvolvido com foco nos requisitos do desafio: execução via CLI, controle de concorrência, quantidade exata de requisições, relatório dos resultados e execução obrigatória através de Docker.
