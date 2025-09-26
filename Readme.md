# open-api-101


Пример использования библиотек для генерации Go кода из [OAD](https://spec.openapis.org/oas/v3.1.0.html#openapi-document) файла (*.yaml). Подход design-first (он [рекомендуется](https://learn.openapis.org/best-practices.html)).

В репозитории примеры с библиотеками:  
- [go-swagger](https://github.com/go-swagger/go-swagger)  (spec v2.0)
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) (spec v3.0.0)
- [ogen](https://github.com/ogen-go/ogen) (spec v3.0.2)

(В скобках - версии спеки, которую реализуют. Явно не нашел, где они это заявляют, поэтому взял версию из файлов-примеров в репозитории)

### С чего начать:
1. Разобраться в структуре OAD файла (openapi.yaml), прочитав это: https://learn.openapis.org/specification/structure.html
2. Смотреть примеры в этом репозитории

Репозиторий имеет отдельную директорию для каждой библиотеки. Внутри - директории отдельно для сервера и клиента. 


### Спецификация OpenAPI:
   https://spec.openapis.org/oas/v3.1.0.html