package entities

type Rocket struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Category string `json:"category"`
}
