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

### Variáveis de ambiente

O goose suporta variáveis de ambiente para evitar repetir a string de conexão e o driver em cada comando. O projeto já as define em `.envrc` (carregado automaticamente pelo [direnv](https://direnv.net/)):

| Variável               | Valor no projeto                                      | Descrição                                        |
| ---------------------- | ----------------------------------------------------- | ------------------------------------------------ |
| `GOOSE_DRIVER`         | `postgres`                                            | Driver do banco de dados                         |
| `GOOSE_DBSTRING`       | `postgres://gator:...@192.168.1.27:42104/gator`       | String de conexão (DSN)                          |
| `GOOSE_MIGRATION_DIR`  | `sql/schema`                                          | Diretório onde ficam os arquivos de migration    |
| `GOOSE_TABLE`          | `goose_db_version`                                    | Tabela interna de controle de versão do goose    |

Com essas variáveis definidas, **não é necessário passar o driver, a DSN ou o diretório** na linha de comando — o goose os lê automaticamente.

### Executar migrations

```bash
goose up
```

### Reverter migrations

Para reverter **apenas a última migration** aplicada (uma por vez):

```bash
goose down
```

### Criar nova migration

```bash
goose create -s <nome_da_migration> sql
```

> **Exemplo:** `goose create -s add_feeds_table sql`  
> A flag `-s` usa **numeração sequencial**, criando `00002_add_feeds_table.sql` — padrão do projeto (ex: `001_users.sql`).
> Sem ela, o goose usa timestamp como prefixo: `20170506082420_add_feeds_table.sql`.

#### Recomendações ao escrever migrations

**Nomes de arquivo**
- **Use sempre a flag `-s`** ao criar migrations: ela gera prefixos sequenciais (`00001_`, `00002_`) em vez de timestamps, mantendo a ordem explícita e legível no repositório.
- **Nomes descritivos:** o sufixo do arquivo deve deixar claro o que a migration faz (ex: `add_feeds_table`, `drop_legacy_column`), não nomes genéricos como `update` ou `change`.
- **Nunca renomeie um arquivo já aplicado:** o goose rastreia as migrations pelo nome. Renomear causa inconsistências entre ambientes.

**Estrutura do arquivo**
- **`-- +goose Up` é obrigatório:** cada arquivo deve ter exatamente uma anotação `-- +goose Up`. Todo SQL abaixo dela é executado no `up`.
- **`-- +goose Down` é opcional, mas recomendado:** todo SQL abaixo dela é executado no `down`. Sem ela, o rollback daquela migration não faz nada.
- **Uma mudança por migration:** evite agrupar alterações não relacionadas. Migrations menores são mais fáceis de depurar e reverter.
- **Nunca edite uma migration já aplicada:** se precisar corrigir algo, crie uma nova migration.

**Comandos úteis**
- **`goose validate`** — valida os arquivos de migration sem executar nada. Útil para CI/CD e para verificar antes de aplicar.
- **`goose fix`** — converte prefixos de timestamp para numeração sequencial. Útil se migrations com timestamps foram criadas acidentalmente.

### Ver status das migrations

```bash
goose status
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
