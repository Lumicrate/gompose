package auth

type UserModel struct {
	ID       string `gorm:"primaryKey" json:"id" bson:"id,omitempty"`
	Email    string `gorm:"unique" json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (u *UserModel) GetID() string             { return u.ID }
func (u *UserModel) GetEmail() string          { return u.Email }
func (u *UserModel) GetHashedPassword() string { return u.Password }
