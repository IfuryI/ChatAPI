# Тестовое задание для стажировки в Avito

### Для запуска нужно ввести команды:
1. sudo docker build --no-cache -t avito .
2. sudo docker run -p 9000:9000 avito

Помимо основного задания я добавил немного тестов и небольшое логирование. 

Также в репозитории пристутствует конфиг для Github Actions, в котором я запускаю линтер и тесты. (Можно посмотреть на прохождение тестов и линтеров во вкладке Actions в репозитории)

## Пример использования:
### Сборка и запуск сервера:
    >> sudo docker build --no-cache -t avito .
    ....
    >> sudo docker run -p 9000:9000 avito
    ....
### В новой вкладке терминала:
##### Создание пользователей:
    >> curl --header "Content-Type: application/json" \ 
    --request POST \
    --data '{"username": "user_227"}' \
    http://localhost:9000/users/add
    {"id":"9"}
    
    >> curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"username": "user_229"}' \
    http://localhost:9000/users/add
    {"id":"11"}
    
##### Создание чата и отправка сообщения:
    >> curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "chat_5", "users": ["11", "9"]}' \  
    http://localhost:9000/chats/add
    {"id":"10"}
    
    >> curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"chat": "10", "author": "11", "text": "hi ALL"}' \
    http://localhost:9000/messages/add
    {"id":"12"}
    
    >> curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"chat": "10", "author": "9", "text": "hi 11!"}' \
    http://localhost:9000/messages/add
    {"id":"13"}
    
##### Получение чатов пользователя и сообщений чата:

    >> curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"user": "11"}' \       
    http://localhost:9000/chats/get
    [{"chat_id":"10","name":"chat_5","users":["9","11"],"created_at":"2021-07-24T19:54:48.133668Z"}]%

    >> curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"chat": "10"}' \                                 
    http://localhost:9000/messages/get
    [{"message_id":"13","chat":"10","author":"9","text":"hi 11!","created_at":"2021-07-24T19:56:27.734977Z"},{"message_id":"12","chat":"10","author":"11","text":"hi ALL","created_at":"2021-07-24T19:55:13.215147Z"}]
