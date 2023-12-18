package usersv1

import (
	"context"
	"strings"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/zitadel-data-loader/internal/config"
	"github.com/zitadel/passwap"
	"github.com/zitadel/passwap/pbkdf2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user"
	"google.golang.org/grpc"
)

func ImportUsers(ctx context.Context, cc *grpc.ClientConn, amount int) error {
	swapper := passwap.NewSwapper(pbkdf2.NewSHA256(pbkdf2.RecommendedSHA256Params))
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

	_, err := client.ImportHumanUser(ctx, newImportHumanUserRequest(swapper))
	return err
}

func newImportHumanUserRequest(swapper *passwap.Swapper) *management.ImportHumanUserRequest {
	firstName := fake.FirstName()
	lastName := fake.LastName()
	userNickName := strings.Join([]string{firstName, lastName}, "_")
	encoded, err := swapper.Hash(defaultPassword)
	if err != nil {
		panic(err)
	}

	return &management.ImportHumanUserRequest{
		UserName: userNickName,
		Profile: &management.ImportHumanUserRequest_Profile{
			FirstName:         firstName,
			LastName:          lastName,
			NickName:          userNickName,
			DisplayName:       strings.Join([]string{firstName, lastName}, " "),
			PreferredLanguage: fake.LanguageAbbreviation(),
			Gender:            user.Gender(fake.IntRange(int(user.Gender_GENDER_FEMALE), int(user.Gender_GENDER_DIVERSE))),
		},
		Email: &management.ImportHumanUserRequest_Email{
			Email:           strings.Join([]string{userNickName, fake.DomainName()}, "@"),
			IsEmailVerified: true,
		},
		HashedPassword: &management.ImportHumanUserRequest_HashedPassword{
			Value: encoded,
		},
	}
}
