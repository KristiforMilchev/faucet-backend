package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"faucet/api/interfaces"
	"faucet/api/models"
)

type Accounts struct {
	ConnectionString string
	Storage          interfaces.Storage
}

func (accountService *Accounts) OpenConnection(storage *interfaces.Storage) bool {

	accountService.Storage = *storage
	accountService.Storage.Open(accountService.ConnectionString)
	return true

}

func (accountService *Accounts) UserExists(address string) (models.Account, bool) {
	var account models.Account
	data := accountService.Storage.Single("select id, address, last_lease from public.accounts where address = $1", []interface{}{&address})

	err := data.Scan(&account.Id, &account.Address, &account.LastLease)
	if err != nil {
		fmt.Printf("Failed to fetch account with email %v", address)
		fmt.Println(err)
		return account, false
	}

	fmt.Println(&account)
	return account, true
}

func (accountService *Accounts) CreateUser(address string) (uuid.UUID, error) {
	sql := `
		INSERT INTO public.accounts(
		address, last_lease, total_lease)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	queryResult := accountService.Storage.Single(sql, []interface{}{
		&address,
		time.Now().Add((time.Hour * 200) * -1),
		0,
	})

	var id uuid.UUID
	err := queryResult.Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (accountService *Accounts) UpdateDrip(address string) bool {
	sql := `
		UPDATE public.accounts
		SET last_lease=$1
		WHERE address = $2;
	`

	_, err := accountService.Storage.Exec(sql, []interface{}{
		time.Now(),
		address,
	})

	return err == nil
}

func (accountService *Accounts) Get() []models.Account {

	rows := accountService.Storage.Where("SELECT * from public.accounts", []interface{}{})

	var accounts []models.Account

	for rows.Next() {
		var account models.Account
		rows.Scan(&account.Id, &account.Address, &account.Name, &account.LastLease, &account.TotalLease)
		accounts = append(accounts, account)
		fmt.Println(account)
	}

	return accounts
}

func (accountService *Accounts) Close() bool {
	accountService.Storage.Close()
	return true
}
