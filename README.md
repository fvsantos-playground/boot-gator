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
