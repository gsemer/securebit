package services

import "securebit/domain"

type UserService struct {
	ur domain.UserRepository
}

func NewUserService(ur domain.UserRepository) *UserService {
	return &UserService{ur: ur}
}

func (us *UserService) Create(user domain.User) (domain.User, error) {
	return us.ur.Create(user)
}

func (us *UserService) Get(username string) (domain.User, error) {
	return us.ur.Get(username)
}

func (us *UserService) Delete(user domain.User) error {
	return us.ur.Delete(user)
}
