-- Создание таблицы для хранения выражений и ответов пользователей
CREATE TABLE user_expressions (
                                  id INTEGER PRIMARY KEY,
                                  expression TEXT,
                                  responses TEXT,
                                  user_name TEXT
);

-- Создание таблицы для хранения текущего значения счетчика ID
CREATE TABLE id_counter (
                            current_id INTEGER
);

-- Инициализация счетчика ID значением 0
INSERT INTO id_counter (current_id) VALUES (0);