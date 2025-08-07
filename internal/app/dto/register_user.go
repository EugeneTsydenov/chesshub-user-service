package dto

type (
	RegisterUserInputDTO struct {
		Email       string
		PublicName  string
		Password    string
		CountryCode string
		Language    string
	}

	RegisterUserOutputDTO struct {
		Message string
	}
)
