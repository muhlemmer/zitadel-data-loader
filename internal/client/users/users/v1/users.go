package usersv1

import (
	"context"
	"strings"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/zitadel-data-loader/internal/config"
	"github.com/zitadel/passwap"
	"github.com/zitadel/passwap/bcrypt"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user"
	"google.golang.org/grpc"
)

func ImportUsers(ctx context.Context, cc *grpc.ClientConn, amount int) error {
	swapper := passwap.NewSwapper(bcrypt.New(4))
	client := management.NewManagementServiceClient(cc)

	for i := 0; i < amount; i++ {
		if err := importUser(ctx, client, swapper); err != nil {
			return err
		}
	}
	return nil
}

const defaultPassword = "Password1!"

func importUser(ctx context.Context, client management.ManagementServiceClient, swapper *passwap.Swapper) error {
	ctx, cancel := context.WithTimeout(ctx, config.Global.Timeout)
	defer cancel()

	firstName := fake.FirstName()
	lastName := fake.LastName()
	encoded, err := swapper.Hash(defaultPassword)
	if err != nil {
		return err
	}

	_, err = client.ImportHumanUser(ctx, &management.ImportHumanUserRequest{
		UserName: fake.Username(),
		Profile: &management.ImportHumanUserRequest_Profile{
			FirstName:         fake.FirstName(),
			LastName:          fake.LastName(),
			NickName:          fake.PetName(),
			DisplayName:       strings.Join([]string{firstName, lastName}, " "),
			PreferredLanguage: fake.LanguageAbbreviation(),
			Gender:            user.Gender(fake.IntRange(int(user.Gender_GENDER_FEMALE), int(user.Gender_GENDER_DIVERSE))),
		},
		Email: &management.ImportHumanUserRequest_Email{
			Email:           fake.Email(),
			IsEmailVerified: true,
		},
		HashedPassword: &management.ImportHumanUserRequest_HashedPassword{
			Value: encoded,
		},
	})
	return err
}
