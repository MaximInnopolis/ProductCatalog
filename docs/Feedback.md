# Ревью проекта

## Миграции

   **Проблема:**  
   Мы выкатили сервис в Прод. Теперь мы хотим выкатить апдейт, но для него нужны изменения в схеме БД. Можно внести изменения в `migration.go`, но если что-то пойдет не так, не ясно как откатывать.

   **Задание:**  
   Изменить механизм накатывания миграций. Механизм должен поддерживать rollback и версионирование.

## База данных

   **Проблема:**
   1. В пакете `database` есть глобальная переменная `db`. Это антипаттерн, т.к. все необходимые зависимости должны передаваться в соответствующий слой. В противном случае, любая часть кода может (в теории) обратиться в БД, хотя это не её область ответственности.
   2. Можно создать две категории с одним и тем же имененм.

   **Задание:**
   1. Продумать и реализовать механизм "прокидывания" зависимостей в каждый из слоёв приложения (`api`, `service`, `storage`). Тут важно понимать "кому и что нужно", например, HTTP router нужен только слою API (должен ли он создаваться в `main` как зависимость?).
   2. Необходимо обработать этот кейс. Подумать, а должен ли пользователь получить ошибку в этом случае (`POST`/`PUT`)?

## HTTP

   **Проблема:**
   1. При `GET` запросах нет проверки авторизации, т.е. кто угодно может получить информацию из нашей системы.
   2. Метод `auth.RequireValidToken` в `api` вызывает во многих `handler-ах`, т.о. если мы добим ещё один необходимо не забыть добавить туда проверки авторизации.
   3. Пакет `auth` использует термины HTTP, хотя он не в пакете `api`. Возможно, механизм авторизации должен быть "отвязан" от конкретного типа API (`HTTP`/`gRPC`/`Socket`).

   **Задание:**  
   "Автоматизировать" проверку токена. Каждый `handler`, в идеале, вообще не должен знать что есть какая-либо проверка авторизации, т.к. это не его зона ответственности.

## Логи (на подумать)

   **Проблемы:**
   1. Сейчас запись в лог не имеет контекста, т.е. не ясно из какой части кода этот лог.
   2. У нас может быть множество вызовов одного и того же `endpoint-а`, но в логах мы не сможем понять к какому конкретно вызову сообщение.

   **Задание:**
   1. Подумать, как сделать лог более читаемым, чтобы по нему было ясно в каком месте он "возник".
   2. Как привязать запись в логе к конкретному вызову какого-либо `endpoint-а`?