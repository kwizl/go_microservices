package data

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Returns all users
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

// Returns a user based on their email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order where email = $1 order by email`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Returns a user based on their id
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order where id = $1 order by id`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Changes user details
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `update users set
		email = 1$
		first_name = 2$
		last_name = 3$
		user_active = 4$
		updated_at = 5$
		where id = $6`

	_, err := db.ExecContext(ctx, query, u.Email, u.FirstName, u.LastName, u.Active, u.UpdatedAt)

	return err
}

// Returns a user based on their email
func (u *User) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from users where id = $1`
	_, err := db.ExecContext(ctx, query, id)

	return err
}

// Create a user
func (u *User) Insert(user *User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	if err != nil {
		return 0, err
	}

	query := `insert into users (first_name, last_name, email, user_active, password, created_at) 
	values ($1, $2, $3, $4, $5, $6) returning *`

	row := db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Active, hashedPassword, user.CreatedAt)
	exe := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user.ID, exe
}

// Changes user details
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	query := `update users set
		password = 1$
		where id = $2`

	_, exe := db.ExecContext(ctx, query, hashedPassword, u.ID)

	return exe
}

// // Changes user details
// func (u *User) PasswordMatches(plainText string) (bool, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)

// 	if err != nil {
// 		return false, err
// 	}

// 	query := `select password from users where id = $1`

// 	var user User
// 	row := db.QueryRowContext(ctx, query, u.ID)

// 	exe := row.Scan(
// 		&user.ID,
// 		&user.FirstName,
// 		&user.LastName,
// 		&user.Password,
// 		&user.Active,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 	)

// 	if exe != nil {
// 		return false, exe
// 	}

// 	conv := []byte(user.Password)
// 	if conv != hashedPassword {
// 		return false, exe
// 	}

// 	return true, nil
// }
