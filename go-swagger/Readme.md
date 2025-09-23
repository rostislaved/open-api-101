Есть несколько вариантов реализовать хендлеры:
- переопределять хендлеры в go-swagger/server/generated/restapi/configure_users_api.go
- в мейне присваивать к полям api (экземпляр созданный через operations.New***)