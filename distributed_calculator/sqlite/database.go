package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase() (*sql.DB, error) {
    // Подключение к базе данных SQLite
    db, err := sql.Open("sqlite3", "./calculator.db")
    if err != nil {
        return nil, err
    }

    // Создание таблиц для пользователей и вычислений
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            login TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS calculations (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            expression TEXT NOT NULL,
            result TEXT NOT NULL,
            FOREIGN KEY(user_id) REFERENCES users(id)
        );
    `)

    return db, err
}

func Migrate(db *sql.DB) error {
    // Примените миграции для базы данных, если необходимо
    // ...

    return nil
}

func Close(db *sql.DB) error {
    return db.Close()
}

func (u *User) Save(db *sql.DB) error {
    // Хэширование пароля
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // Сохранение пользователя в базе данных
    _, err = db.Exec("INSERT INTO users (login, password_hash) VALUES (?, ?)", u.Login, hashedPassword)

    return err
}

    // Получение хэша пароля из базы данных
    var hashedPassword []byte
    err := db.QueryRow("SELECT password_hash FROM users WHERE login = ?", u.Login).Scan(&hashedPassword)
    if err != nil {
        return false
    }

    // Сравнение хэша пароля из базы данных с введённым паролем
    return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)) == nil
}

func (c *Calculator) Save(db *sql.DB) error {
    // Получение ID пользователя
    var userID int
    err := db.QueryRow("SELECT id FROM users WHERE login = ?", c.User.Login).Scan(&userID)
    if err != nil {
        return err
    }

    // Сохранение вычисления в базе данных
    _, err = db.Exec("INSERT INTO calculations (user_id, expression, result) VALUES (?, ?, ?)", userID, c.Expression, c.Result)

    return err
}

func NewUser(db *sql.DB, input models.User) (*User, error) {
    // Создание нового пользователя
    user := &User{
        Login:    input.Login,
        Password: input.Password,
    }

    // Сохранение пользователя в базе данных
    err := user.Save(db)

    return user, err
}

func FindUserByLogin(db *sql.DB, login string) (*User, error) {
    // Получение пользователя из базы данных по логину
    var user User
    err := db.QueryRow("SELECT id, login, password_hash FROM users WHERE login = ?", login).Scan(&user.ID, &user.Login, &user.Password)
    if err != nil {
        return nil, err
    }

    return &user, nil
}
