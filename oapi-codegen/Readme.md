https://github.com/oapi-codegen/oapi-codegen

Установка:
```sh
  go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

Заметки:
- Для strict server в разных папках (chi, gin, http, etc) отличаются только main функции. Остальное всё одинаково.
- Нельзя называть типы в виде operationId+Request\Response. Генератор создает те же названия в клиенте и происходит redeclaration, чтобы это избежать есть workaround (https://github.com/oapi-codegen/oapi-codegen/issues/1474). Он применен в конфиге клиента.
- Чтобы ClientWithResponse считывал response в хендлерах сервера надо явно выставлять Content-Type application/json. Легко пропустить это.