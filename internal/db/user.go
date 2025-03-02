package db

func AddUser(u *User) error {
	return db.Create(&u).Error
}

func IsUserExists(username string) (bool, error) {
	tx := db.Limit(1).Where("Username = ?", username).Find(&User{})
	return tx.RowsAffected > 0, tx.Error
}

func (u *User) Fill() error {
	if u.Username != "" {
		return db.Where("Username = ?", u.Username).First(&u).Error
	} else if u.ID != 0 {
		return db.First(&u, u.ID).Error
	}

	return ErrNoIdentifier
}

func FetchUsers(query string, limit int, page int) (*[]User, error) {
	var users *[]User
	tx := db.Limit(limit).Offset(page * limit)
	tx.Where("Username LIKE ?", query)
	tx.Find(&users)

	return users, tx.Error
}
