# T-Bank Shariki Case

Материалы по кейсу платформы безопасной монетизации обезличенных клиентских данных:

- [Архитектурный документ](/Users/hanq/tbankshariki/docs/architecture.md)
- [Выдача данных компаниям по API](/Users/hanq/tbankshariki/docs/company-api.md)
- [Use Case диаграмма](/Users/hanq/tbankshariki/docs/diagrams/use-case.puml)
- [DFD диаграмма](/Users/hanq/tbankshariki/docs/diagrams/dfd.puml)
- [Черновик OpenAPI](/Users/hanq/tbankshariki/docs/openapi.yaml)

Документы заточены под MVP на Go и пригодны как база для финального питча.

## Go Demo API

В репозитории есть минимальный runnable backend на Go с in-memory данными и раздачей Swagger.

Запуск:

```bash
go run ./cmd/api
```

По умолчанию сервис стартует на `http://127.0.0.1:8080`.

Полезные ссылки:

- `http://127.0.0.1:8080/swagger/`
- `http://127.0.0.1:8080/docs/openapi.yaml`
- `http://127.0.0.1:8080/v1/company/products`
