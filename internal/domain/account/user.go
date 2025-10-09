package account

import (
	"github.com/abdelrahman146/kyora/internal/utils"

	"gorm.io/gorm"
)

const (
	UserTable  = "users"
	UserStruct = "User"
	UserAlias  = "usr"
)

type User struct {
	gorm.Model
	ID             string       `gorm:"column:id;primaryKey;type:text" json:"id"`
	FirstName      string       `gorm:"column:first_name;type:text;not null" json:"firstName"`
	LastName       string       `gorm:"column:last_name;type:text;not null" json:"lastName"`
	Email          string       `gorm:"column:email;type:text;not null;unique" json:"email"`
	PasswordHash   string       `gorm:"column:password_hash;type:text;not null" json:"-"`
	OrganizationID string       `gorm:"column:organization_id;type:text;not null" json:"organizationId"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"organization"`
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(UserAlias)
	}
	return
}

type CreateUserRequest struct {
	FirstName      string `json:"firstName" binding:"required"`
	LastName       string `json:"lastName" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
	OrganizationID string `json:"organizationId" binding:"required,uuid"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName" binding:"omitempty"`
	LastName  string `json:"lastName" binding:"omitempty"`
}
