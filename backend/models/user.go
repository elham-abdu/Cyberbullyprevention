// models/user.go
package models

import "time"

type User struct {
    ID           uint
    Email        string
    PasswordHash string
    Role         string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
