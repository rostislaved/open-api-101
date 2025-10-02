https://github.com/oapi-codegen/oapi-codegen

Установка:
```sh
  go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

Заметки:
- Для strict server в разных папках (chi, gin, http, etc) отличаются только main функции и мидлвары валидации. Остальное всё одинаково (handlers, usecases).
- В этом генераторе нет валидации из коробки, но её можно добавить, как middleware, который в другом пакете.
- Для инициализации этого middleware в его конструктор нужно передать OAD файл. Это можно сделать или заново распарсив файл или сгенерировав код сервера с флагом -generate spec, который добавляет метод GetSwagger().
- Мидлвар валидации для net/http внутри себя использует gorilla/mux
- Мидлвар валидации добавлен только для net/http, chi и gin. Для других мидлвары можно найти тут: https://github.com/oapi-codegen/oapi-codegen#requestresponse-validation-middleware
- Нельзя называть типы в виде operationId+Request\Response. Генератор создает те же названия в клиенте и происходит redeclaration, чтобы это избежать есть workaround (https://github.com/oapi-codegen/oapi-codegen/issues/1474). Он применен в конфиге клиента.
- Чтобы ClientWithResponse считывал response в хендлерах сервера надо явно выставлять Content-Type application/json. Легко пропустить это.
