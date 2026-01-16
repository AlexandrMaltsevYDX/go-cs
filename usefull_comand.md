# Usfull comands

## Rebuild

```bash
go mod tidy && go build -o ./tmp/main ./cmd
```

## Run

```bash
http POST localhost:3000/api/example id:=1 example_column="Hello World"
```

**Разбор команды:**
 *`POST` — метод запроса (HTTPie выбирает его автоматически при наличии данных,
 но лучше явно указывать).
 *`localhost:3000/api/example` — URL вашего эндпоинта.
 *`id:=1` — синтаксис `:=` используется для передачи **нестроковых** значений
 (чисел, булевых, JSON-объектов). В Go структуре поле `Id` имеет тип `int`.
 *`example_column="Hello World"` — синтаксис `=` используется для строковых значений
