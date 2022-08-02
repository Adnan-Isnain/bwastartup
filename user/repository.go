package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
<<<<<<< HEAD
	FindByEmail(email string) (User, error)
=======
>>>>>>> a898100e3f8bb92bcb65140d46c0a4a476df9d38
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create((&user)).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.Where("email = ? ", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
