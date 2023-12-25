package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Permissions []string

// Include returns whether the given code exists in the
// permissions slice.
func (p Permissions) Include(code string) bool {
	for _, v := range p {
		if v == code {
			return true
		}
	}
	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	stmt := `SELECT p.permission
	FROM permissions p
	INNER JOIN users_permissions up ON up.permission_id = p.id
	INNER JOIN users u ON u.id = up.user_id
	WHERE u.id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	permissions := Permissions{}

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m PermissionModel) AddForUser(userID int64, permissions ...string) error {
	numPermissions := len(permissions)
	inParams := strings.Repeat(",?", numPermissions-1)

	stmt := fmt.Sprintf(`INSERT INTO users_permissions
	SELECT ?, permissions.id
	FROM permissions
	WHERE permissions.permission IN (?%s)`, inParams)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := make([]any, 0, numPermissions+1)
	args = append(args, userID)
	for i := range permissions {
		args = append(args, permissions[i])
	}

	_, err := m.DB.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}

	return nil
}
