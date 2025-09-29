# ProcrastiGo API 🦥

API для генерации оправданий прокрастинации. Когда дедлайн горит, а делать ничего не хочется!

## Быстрый старт

```bash
git clone https://github.com/your-username/procrastigo.git
cd procrastigo
go run cmd/api/main.go
```

## Использование

```bash
# случайное оправдание
curl http://localhost:8080/api/v1/excuses/random

# технические оправдания
curl "http://localhost:8080/api/v1/excuses?category=tech"

# добавить своё оправдание
curl -X POST http://localhost:8080/api/v1/excuses \
  -H "Content-Type: application/json" \
  -d '{"text":"Гит сломался", "category":"tech"}'
```

