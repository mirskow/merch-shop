package model

type User struct {
	ID           int    `json:"-" db:"id"`
	Username     string `json:"username" binding:"required" db:"username"`
	Coin         int    `json:"coin"`
	PasswordHash string `json:"password" binding:"required" db:"password_hash"`
}

type Merch struct {
	ID   int    `json:"-" db:"id"`
	Item string `json:"item" binding:"required" db:"item"`
	Cost int    `json:"cost" db:"cost"`
}

type Purchases struct {
	ID       int    `json:"-" db:"id"`
	UserID   int    `json:"user_id" db:"user_id"`
	ItemID   int    `json:"item_id" db:"item_id"`
	ItemName string `json:"item_type" db:"-"`
	Quantity int    `json:"quantity" db:"-"`
}

type Transaction struct {
	FromUserID   int    `json:"fromUser" db:"from_user"`
	ToUserID     int    `json:"toUser" db:"to_user"`
	Amount       int    `json:"amount" db:"amount"`
	SenderName   string `json:"sender" db:"-"`
	ReceiverName string `json:"receiver" db:"-"`
}
