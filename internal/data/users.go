package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"ratrace.darieldejesus.com/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousUser = &User{}

type password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// IsAnonymous returns true if the referenced user is anonymous.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// Set calculates the bcrypt hash of a plain text password and stores
// both the hash and the plain password in the struct.
func (p *password) Set(plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plainPassword
	p.hash = hash

	return nil
}

// Matches checks whether the provided plain text password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *password) Matches(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 255, "name", "must not be more than 255 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePassword(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

// isDuplicatedEmailError returns true if the given error comes
// from the unique constraint in the users table
func isDuplicatedEmailError(err error) bool {
	var mySQLError *mysql.MySQLError
	if errors.As(err, &mySQLError) {
		return mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_unique_email")
	}
	return false
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	stmt := `INSERT INTO users (name, email, password_hash, activated, created_at)
	VALUES (?, ?, ?, ?, UTC_TIMESTAMP())`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(
		ctx,
		stmt,
		user.Name,
		user.Email,
		string(user.Password.hash),
		user.Activated,
	)
	if err != nil {
		switch {
		case isDuplicatedEmailError(err):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	stmt := `SELECT id, name, email, password_hash, activated, version, created_at
	FROM users
	WHERE email = ?`

	user := &User{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (m UserModel) Update(user *User) error {
	stmt := `UPDATE users
	SET name = ?, email = ?, password_hash = ?, activated = ?, version = version + 1
	WHERE id = ? AND version = ?`

	args := []any{
		user.Name,
		user.Email,
		string(user.Password.hash),
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, args...)
	if err != nil {
		switch {
		case isDuplicatedEmailError(err):
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlainText string) (*User, error) {
	hash := sha256.Sum256([]byte(tokenPlainText))
	tokenHash := hex.EncodeToString(hash[:])

	stmt := `SELECT users.id, users.name, users.email, users.password_hash, users.activated, users.version
	FROM users
	INNER JOIN tokens ON users.id = tokens.user_id
	WHERE tokens.hash = ? AND tokens.scope = ? AND tokens.expiry > ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, stmt, tokenHash, tokenScope, time.Now()).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
