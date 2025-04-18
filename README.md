# Для запуска использовать ```make up```

## 2 rest маршрута
### 1-й маршрут(get): http://localhost:8080/auth/token?guid=
### Принимает query параметр guid: uuid пользователя
### Пример запроса ```http://localhost:8080/auth/token?guid=0d44f783-e1c0-4d49-bb96-80ea8abc0efa```
### Возвращает json ввида: 
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMGQ0NGY3ODMtZTFjMC00ZDQ5LWJiOTYtODBlYThhYmMwZWZhIiwiaXBfYWRkcmVzcyI6IjE3Mi4xOS4wLjE6NDk2MjIiLCJKVEkiOiIxOTA3ZWZmOS0xYTJmLTQzM2QtYmVmMi05MjMwOTVlYjE3Y2UiLCJleHAiOjE3NDQ5NDQ0MTgsImlhdCI6MTc0NDk0MzUxOH0.uhY4yZxTaNt3BKoEAaSQyIGq8KGXqqPRfPzjeVkgomRd3bBHhJp1vBXT3RSnfJgnjDnhSjxxxH2jUvw6ltSdHQ",
    "refresh_token": "N2FmMjNhMzctNzkwNy00ZWNlLWE2ODItYTE0Mjc5NDQ0ZTc4"
}


### 2-й маршрут(post):http://localhost:8080/auth/token/refresh
### Принимает тело запроса в виде json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMGQ0NGY3ODMtZTFjMC00ZDQ5LWJiOTYtODBlYThhYmMwZWZhIiwiaXBfYWRkcmVzcyI6IjE3Mi4xOS4wLjE6NDk2MjIiLCJKVEkiOiIxOTA3ZWZmOS0xYTJmLTQzM2QtYmVmMi05MjMwOTVlYjE3Y2UiLCJleHAiOjE3NDQ5NDQ0MTgsImlhdCI6MTc0NDk0MzUxOH0.uhY4yZxTaNt3BKoEAaSQyIGq8KGXqqPRfPzjeVkgomRd3bBHhJp1vBXT3RSnfJgnjDnhSjxxxH2jUvw6ltSdHQ",
    "refresh_token": "N2FmMjNhMzctNzkwNy00ZWNlLWE2ODItYTE0Mjc5NDQ0ZTc4"
}
Возращает новую пару токенов
