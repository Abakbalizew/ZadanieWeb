Инструкции по запуску и тестированию приложения:

1. Создаем базу данных с названием zadaniedb:

CREATE database zadaniedb

2. Создаём таблицу users:

CREATE TABLE users (
    userid UUID PRIMARY KEY,
    email CHARACTER VARYING (30) NOT NULL,
    passwordhash BYTEA NOT NULL,
    role_ CHARACTER VARYING (30) CHECK(role_ IN ('Author', 'Reader')) NOT NULL,
	accessToken TEXT,
    refreshToken TEXT
);

3. Создаём таблицу posts:

CREATE TABLE posts (
    postid UUID PRIMARY KEY,
    authorid UUID REFERENCES users(userid) ON DELETE CASCADE,
    idempotencykey TEXT UNIQUE NOT NULL,
    title CHARACTER VARYING (100) NOT NULL,
    content TEXT NOT NULL,
    createdat TIMESTAMPTZ DEFAULT NOW(),
    updatedat TIMESTAMPTZ DEFAULT NOW(),
    status CHARACTER VARYING (10) CHECK(status IN ('Draft', 'Published')) NOT NULL
);

4. Создаём таблицу images:

CREATE TABLE images (
    imageid UUID PRIMARY KEY,
    postid UUID REFERENCES posts(postid) ON DELETE CASCADE,
    imageurl TEXT NOT NULL,
    createdat TIMESTAMPTZ DEFAULT NOW()
);

5. Запускаем с помощью 'go run main.go'. Теперь можно перейти на автоматически сгенерированную
документацию по адресу 'http://localhost:8080/swagger/index.html'

6. Главная страница сайта находится по адресу 'http://localhost:8080/api'

7. Во вкладке Зарег/Войти можно авторизоваться, или, 
нажав на соответствующую ссылку под формой, 
перейти на вкладку регистрации и зарегистрироваться.

8. Вас перенесёт на главную страницу, однако теперь вы будете аутентифицированы.
После перехода на вкладку постов, если вы автор, можете опубликовать пост или добавить его
в черновик. При редактировании поста (редактировать пост можно, если он опубликован или
находится в черновике) можно добавить картинку с помощью url-адреса.
Удалить пост можно также в редактировании поста.

Примечание: написан юнит-тест, который проверяет правильность создания токенов, тест также
проверяет работу функции, которая парсит токены (создаётся 2 токена: просроченный и валидный,
токены парсятся и возвращается либо строка "Valid", либо "Not valid").
Тест находится в jwttokens/jwttokens_test.go