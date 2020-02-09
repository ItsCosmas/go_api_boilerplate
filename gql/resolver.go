package gql

import (
	"context"
	"errors"
	"go_api_boilerplate/domain/user"
	"go_api_boilerplate/gql/gen"
	"go_api_boilerplate/services/authservice"
	"go_api_boilerplate/services/userservice"
)

// Resolver struct
type Resolver struct {
	UserService userservice.UserService
	AuthService authservice.AuthService
}

// Mutation graphql
func (r *Resolver) Mutation() gen.MutationResolver {
	return &mutationResolver{r}
}

// Query graphql
func (r *Resolver) Query() gen.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

// ================ Query ================

func (r *queryResolver) User(ctx context.Context, id int) (*gen.User, error) {
	user, err := r.UserService.GetByID(uint(id))
	if err != nil {
		return nil, err
	}

	return &gen.User{
		ID:        int(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		Active:    user.Active,
	}, nil
}

func (r *queryResolver) UserProfile(ctx context.Context) (*gen.User, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, errors.New("Unauthorized: Token is invlaid")
	}

	user, err := r.UserService.GetByID(userID.(uint))
	if err != nil {
		return nil, err
	}

	return &gen.User{
		ID:        int(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		Active:    user.Active,
	}, nil
}

// ================ Mutation ================

func (r *mutationResolver) Login(ctx context.Context, input gen.RegisterLogin) (*gen.RegisterLoginOutput, error) {
	user, err := r.UserService.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	err = r.UserService.ComparePassword(input.Password, user.Password)
	if err != nil {
		return nil, err
	}

	token, err := r.AuthService.IssueToken(*user)
	if err != nil {
		return nil, err
	}

	return &gen.RegisterLoginOutput{
		Token: token,
		User: &gen.User{
			ID:        int(user.ID),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			Active:    user.Active,
		},
	}, nil
}

func (r *mutationResolver) Register(ctx context.Context, input gen.RegisterLogin) (*gen.RegisterLoginOutput, error) {
	userDomain := &user.User{
		Email:    input.Email,
		Password: input.Password,
	}

	err := r.UserService.Create(userDomain)
	if err != nil {
		return nil, err
	}

	token, err := r.AuthService.IssueToken(*userDomain)
	if err != nil {
		return nil, err
	}

	return &gen.RegisterLoginOutput{
		Token: token,
		User: &gen.User{
			ID:        int(userDomain.ID),
			FirstName: userDomain.FirstName,
			LastName:  userDomain.LastName,
			Email:     userDomain.Email,
			Role:      userDomain.Role,
			Active:    userDomain.Active,
		},
	}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input gen.UpdateUser) (*gen.User, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, errors.New("Unauthorized: Token is invlaid")
	}

	usr, err := r.UserService.GetByID(userID.(uint))
	if err != nil {
		return nil, err
	}

	if input.Email != "" {
		usr.Email = input.Email
	}
	if input.FirstName != nil {
		usr.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		usr.LastName = *input.LastName
	}

	err = r.UserService.Update(usr)
	if err != nil {
		return nil, err
	}

	return &gen.User{
		ID:        int(usr.ID),
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Email:     usr.Email,
		Role:      usr.Role,
		Active:    usr.Active,
	}, nil
}