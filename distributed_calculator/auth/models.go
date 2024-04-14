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
	// Implement the save method for the SQLite database
}

// CheckPassword method for User model
func (u *User) CheckPassword(password string) bool {
	// Implement the password checking logic
}

// Save method for Calculator model
func (c *Calculator) Save() error {
	// Implement the save method for the SQLite database
}

// NewUser function
func NewUser(input models.User) (*User, error) {
	// Implement the logic for creating a new user
}

// FindUserByLogin function
func FindUserByLogin(login string) (*User, error) {
	// Implement the logic for finding a user by login
}
