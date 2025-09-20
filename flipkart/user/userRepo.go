package user

import "sync"

type UserRepo struct {
	Mu       sync.Mutex
	UserRepo map[string]*User
}

func CreateNewUserRepo() *UserRepo {
	return &UserRepo{
		UserRepo: make(map[string]*User),
	}
}

func (userRepo *UserRepo) Register(user *User) {
	userRepo.Mu.Lock()
	defer userRepo.Mu.Unlock()
	userRepo.UserRepo[user.UserID] = user
}
