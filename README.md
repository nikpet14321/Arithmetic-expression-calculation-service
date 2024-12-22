# Arithmetic-expression-calculation-service

- это веб-сервис на Go, который вычисляет арифметические выражения, переданные пользователем в HTTP-запросе.


Описание:

Сервис принимает POST-запрос на эндпоинт:
```
/api/v1/calculate
```

в формате JSON, где тело запроса имеет вид:
```
{
  "expression": "арифметическое выражение"
}
```

Возможный результат:

1. Если все было корректно:
    В ответ сервис возвращает 200 OK + JSON вида:
  ```
  {
    "result": "числовое значение"
  }
  ```
2. Eсли в выражении содержатся некорректные символы, лишние операторы, или происходит деление на ноль:
    422 Unprocessable Entity + JSON вида:
   ```
   {
    "error": "Expression is not valid"
   }
   ```
3. Eсли возникла непредвиденная ошибка (например, сбой в работе сервера):
    500 Internal Server Error + JSON вида:
   ```
   {
    "error": "Internal server error"
   }
   ```
   
Пример правильного использования на Linux:

Откройте терминал (Terminal) и введите команду целиком (можно в одну строку)

```
curl --location 'http://localhost:8080/api/v1/calculate' \
  --header 'Content-Type: application/json' \
  --data '{
    "expression": "2+2*2"
  }'
```
(Надеюсь работает... у меня винда)


Пример правильного использования на Windows:

1. Откройте PowerShell и введите команду, но изменённую под Windows

```
Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/calculate" `
  -Method POST `
  -Body '{"expression": "2+2*2"}' `
  -ContentType "application/json"

```
2. Откройте cmd (Win + R) и введите команду, но изменённую под Windows
```
curl --location "http://localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\":\"2+2*2\"}"
```


