# CDTO #

### CLI Комманды

#### Добавление пользователя
 
 Для добавления пользователя наберите комманду:
```
go run main.go --config=./config/config.json useradd -e exmaple@exmaple.com  -ln Ivanov -fn Ivan -mn Ivanovich -r manager,lecturer
```
#### Тестирование email рассылки

 Для отправки тестового письма наберите комманду:

```
go run main.go --config ./config/config.json smtp-test example@example.com
```
#### Добавление тестовых данных в БД

 Добавляет тестовые данные в таблицу
 Все данные добавляются одной транзакцией (все или ничего)

```
make add_test_data
```
Список добавляемых данных:
 - Manager (manager@cdto.ru)
 - Moderator (moderator@cdto.ru)
 - WorkGroup
 - Student (student@cdto.ru)
 - A two events for the studyGroup with id=1
 - A two lecturers (lecturer@cdto.ru & lecturer2@cdto.ru
 - Links event-lecturer, event2-lecturer, event2-lecturer2
 - User (user@cdto.ru) user with no roles, used just as organizer of the first module
 - Track, themes, topics, sections
All users have the 'qwe123' passwords

### Swagger

Swagger — это фрэймворк для описания, документирования и визуализации REST API.

#### Доступ к документации

Для того, чтобы получить доступ к документации, необходимо перейти по ссылке вида `https://{{ .Domen }}/api/swagger/index.html` и ввести авторизационные данные 

       Логин: AreSpIelDOwRIETfAN
       Пароль: 2Nk7WFD%rm4g5y3e-r-5

Документация доступна на всех окружениях.

#### Генерация документации локально

Для того, чтобы сгенерировать документацию локально, необходимо выполнить команду 

```
make generate_swagger_docs
```

### Deployment

#### Prerequisites

For deploying this project Ansible is required. Ansible docs can be founded here: https://docs.ansible.com/

Currently, we have a two playbooks:

 - main.yml
 - application.yml

`main.yml` preserves common environment and deploy selected version of application.
`application.yml` just deploy selected version

#### How to deploy a project to your's local machine into docker container

You can use docker image for testing a deploy project and ensure that all works correct.
Prepare related software, ansible and docker must be installed.

Just follow next steps:
 - update permissions for correct using rsa key by `sudo chmod 0600` for ./id_rsa & ./id_rsa.pub
 - create docker container by command `make build`
 - start running scripts by command `ansible-playbook --key-file id_rsa -i inventories/dev/hosts.yml playbooks/main.yml`
 - Find container IP use command `docker inspect --format '{{ .NetworkSettings.IPAddress }}' "$@" <container_id>`
 - Go to containerIP:8000

For ssh access to the container use command `make ssh`

After deploy backEnd binary will be in /usr/local/bin folder, if you want to run it with some cli commands -
use `/usr/local/bin/{{ project_name }} {{  cli_command }}`

#### How to deploy a project to a staging server

Prepare your's ssh config file, specify host, port and user name for server which will be used for deploy.

Example of specified ssh config:
```
cat ~/.ssh/config

Host project.demo.atnr.pro
    HostName 123.123.123.123
    User root
    IdentityFile ~/.ssh/id_rsa
```

 - Write previously specified server name in the ./inventories/staging/hosts.yml
 - Start scripts by `ansible-playbook -i inventories/staging/host.yml playbooks/main.yml`
 - Go to your host in browser
 
If you want to deploy application from a custom branch (default branch is 'master') 
You can specify it in the environment variable `BITBUCKET_BRANCH`

After deploy backEnd binary will be in /usr/local/bin folder, if you want to run it with some cli commands -
use `/usr/local/bin/{{ project_name }} {{  cli_command }}`

#### How to

##### Check your service status
  >systemctl status {{ project_name }}_http_server

##### Restart a service
  >systemctl restart {{ project_name }}_http_server

##### See service logs
  >journalctl --since 10:00 -u {{ project_name }}_http_server

examples of 'since' parameter:
 - yesterday
 - today
 - 2018-10-20
 - "2018-10-20 12:48"
 note that you and server most likely in different timezones

##### Attach to database console
  > psql -h 127.0.0.1 -U {{ database_user }} {{ database_name }}

### Данные config.json

#### Скелет конфигурации

    {
      "port": {{ backend_port }},
      "logLevel": "info",
      "fqdn" : "{{ fqdn }}",
      "topicStorageLocation": "{{ topic_storage_location }}",
      "database" : {
        "name": "{{ postgres_database }}",
        "host": "{{ postgres_host }}",
        "port": {{ postgres_port }},
        "user": "{{ postgres_user }}",
        "password": "{{ postgres_password }}",
        "enableSSL": {{ postgres_enable_ssl|lower }}
      },
      "reCaptcha": {
        "secretKey": "{{ recaptcha_secret_key }}"
      },
      "jwt": {
        "secretKey": "{{ jwt_secret_key }}",
        "tokenExpHours": {{ jwt_token_exp_hours }}
      },
      "smtp": {
        "host": "{{ smtp_host }}",
        "port": {{ smtp_port }},
        "from": "{{ smtp_from }}",
        "username" : "{{ smtp_username }}",
        "password": "{{ smtp_password }}",
        "enableTLS": {{ smtp_enable_tls|lower }}
      },
      "revisionFields": {{ revision_fields | to_json }},
      "googlePlayMobileAppUrl" : "{{ google_play_url }}",
      "appStoreMobileAppUrl" : "{{ app_store_url }}",
      "emailWorkersCount" : {{ email_workers_count }},
      "googlePlayMobileAppID" : "{{ google_play_id }}",  - id приложения в play market
      "appStoreMobileAppID" : "{{ app_store_id }}",  - id приложения в itunes
      "appMonstaToken" : "{{ app_monsta_token }}",  - Сервис предоставляет доступ к API с информацией о приложениях в play market и itunes. Для получения бесплатного API с ограниченным количеством запросов необходимо в форме по ссылке ввести email и токен придет на почту. https://appmonsta.com/dashboard/get_api_key/
      "scheduleUpdateMobileVersion" : "{{ schedule_version_update }}",  - настройка частоты обновления версии приложения по крону
      "postTrackNumberUrl": "https://www.pochta.ru/tracking#"  -Ссылка на отслеживание посылки на сайте почты России, где после # подставляется трек номер
    }
    
#### Cron manager settings 
   
    Field name   | Mandatory? | Allowed values  | Allowed special characters
    ----------   | ---------- | --------------  | --------------------------
    Seconds      | Yes        | 0-59            | * / , -
    Minutes      | Yes        | 0-59            | * / , -
    Hours        | Yes        | 0-23            | * / , -
    Day of month | Yes        | 1-31            | * / , - ?
    Month        | Yes        | 1-12 or JAN-DEC | * / , -
    Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
    
You may use one of several pre-defined schedules in place of a cron expression.

    Entry                  | Description                                | Equivalent To
    -----                  | -----------                                | -------------
    @yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
    @monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
    @weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
    @daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
    @hourly                | Run once an hour, beginning of hour        | 0 0 * * * *


### Preprod окружение
 Preprod (препродакшен) - окружение, максимально повторяющее Production окружение.
 
 Отличительной особенностью данного окружения является то, что Препрод не содержит персональных данных пользователей.
 Это позволяет допускать к тестированию сотрудников, не обладающим специальным доступом.

 Детали по окружению [читайте здесь (./etc/preprod/README.md)](./etc/preprod/README.md)
 
 
 ### Развертывание проекта локально
 ##### Развертывание в среде Linux
 ###### Предварительные настройки
 Установить пакет yarn (для деплоя фронта)
 https://yarnpkg.com/lang/en/docs/install/#debian-stable
 
 Или выбрать систему из списка.
 
 Для деплоя бэка нужно установить go1.12.7 (на более ранних и поздних версиях система не работает, потому что используется плагин КАС собранный именно в этой версии).
 Скачать старую версию go можно
 ```
wget https://dl.google.com/go/go1.12.7.linux-amd64.tar.gz
```
 Распаковать архив 
```
sudo tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
```
Установить path на папку с проектом
  
Сам проект лежит в папке /src/cdto-platform
Установить PATH в файле /etc/profile
В конец файла дописать 
```
export GOPATH="/home/${user}/go"
export PATH=$PATH:/usr/local/go/bin
```
Вместо переменной прописать имя пользователя и соответствующий путь.

Установить докер последней версии:
 ```
 sudo apt-get update
 //удалит старые версии докера, если есть
 sudo apt-get remove docker docker-engine docker.io
 sudo apt install docker.io
 sudo systemctl start docker
```
 Проверьте, что докер установлен правильно docker --version
 Поскольку по умолчанию докер ставится с root правами, могут возникнуть конфликты. Необходимо снять запуск докера через рут и запускать от текущего пользователя.
 
Смотри тут -> https://itsecforu.ru/2018/04/12/%D0%BA%D0%B0%D0%BA-%D0%B8%D1%81%D0%BF%D0%BE%D0%BB%D1%8C%D0%B7%D0%BE%D0%B2%D0%B0%D1%82%D1%8C-docker-%D0%B1%D0%B5%D0%B7-sudo-%D0%BD%D0%B0-ubuntu/
 
 Внимание! Требуется перезагрузка системы после смены прав!
 
 ##### Развертывание Фронта
 Склонировать проект.
 
 Зайти в папку ${workspace}/fe
 
 Выполнить команду yarn
 
 Выполнить команду yarn dev
 
 
 Настройка окружения, куда смотрит фронт осуществляется в файле vue.config.js.
 
 Переключиться на локальный бэк можно изменив параметр target на 'http://localhost:7834',
 ``` 
 '/api': {
           // target: 'https://095e3ef4-5099-4208-9632-d5c314e5e3ed.mock.pstmn.io',
           // target: 'https://cdto-platform.demo.atnr.pro',
           target: 'http://localhost:7834',
           // target: 'https://my-cdto.gspm.ranepa.ru', 
 }
```
 
 
 ##### Развертывание бэка
 После всех предварительных настроек необходимо создать базу и добавить тестовые данные
```
make recreate_db && make add_test_data
```
 
 Поднять бэк 
 ```
 make up
 ```
 В случае, если есть ошибки миграций, то запускаем команду
  ``` 
 make migrate
 ```
 
 Если вдруг конфликты миграций между ветками пересоздаем БД командой
  ```
 make recreate_db
 ```

 
 Для добавления тестовых данных запускаем команду
  ``` 
 make add_test_data
 ```
 
 Для подписания документов потребуется пакет wkhtmltopdf.  
 ```
sudo apt install wkhtmltopdf
 ```
 
 Для доступа к БД смотрите настройки в файле ${project}/config/config.json
 ```
 "database": {
     ...
   },
 ```
 
