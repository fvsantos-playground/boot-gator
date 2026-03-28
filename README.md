# boot-gator

## Conexão ao banco de dados

```bash
psql "postgres://gator:gator_pswd@192.168.1.27:42104/gator"
```

## Migrations

As migrations são gerenciadas com o [goose](https://github.com/pressly/goose).

### Instalação

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Executar migrations

```bash
cd sql/schema
goose postgres "postgres://gator:gator_pswd@192.168.1.27:42104/gator" up
```

### Reverter migrations

```bash
cd sql/schema
goose postgres "postgres://gator:gator_pswd@192.168.1.27:42104/gator" down
```

### Ver status das migrations

```bash
cd sql/schema
goose postgres "postgres://gator:gator_pswd@192.168.1.27:42104/gator" status
```

## sqlc

O [sqlc](https://sqlc.dev/) gera código Go tipado a partir de queries SQL.  
Em vez de escrever código de acesso ao banco manualmente, você escreve SQL puro e o sqlc cria as funções e structs correspondentes automaticamente.

### Instalação

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Configuração

O projeto já possui o arquivo `sqlc.yaml` na raiz com as seguintes definições:

```yaml
version: '2'
sql:
  - schema: 'sql/schema' # onde ficam as migrations (DDL)
    queries: 'sql/queries' # onde ficam as queries SQL anotadas
    engine: 'postgresql'
    gen:
      go:
        out: 'internal/database' # diretório de saída do código gerado
```

### Como escrever queries

As queries ficam em `sql/queries/`. Cada query precisa de um comentário de anotação no topo indicando o nome da função e o tipo de retorno:

| Anotação    | Descrição                                |
| ----------- | ---------------------------------------- |
| `:one`      | Retorna uma única linha (`T, error`)     |
| `:many`     | Retorna múltiplas linhas (`[]T, error`)  |
| `:exec`     | Executa sem retorno (`error`)            |
| `:execrows` | Retorna linhas afetadas (`int64, error`) |

**Exemplo** (`sql/queries/users.sql`):

```sql
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4)
RETURNING *;
```

### Gerar o código

Execute na raiz do projeto:

```bash
sqlc generate
```

O código gerado é criado em `internal/database/` e **não deve ser editado manualmente**. Sempre altere a query SQL e regenere.

### Fluxo de trabalho

```
Escrever/alterar query em sql/queries/*.sql
        ↓
sqlc generate
        ↓
Usar a função gerada em internal/database/
```
