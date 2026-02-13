// models/post.go
package models

import "time"

type Post struct {
    ID            uint
    UserID        uint
    Content       string
    ToxicityScore int
    IsFlagged     bool
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
