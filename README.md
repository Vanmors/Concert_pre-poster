### Concert pre-poster

Приложение, которое помогает взаимодействовать слушателю и музыканту.  
Зачастую музыкант не знает соберёт ли он достаточное количество людей в определённом городе, чтобы дать концерт, или нет.  
Данное приложение даёт возможность музыканту создавать пред-афишу и добавить в неё пару дат, в которые есть возможность организовать концерт.  
Пользователь же голосует за даты, которые ему удобны и указывает цену, которую он готов отдать за билет.  
Во время проведения голосования музыкант видит рнезультаты голосования и может оценить, ехать ему в данный город с концертом или нет.

Было реализовано REST API с помощью фреймворка Gorilla. В качестве базы данных используется PostgreSQL, в ней хранятся данные созданных афиш, голосований и результатов голосований. Для создания интерфесов пользователя была использована шаблонизация в Go.  
Переде использованием пользователь должен зарегистрироваться. После этого ему будет предоставлен «session_id», необходимый для идентификации. Приложение использует файлы cookie для аутентификации и Redis для их хранения.

### Endpoints

/billborads - вывод списка мероприятий  
/role - выбор роли  
/make_vote/{number} - предоставляет пользователю список дат на выбор и возможность ввести сумму билета, number - номер пред-афиши  
/create_voting/{number} - предоставляет артисту возможность создать голосование для определённой даты  
/create_billboard - предоставляет артисту возможность создать пред-афишу для своего концерта  
/result_voting/{number} - выводит для артиста список результатов по определённой пред-афише


### План по рефакторингу приложения

#### 1ая фаза:
1) Отоборажать голосование на стороне пользователя (Done, для успещного сохраннения нужен созданный пользователь)
2) Исправить редирект на кнопке home (Done)
3) Добавить логирование и обработку ошибок. (Done)
4) Добавить миграции в проект (Done)

#### 2ая фаза:
1) Рефакторинг авторизации. Сделать четкое разделение регистрации для юзера и артиста.  (Done)

#### 3ая фаза:
1) Добавить микросервис для создания и хранения статей и комментов к ним. Общение между сервисами через GRPC контракт (Done)
