package user

import "context"

type UserUseCase struct {
	userRepo UserRepository
}

func NewUserUseCase(uRepo UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: uRepo,
	}
}

func (u *UserUseCase) GetMeUseCase(ctx context.Context, userId int) (*UserEntity, error) {
	return u.userRepo.GetFullUserById(ctx, userId)
}

func (u *UserUseCase) GetUserById(ctx context.Context, userId int) (*UserEntity, error) {
	return u.userRepo.GetFullUserById(ctx, userId)
}

func (u *UserUseCase) GetAllUsersUseCase(ctx context.Context) ([]*UserEntity, error) {
	return u.userRepo.GetAllUsers(ctx)
}

func (u *UserUseCase) DeleteUserUseCase(ctx context.Context, userId int) error {
	return u.userRepo.DeleteUser(ctx, userId)
}