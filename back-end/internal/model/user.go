package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey; autoIncrement;column:id" json:"id"`
	Username  string `gorm:"type:varchar(50);unique;not null;column:username" json:"username"`
	Nickname  string `gorm:"type:varchar(100);column:nickname;default:user" json:"nickname"`
	AvatarURL string `gorm:"type:text;column:avatarUrl;default:https://www.svgrepo.com/show/452030/avatar-default.svg" json:"avatarUrl"`

	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`

	Login Login `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID" json:"login"`
}

func (User) TableName() string {
	return "user"
}
