package optimistic

type Optimistic struct {
	Id      int64   `gorm:"column:id; primary_key; AUTO_INCREMENT" json:"id"`
	UserId  string  `gorm:"column:user_id; default:0; not null" json:"user_id"`
	Amount  float32 `gorm:"column:amount; not null" json:"amount"`
	Version int64   `gorm:"column:version; default:0; not null" json:"version"`
}
