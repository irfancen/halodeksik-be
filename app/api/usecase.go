package api

type AllUseCases struct{}

func InitializeUseCases(allRepo *AllRepositories) *AllUseCases {
	return &AllUseCases{}
}
