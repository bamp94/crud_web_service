###1.  Web-сервис(CRUD)

Написать сервис, который будет слушать входящие запросы по HTTP, преобразовывать их в запрос к соответствующей функции Postgres (по схеме трансляции, приведённой ниже), выполнять запрос и возвращать ответ клиенту.

Как плюс: ограничить максимальное количество одновременных коннектов к БД.

Настройки соединения с сервером Postgres читать из ENV среды:

    - port - (int) порт, на котором слушать запросы
    - endpoint - (string) название API
    - host - (string) hostname, где установлен Postgres
    - user - (string) имя пользователя Postgres
    - password - (string) пароль пользователя Postgres
    - schema - (string) схема в Postgres

Трансляция запроса в вызов Postgres функции

    Формат запроса к сервису:
    HTTP_METHOD server:port/endpoint/vversion[/object/id...]]/destination/[id] , где
    HTTP_METHOD - одно из: GET, POST, PUT, DELETE
    server - сервер, где запущен веб-сервис
    port - порт
    endpoint - значение из config-файла
    version - номер версии API, число
    /object/id - необязательный повторяющийся параметр, определяющий путь в иерархии объектов
    /destination/ - конечный объект
id - id конечного объекта. Обязателен для методов PUT, DELETE, не указывается для POST. Для GET -- если указан, то возвращает элемент с данным id, если не указан, возвращает полный список элементов.

Правила трансляции

    запрос в Postgres = select * from схема.[object1[_object2]...]_destination_method( [id1[, id2]... ,] id[, params])

В зависимости от HTTP метода к имени функции добавляется cуффикс method:

    для GET - get
    для POST - ins
    для PUT - upd
    для DELETE - del
Для POST и PUT методов в теле запроса принимается JSON, который передаётся в Postgres в качестве параметра params.

Все методы должны возвращать результат работы соответствующей Postgres функции с ContentType = 'application/json'

    create schema if not exists test;
    
    create sequence if not exists test.seq_users;
    create sequence if not exists test.seq_comments;
    
    create table if not exists test.users
    (
      id int not null default nextval('test.seq_users'::regclass),
      name varchar not null,
      email varchar not null,
      constraint "PK_users" primary key (id),
      constraint "UQ_users_email" unique (email),
      constraint "CHK_users_email" check (email like '%@%')
    );

    create table if not exists test.comments
    (
      id int not null default nextval('test.seq_comments'::regclass),
      id_user int not null,
      txt varchar not null,
      constraint "PK_comments" primary key (id)
    );

#### Запуск
Запуск приложения производится коммандой 

    make up

Применить миграции 

    make migrate

Пересоздать базу данных

    make recreate_db

#### Тестирование запросов

Вот список доступных для тестирования комманд:

        curl -X POST "http://localhost:7856/api/v1/user/" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"name\": \"foo\", \"email\": \"foofoo@fo.ru\"}"
        
        curl -X PUT "http://localhost:7856/api/v1/user/1" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"name\": \"test\", \"email\": \"test@test.ru\"}"v
        
        curl -X GET "http://localhost:7856/api/v1/user/1" -H "accept: application/json" -H "Content-Type: application/json"
        
        curl -X DELETE "http://localhost:7856/api/v1/user/1" -H "accept: application/json" -H "Content-Type: application/json"

###2. Задачи

2.1 У нас есть множество каналов для чтения. Мы не знаем заранее какие из них активны а какие нет. 
Нужно реализовать функцию в которую мы можем передать любое множество каналов а в ответе должны получить 1 канал в который будут 
перенаправляться сообщения из других каналов.
2.2 На вход функция получает a = []int{-1,2,22,213123123,0,1} и исходное значение. Необходимо вернуть набор комбинаций из 3-ёх чисел 
которые в сумме дадут исходное значение. Комбинации должны быть уникальными(1,0,2) == (2,1,0)

#### Запуск
Запуск тестирования каждого задания производится из папки tasks. Для запуска программы тестированиея задачи 2.1:

        go run channels.go

Для запуска программы тестированиея задачи 2.2:

        go run combinations.go

