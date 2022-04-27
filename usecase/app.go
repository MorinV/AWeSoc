package usecase

import "AWesomeSocial/internal"

const dateLayoutIso = "2006-01-02"

type Application struct {
	repositoryRegistry internal.RepositoryRegistry
}

func New(repositoryRegistry internal.RepositoryRegistry) *Application {
	return &Application{repositoryRegistry: repositoryRegistry}
}
