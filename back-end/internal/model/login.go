package model

type Login struct {
	UserID       uint   `gorm:"primaryKey;column:user_id" json:"userId"`
	PasswordHash string `gorm:"type:text;not null;column:password_hash" json:"passwordHash"`
	// LastLoginAt  *time.Time `gorm:"column:last_update_at" json:"lastLoginAt,omitempty"`
	Salt string `gorm:"type:text;not null;column:salt" json:"salt"`

	// khoa ngoai
	User *User `gorm:"constraint:OnUpdata:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID" json:"-"`
}

func (Login) TableName() string {
	return "login"
}
