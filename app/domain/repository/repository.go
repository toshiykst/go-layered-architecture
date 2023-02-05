//go:generate mockgen -source=$GOFILE -package=mock$GOPACKAGE -destination=../../mock/domain/$GOPACKAGE/$GOFILE
package repository

type Repository interface {
	RunTransaction(f func(Transaction) error) error

	User() UserRepositoryQuery
	Group() GroupRepositoryQuery
}

type Transaction interface {
	User() UserRepositoryCommand
	Group() GroupRepositoryCommand
}
