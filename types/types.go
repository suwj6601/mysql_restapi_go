package types

type User struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Age       int     `json:"age"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id string) (*User, error)
	GetUserRoleByID(id string) (bool, error)
	RegisterUser(User) error
	CreateUser(User) error
	DeleteUser(userID string) error
	LoginUser(email, password string) (string, error)
}

type RegisterUserPayload struct {
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Balance  float64 `json:"balance"`
}
