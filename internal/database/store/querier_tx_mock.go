package store

import (
	"context"
	"database/sql"
)

type MockQuerierTx struct {
	AutoUpdateBalanceFunc              func(ctx context.Context, arg AutoUpdateBalanceParams) (int64, error)
	CreateAccountFunc                  func(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateCategoryFunc                 func(ctx context.Context, arg CreateCategoryParams) (Category, error)
	CreateOccurrenceFunc               func(ctx context.Context, arg CreateOccurrenceParams) (RecurringTransactionOccurrence, error)
	CreateRecurringTransactionFunc     func(ctx context.Context, arg CreateRecurringTransactionParams) (RecurringTransaction, error)
	CreateTokenFunc                    func(ctx context.Context, arg CreateTokenParams) (Token, error)
	CreateTransactionFunc              func(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	CreateUserFunc                     func(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAccountByIdFunc              func(ctx context.Context, arg DeleteAccountByIdParams) (sql.Result, error)
	DeleteCategoryByIdFunc             func(ctx context.Context, arg DeleteCategoryByIdParams) (sql.Result, error)
	DeleteRecurringTransactionFunc     func(ctx context.Context, arg DeleteRecurringTransactionParams) (sql.Result, error)
	DeleteTransactionByIDFunc          func(ctx context.Context, arg DeleteTransactionByIDParams) (sql.Result, error)
	DeleteUserByIdFunc                 func(ctx context.Context, id int32) (sql.Result, error)
	GetAccountByIDFunc                 func(ctx context.Context, arg GetAccountByIDParams) (Account, error)
	GetAccountSumBalanceFunc           func(ctx context.Context, arg GetAccountSumBalanceParams) (GetAccountSumBalanceRow, error)
	GetActiveRecurringTransactionsFunc func(ctx context.Context, arg GetActiveRecurringTransactionsParams) ([]RecurringTransaction, error)
	GetAllAccountsFunc                 func(ctx context.Context) ([]Account, error)
	GetAllCategoriesFunc               func(ctx context.Context, arg GetAllCategoriesParams) ([]Category, error)
	GetAllTransactionsFunc             func(ctx context.Context, userID int32) ([]GetAllTransactionsRow, error)
	GetAllUsersFunc                    func(ctx context.Context) ([]User, error)
	GetCategoryByIDFunc                func(ctx context.Context, arg GetCategoryByIDParams) (Category, error)
	GetCategoryByNameFunc              func(ctx context.Context, arg GetCategoryByNameParams) (Category, error)
	GetLastOccurrenceFunc              func(ctx context.Context, recurringTransactionID int32) (RecurringTransactionOccurrence, error)
	GetOccurrenceForDateFunc           func(ctx context.Context, arg GetOccurrenceForDateParams) (RecurringTransactionOccurrence, error)
	GetOccurrencesFunc                 func(ctx context.Context, recurringTransactionID int32) ([]RecurringTransactionOccurrence, error)
	GetRecurringTransactionByIDFunc    func(ctx context.Context, arg GetRecurringTransactionByIDParams) (RecurringTransaction, error)
	GetTransactionByIDFunc             func(ctx context.Context, arg GetTransactionByIDParams) (GetTransactionByIDRow, error)
	GetTransactionsByAccountIDFunc     func(ctx context.Context, arg GetTransactionsByAccountIDParams) ([]GetTransactionsByAccountIDRow, error)
	GetUserAccountsFunc                func(ctx context.Context, userID int32) ([]Account, error)
	GetUserByIDFunc                    func(ctx context.Context, id int32) (User, error)
	GetUserByUsernameFunc              func(ctx context.Context, username string) (User, error)
	GetUserFromTokenFunc               func(ctx context.Context, arg GetUserFromTokenParams) (GetUserFromTokenRow, error)
	GetUserRecurringTransactionsFunc   func(ctx context.Context, userID int32) ([]RecurringTransaction, error)
	UpdateAccountByIdFunc              func(ctx context.Context, arg UpdateAccountByIdParams) (sql.Result, error)
	UpdateBalanceFunc                  func(ctx context.Context, arg UpdateBalanceParams) (int64, error)
	UpdateCategoryByIdFunc             func(ctx context.Context, arg UpdateCategoryByIdParams) (sql.Result, error)
	UpdateRecurringTransactionFunc     func(ctx context.Context, arg UpdateRecurringTransactionParams) (sql.Result, error)
	UpdateTransactionByIdFunc          func(ctx context.Context, arg UpdateTransactionByIdParams) (sql.Result, error)
	UpdateUserByIdFunc                 func(ctx context.Context, arg UpdateUserByIdParams) (sql.Result, error)

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
func (m *MockQuerierTx) AutoUpdateBalance(ctx context.Context, arg AutoUpdateBalanceParams) (int64, error) {
	if m.AutoUpdateBalanceFunc != nil {
		return m.AutoUpdateBalanceFunc(ctx, arg)
	}
	return 0, nil
}

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

// Recurring Transactions
func (m *MockQuerierTx) CreateOccurrence(ctx context.Context, arg CreateOccurrenceParams) (RecurringTransactionOccurrence, error) {
	if m.CreateOccurrenceFunc != nil {
		return m.CreateOccurrenceFunc(ctx, arg)
	}
	return RecurringTransactionOccurrence{}, nil
}

func (m *MockQuerierTx) CreateRecurringTransaction(ctx context.Context, arg CreateRecurringTransactionParams) (RecurringTransaction, error) {
	if m.CreateRecurringTransactionFunc != nil {
		return m.CreateRecurringTransactionFunc(ctx, arg)
	}
	return RecurringTransaction{}, nil
}

func (m *MockQuerierTx) DeleteRecurringTransaction(ctx context.Context, arg DeleteRecurringTransactionParams) (sql.Result, error) {
	if m.DeleteRecurringTransactionFunc != nil {
		return m.DeleteRecurringTransactionFunc(ctx, arg)
	}
	return nil, nil
}

func (m *MockQuerierTx) GetActiveRecurringTransactions(ctx context.Context, arg GetActiveRecurringTransactionsParams) ([]RecurringTransaction, error) {
	if m.GetActiveRecurringTransactionsFunc != nil {
		return m.GetActiveRecurringTransactionsFunc(ctx, arg)
	}
	return []RecurringTransaction{}, nil
}

func (m *MockQuerierTx) GetLastOccurrence(ctx context.Context, recurringTransactionID int32) (RecurringTransactionOccurrence, error) {
	if m.GetLastOccurrenceFunc != nil {
		return m.GetLastOccurrenceFunc(ctx, recurringTransactionID)
	}
	return RecurringTransactionOccurrence{}, nil
}

func (m *MockQuerierTx) GetOccurrenceForDate(ctx context.Context, arg GetOccurrenceForDateParams) (RecurringTransactionOccurrence, error) {
	if m.GetOccurrenceForDateFunc != nil {
		return m.GetOccurrenceForDateFunc(ctx, arg)
	}
	return RecurringTransactionOccurrence{}, nil
}

func (m *MockQuerierTx) GetOccurrences(ctx context.Context, recurringTransactionID int32) ([]RecurringTransactionOccurrence, error) {
	if m.GetOccurrencesFunc != nil {
		return m.GetOccurrencesFunc(ctx, recurringTransactionID)
	}
	return []RecurringTransactionOccurrence{}, nil
}

func (m *MockQuerierTx) GetRecurringTransactionByID(ctx context.Context, arg GetRecurringTransactionByIDParams) (RecurringTransaction, error) {
	if m.GetRecurringTransactionByIDFunc != nil {
		return m.GetRecurringTransactionByIDFunc(ctx, arg)
	}
	return RecurringTransaction{}, nil
}

func (m *MockQuerierTx) GetUserRecurringTransactions(ctx context.Context, userID int32) ([]RecurringTransaction, error) {
	if m.GetUserRecurringTransactionsFunc != nil {
		return m.GetUserRecurringTransactionsFunc(ctx, userID)
	}
	return []RecurringTransaction{}, nil
}

func (m *MockQuerierTx) UpdateRecurringTransaction(ctx context.Context, arg UpdateRecurringTransactionParams) (sql.Result, error) {
	if m.UpdateRecurringTransactionFunc != nil {
		return m.UpdateRecurringTransactionFunc(ctx, arg)
	}
	return nil, nil
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
