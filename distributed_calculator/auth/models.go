package models

// User struct definition
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Calculator struct definition
type Calculator struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

// Save method for User model
func (u *User) Save() error {
    // Подключение к базе данных
    db, err := InitDB()
    if err != nil {
        return err
    }
    defer db.Close()

    // Подготовка SQL запроса
    sqlStmt := `
        INSERT INTO users (login, password)
        VALUES (?, ?);
    `

    stmt, err := db.Prepare(sqlStmt)
    if err != nil {
        return err
    }

    // Выполнение SQL запроса
    _, err = stmt.Exec(u.Login, u.Password)
    if err != nil {
        return err
    }

    return nil
}

// CheckPassword method for User model
func (u *User) CheckPassword(password string) bool {
	// Implement the password checking logic
}

unc (c *Calculator) Save() error {
    // Подключение к базе данных
    db, err := InitDB()
    if err != nil {
        return err
    }
    defer db.Close()

    // Подготовка SQL запроса
    stmt, err := db.Prepare("INSERT INTO calculations (user_id, expression, result) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }

    // Выполнение SQL запроса
    _, err = stmt.Exec(c.UserID, c.Expression, c.Result)
    if err != nil {
        return err
    }

    return nil
}
func NewUser(input models.User) (*User, error) {
    // Создание нового пользователя
    newUser := &User{
        Login:    input.Login,
        Password: input.Password,
    }

    // Сохранение пользователя в базе данных
    err := newUser.Save()
    if err != nil {
        return nil, err
    }

    return newUser, nil
}

func FindUserByLogin(login string) (*User, error) {
    // Подключение к базе данных
    db, err := InitDB()
    if err != nil {
        return nil, err
    }
    defer db.Close()

    // Подготовка SQL запроса
    row := db.QueryRow("SELECT id, login, password FROM users WHERE login = ?", login)

    // Инициализация переменных для хранения данных пользователя
    var id int
    var foundLogin, password string

    // Извлечение данных из результата запроса
    err = row.Scan(&id, &foundLogin, &password)
    if err != nil {
        return nil, err
    }

    // Создание объекта пользователя
    user := &User{
        ID:       id,
        Login:    foundLogin,
        Password: password,
    }

    return user, nil
}
