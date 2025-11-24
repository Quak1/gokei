package store

import (
	"context"
	"database/sql"
)

type MockQuerierTx struct {
	CreateAccountFunc              func(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateCategoryFunc             func(ctx context.Context, arg CreateCategoryParams) (Category, error)
	CreateTokenFunc                func(ctx context.Context, arg CreateTokenParams) (Token, error)
	CreateTransactionFunc          func(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	CreateUserFunc                 func(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAccountByIdFunc          func(ctx context.Context, arg DeleteAccountByIdParams) (sql.Result, error)
	DeleteCategoryByIdFunc         func(ctx context.Context, arg DeleteCategoryByIdParams) (sql.Result, error)
	DeleteTransactionByIDFunc      func(ctx context.Context, arg DeleteTransactionByIDParams) (sql.Result, error)
	DeleteUserByIdFunc             func(ctx context.Context, id int32) (sql.Result, error)
	GetAccountByIDFunc             func(ctx context.Context, arg GetAccountByIDParams) (Account, error)
	GetAccountSumBalanceFunc       func(ctx context.Context, arg GetAccountSumBalanceParams) (GetAccountSumBalanceRow, error)
	GetAllAccountsFunc             func(ctx context.Context) ([]Account, error)
	GetAllCategoriesFunc           func(ctx context.Context, arg GetAllCategoriesParams) ([]Category, error)
	GetAllTransactionsFunc         func(ctx context.Context, userID int32) ([]GetAllTransactionsRow, error)
	GetAllUsersFunc                func(ctx context.Context) ([]User, error)
	GetCategoryByIDFunc            func(ctx context.Context, arg GetCategoryByIDParams) (Category, error)
	GetCategoryByNameFunc          func(ctx context.Context, arg GetCategoryByNameParams) (Category, error)
	GetTransactionByIDFunc         func(ctx context.Context, arg GetTransactionByIDParams) (GetTransactionByIDRow, error)
	GetTransactionsByAccountIDFunc func(ctx context.Context, arg GetTransactionsByAccountIDParams) ([]GetTransactionsByAccountIDRow, error)
	GetUserAccountsFunc            func(ctx context.Context, userID int32) ([]Account, error)
	GetUserByIDFunc                func(ctx context.Context, id int32) (User, error)
	GetUserByUsernameFunc          func(ctx context.Context, username string) (User, error)
	GetUserFromTokenFunc           func(ctx context.Context, arg GetUserFromTokenParams) (GetUserFromTokenRow, error)
	UpdateAccountByIdFunc          func(ctx context.Context, arg UpdateAccountByIdParams) (sql.Result, error)
	UpdateBalanceFunc              func(ctx context.Context, arg UpdateBalanceParams) (int64, error)
	UpdateCategoryByIdFunc         func(ctx context.Context, arg UpdateCategoryByIdParams) (sql.Result, error)
	UpdateTransactionByIdFunc      func(ctx context.Context, arg UpdateTransactionByIdParams) (sql.Result, error)
	UpdateUserByIdFunc             func(ctx context.Context, arg UpdateUserByIdParams) (sql.Result, error)

	WithTxFunc func(tx *sql.Tx) QuerierTx
}

var _ QuerierTx = (*MockQuerierTx)(nil)

// Account queries
func (m *MockQuerierTx) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	if m.CreateAccountFunc != nil {
		return m.CreateAccountFunc(ctx, arg)
	}
	return Account{}, nil
}

func (m *MockQuerierTx) GetAllAccounts(ctx context.Context) ([]Account, error) {
	if m.GetAllAccountsFunc != nil {
		return m.GetAllAccountsFunc(ctx)
	}
	return []Account{}, nil
}

func (m *MockQuerierTx) GetUserAccounts(ctx context.Context, userID int32) ([]Account, error) {
	if m.GetUserAccountsFunc != nil {
		return m.GetUserAccountsFunc(ctx, userID)
	}
	return []Account{}, nil
}

func (m *MockQuerierTx) UpdateBalance(ctx context.Context, arg UpdateBalanceParams) (int64, error) {
	if m.UpdateBalanceFunc != nil {
		return m.UpdateBalanceFunc(ctx, arg)
	}
	return 0, nil
}

func (m *MockQuerierTx) GetAccountByID(ctx context.Context, arg GetAccountByIDParams) (Account, error) {
	if m.GetAccountByIDFunc != nil {
		return m.GetAccountByIDFunc(ctx, arg)
	}
	return Account{}, nil
}

func (m *MockQuerierTx) DeleteAccountById(ctx context.Context, arg DeleteAccountByIdParams) (sql.Result, error) {
	if m.DeleteAccountByIdFunc != nil {
		return m.DeleteAccountByIdFunc(ctx, arg)
	}
	return NewMockResult(1), nil
}

func (m *MockQuerierTx) GetAccountSumBalance(ctx context.Context, arg GetAccountSumBalanceParams) (GetAccountSumBalanceRow, error) {
	if m.GetAccountSumBalanceFunc != nil {
		return m.GetAccountSumBalanceFunc(ctx, arg)
	}
	return GetAccountSumBalanceRow{}, nil
}

func (m *MockQuerierTx) UpdateAccountById(ctx context.Context, arg UpdateAccountByIdParams) (sql.Result, error) {
	if m.UpdateAccountByIdFunc != nil {
		return m.UpdateAccountByIdFunc(ctx, arg)
	}
	return NewMockResult(1), nil
}

// Category queries
func (m *MockQuerierTx) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	if m.CreateCategoryFunc != nil {
		return m.CreateCategoryFunc(ctx, arg)
	}
	return Category{}, nil
}

func (m *MockQuerierTx) GetAllCategories(ctx context.Context, arg GetAllCategoriesParams) ([]Category, error) {
	if m.GetAllCategoriesFunc != nil {
		return m.GetAllCategoriesFunc(ctx, arg)
	}
	return []Category{}, nil
}

func (m *MockQuerierTx) GetCategoryByID(ctx context.Context, arg GetCategoryByIDParams) (Category, error) {
	if m.GetCategoryByIDFunc != nil {
		return m.GetCategoryByIDFunc(ctx, arg)
	}
	return Category{}, nil
}

func (m *MockQuerierTx) GetCategoryByName(ctx context.Context, arg GetCategoryByNameParams) (Category, error) {
	if m.GetCategoryByNameFunc != nil {
		return m.GetCategoryByName(ctx, arg)
	}
	return Category{}, nil
}

func (m *MockQuerierTx) DeleteCategoryById(ctx context.Context, arg DeleteCategoryByIdParams) (sql.Result, error) {
	if m.DeleteCategoryByIdFunc != nil {
		return m.DeleteCategoryByIdFunc(ctx, arg)
	}
	return NewMockResult(1), nil
}

func (m *MockQuerierTx) UpdateCategoryById(ctx context.Context, arg UpdateCategoryByIdParams) (sql.Result, error) {
	if m.UpdateCategoryByIdFunc != nil {
		return m.UpdateCategoryByIdFunc(ctx, arg)
	}
	return NewMockResult(1), nil
}

// Token queries
func (m *MockQuerierTx) CreateToken(ctx context.Context, arg CreateTokenParams) (Token, error) {
	if m.CreateTokenFunc != nil {
		return m.CreateTokenFunc(ctx, arg)
	}
	return Token{}, nil
}

// Transaction queries
func (m *MockQuerierTx) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	if m.CreateTransactionFunc != nil {
		return m.CreateTransactionFunc(ctx, arg)
	}
	return Transaction{}, nil
}

func (m *MockQuerierTx) GetAllTransactions(ctx context.Context, userID int32) ([]GetAllTransactionsRow, error) {
	if m.GetAllTransactionsFunc != nil {
		return m.GetAllTransactionsFunc(ctx, userID)
	}
	return []GetAllTransactionsRow{}, nil
}

func (m *MockQuerierTx) GetTransactionsByAccountID(ctx context.Context, arg GetTransactionsByAccountIDParams) ([]GetTransactionsByAccountIDRow, error) {
	if m.GetTransactionsByAccountIDFunc != nil {
		return m.GetTransactionsByAccountIDFunc(ctx, arg)
	}
	return []GetTransactionsByAccountIDRow{}, nil
}

func (m *MockQuerierTx) GetTransactionByID(ctx context.Context, arg GetTransactionByIDParams) (GetTransactionByIDRow, error) {
	if m.GetTransactionByIDFunc != nil {
		return m.GetTransactionByIDFunc(ctx, arg)
	}
	return GetTransactionByIDRow{}, nil
}

func (m *MockQuerierTx) DeleteTransactionByID(ctx context.Context, arg DeleteTransactionByIDParams) (sql.Result, error) {
	if m.DeleteTransactionByIDFunc != nil {
		return m.DeleteTransactionByIDFunc(ctx, arg)
	}
	return NewMockResult(1), nil
}

func (m *MockQuerierTx) UpdateTransactionById(ctx context.Context, arg UpdateTransactionByIdParams) (sql.Result, error) {
	if m.UpdateTransactionByIdFunc != nil {
		return m.UpdateTransactionByIdFunc(ctx, arg)
	}
	return NewMockResult(1), nil
}

// User queries
func (m *MockQuerierTx) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, arg)
	}
	return User{}, nil
}

func (m *MockQuerierTx) GetAllUsers(ctx context.Context) ([]User, error) {
	if m.GetAllUsersFunc != nil {
		return m.GetAllUsersFunc(ctx)
	}
	return []User{}, nil
}

func (m *MockQuerierTx) GetUserByID(ctx context.Context, id int32) (User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return User{}, nil
}

func (m *MockQuerierTx) GetUserByUsername(ctx context.Context, username string) (User, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(ctx, username)
	}
	return User{}, nil
}

func (m *MockQuerierTx) GetUserFromToken(ctx context.Context, arg GetUserFromTokenParams) (GetUserFromTokenRow, error) {
	if m.GetUserFromTokenFunc != nil {
		return m.GetUserFromTokenFunc(ctx, arg)
	}
	return GetUserFromTokenRow{}, nil
}

func (m *MockQuerierTx) DeleteUserById(ctx context.Context, id int32) (sql.Result, error) {
	if m.DeleteUserByIdFunc != nil {
		return m.DeleteUserByIdFunc(ctx, id)
	}
	return NewMockResult(1), nil
}

func (m *MockQuerierTx) UpdateUserById(ctx context.Context, arg UpdateUserByIdParams) (sql.Result, error) {
	if m.UpdateUserByIdFunc != nil {
		return m.UpdateUserByIdFunc(ctx, arg)
	}
	return NewMockResult(1), nil
}

// Tx
func (m *MockQuerierTx) WithTx(tx *sql.Tx) QuerierTx {
	if m.WithTxFunc != nil {
		return m.WithTxFunc(tx)
	}
	return nil
}

// implements sql.Result
type MockResult struct {
	lastInsertID int64
	rowsAffected int64
}

func (m MockResult) LastInsertId() (int64, error) {
	return m.lastInsertID, nil
}

func (m MockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func NewMockResult(rowsAffected int64) sql.Result {
	return MockResult{rowsAffected: rowsAffected}
}
