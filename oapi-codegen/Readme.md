https://github.com/oapi-codegen/oapi-codegen

Install:
```sh
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

- Нельзя называть типы в виде operationId+Request\Response. Генератор создает те же названия в клиенте и проходит redeclaration
- Чтобы ClientWithResponse считывал response в хендлерах сервера надо явно выставлять Content-Type application/json. Легко пропустить это.