# Сервер с аутентификацией по access и refresh токенам. 
## Порядок установки
1. `git clone github.com/crimson-catfish/MEDODS-TestTask <destination-directory>`
2. `cd <destination-directory>`
3. Создать файл .env со следующими переменными: DB_USER, DB_PASSWORD, DB_NAME, SECRET_KEY, LOCALHOST_PORT
4. присвоить LOCALHOST_PORT любой свободный порт, остальные переменные - произвольная строка
5. `docker compose up`

## Маршруты сервера
/login
```json
{
  "guid": "sample-guid"
  "email": "example@email.com"
  "password": "very-secret-password"
}
```

/register
```json
{
  "guid": "sample-guid"
  "password": "very-secret-password"
}
```

/
```json
{
  "access_token": "aaaaaaa.bbbbbbb.ccccccc"
}
```

/refresh
```json
{
  "access_token": "aaaaaaa.bbbbbbb.ccccccc"
  "refresh_token": "$aa$bb$ccccccc"
}
```
