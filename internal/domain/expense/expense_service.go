package expense

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
)

type ExpenseService struct {
	expenseRepo   *ExpenseRepository
	recurringRepo *RecurringExpenseRepository
	storeService  *store.StoreService
	atomicProcess *db.AtomicProcess
}

func NewExpenseService(expenseRepo *ExpenseRepository, recurringRepo *RecurringExpenseRepository, storeService *store.StoreService, atomicProcess *db.AtomicProcess) *ExpenseService {
	return &ExpenseService{
		expenseRepo:   expenseRepo,
		recurringRepo: recurringRepo,
		storeService:  storeService,
		atomicProcess: atomicProcess,
	}
}

func (s *ExpenseService) GetExpenseByID(ctx context.Context, storeID, id string) (*Expense, error) {
	exp, err := s.expenseRepo.FindOne(ctx, s.expenseRepo.ScopeStoreID(storeID), s.expenseRepo.ScopeID(id))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func (s *ExpenseService) ListExpenses(ctx context.Context, storeID string, filter *ExpenseFilter, page, pageSize int, orderBy string, ascending bool) ([]*Expense, error) {
	exps, err := s.expenseRepo.List(ctx,
		s.expenseRepo.ScopeStoreID(storeID),
		s.expenseRepo.ScopeFilter(filter),
		db.WithPagination(page, pageSize),
		db.WithSorting(orderBy, ascending),
	)
	if err != nil {
		return nil, err
	}
	return exps, nil
}

func (s *ExpenseService) CountExpenses(ctx context.Context, storeID string, filter *ExpenseFilter) (int64, error) {
	count, err := s.expenseRepo.Count(ctx,
		s.expenseRepo.ScopeStoreID(storeID),
		s.expenseRepo.ScopeFilter(filter),
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *ExpenseService) CreateExpense(ctx context.Context, storeID string, req *CreateExpenseRequest) (*Expense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense := &Expense{
		StoreID:  storeID,
		Amount:   req.Amount,
		Category: req.Category,
		Type:     req.Type,
		Note:     req.Note,
	}
	if err := s.expenseRepo.CreateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, storeID, expenseID string, updates *UpdateExpenseRequest) (*Expense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.expenseRepo.FindOne(ctx, s.expenseRepo.ScopeStoreID(storeID), s.expenseRepo.ScopeID(expenseID))
	if err != nil {
		return nil, err
	}
	if updates.Amount.Sign() != 0 {
		expense.Amount = updates.Amount
	}
	if updates.Category != "" {
		expense.Category = updates.Category
	}
	if updates.Type != "" {
		expense.Type = updates.Type
	}
	if updates.RecurringExpenseID.Valid {
		expense.RecurringExpenseID = updates.RecurringExpenseID
	}
	if updates.Note != "" {
		expense.Note = updates.Note
	}
	if err := s.expenseRepo.UpdateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, storeID, id string) error {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return err
	}
	if _, err := s.GetExpenseByID(ctx, storeID, id); err != nil {
		return err
	}
	if err := s.expenseRepo.DeleteOne(ctx, s.expenseRepo.ScopeStoreID(storeID), s.expenseRepo.ScopeID(id)); err != nil {
		return err
	}
	return nil
}

func (s *ExpenseService) DeleteExpenses(ctx context.Context, storeID string, filter *ExpenseFilter) error {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return err
	}
	if err := s.expenseRepo.DeleteMany(ctx, s.expenseRepo.ScopeStoreID(storeID), s.expenseRepo.ScopeFilter(filter)); err != nil {
		return err
	}
	return nil
}

func (s *ExpenseService) GetRecurringExpenseByID(ctx context.Context, storeID, id string) (*RecurringExpense, error) {
	exp, err := s.recurringRepo.FindOne(ctx, s.recurringRepo.ScopeStoreID(storeID), s.recurringRepo.ScopeID(id))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func (s *ExpenseService) ListRecurringExpenses(ctx context.Context, storeID string, filter *RecurringExpenseFilter, page, pageSize int, orderBy string, ascending bool) ([]*RecurringExpense, error) {
	exps, err := s.recurringRepo.List(ctx,
		s.recurringRepo.ScopeStoreID(storeID),
		s.recurringRepo.ScopeFilter(filter),
		db.WithPagination(page, pageSize),
		db.WithSorting(orderBy, ascending),
	)
	if err != nil {
		return nil, err
	}
	return exps, nil
}

func (s *ExpenseService) CountRecurringExpenses(ctx context.Context, storeID string, filter *RecurringExpenseFilter) (int64, error) {
	count, err := s.recurringRepo.Count(ctx,
		s.recurringRepo.ScopeStoreID(storeID),
		s.recurringRepo.ScopeFilter(filter),
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *ExpenseService) CreateRecurringExpense(ctx context.Context, storeID string, req *CreateRecurringExpenseRequest) (*RecurringExpense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense := &RecurringExpense{
		StoreID:            storeID,
		Frequency:          req.Frequency,
		RecurringEndDate:   req.RecurringEndDate,
		RecurringStartDate: req.RecurringStartDate,
		Amount:             req.Amount,
		Category:           req.Category,
		Note:               req.Note,
	}
	if err := s.recurringRepo.CreateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) UpdateRecurringExpense(ctx context.Context, storeID, expenseID string, updates *UpdateRecurringExpenseRequest) (*RecurringExpense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.recurringRepo.FindOne(ctx, s.recurringRepo.ScopeStoreID(storeID), s.recurringRepo.ScopeID(expenseID))
	if err != nil {
		return nil, err
	}
	if updates.Frequency.Valid {
		expense.Frequency = updates.Frequency
	}
	if updates.RecurringEndDate.Valid {
		expense.RecurringEndDate = updates.RecurringEndDate
	}
	if !updates.RecurringStartDate.IsZero() {
		expense.RecurringStartDate = updates.RecurringStartDate
	}
	if updates.Amount.Sign() != 0 {
		expense.Amount = updates.Amount
	}
	if updates.Category != "" {
		expense.Category = updates.Category
	}
	if updates.Note != "" {
		expense.Note = updates.Note
	}
	if err := s.recurringRepo.UpdateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) DeleteRecurringExpense(ctx context.Context, storeID, id string) error {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return err
	}
	if _, err := s.GetRecurringExpenseByID(ctx, storeID, id); err != nil {
		return err
	}
	if err := s.recurringRepo.DeleteOne(ctx, s.recurringRepo.ScopeStoreID(storeID), s.recurringRepo.ScopeID(id)); err != nil {
		return err
	}
	return nil
}

func (s *ExpenseService) GetRecurringExpenseHistory(ctx context.Context, storeID, recurringExpenseID string) ([]*Expense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	history, err := s.expenseRepo.List(ctx,
		s.expenseRepo.ScopeStoreID(storeID),
		s.expenseRepo.ScopeRecurringExpenseID(recurringExpenseID),
	)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (s *ExpenseService) PauseRecurringExpense(ctx context.Context, storeID, recurringExpenseID string) (*RecurringExpense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.recurringRepo.FindOne(ctx, s.recurringRepo.ScopeStoreID(storeID), s.recurringRepo.ScopeID(recurringExpenseID))
	if err != nil {
		return nil, err
	}
	sm := NewRecurringExpenseStateMachine(expense)
	if !sm.CanTransitionTo(RecurringExpenseStatusPaused) {
		return nil, nil
	}
	if err := sm.TransitionTo(RecurringExpenseStatusPaused); err != nil {
		return nil, err
	}
	if err := s.recurringRepo.UpdateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) ResumeRecurringExpense(ctx context.Context, storeID, recurringExpenseID string) (*RecurringExpense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.recurringRepo.FindOne(ctx, s.recurringRepo.ScopeStoreID(storeID), s.recurringRepo.ScopeID(recurringExpenseID))
	if err != nil {
		return nil, err
	}
	sm := NewRecurringExpenseStateMachine(expense)
	if !sm.CanTransitionTo(RecurringExpenseStatusActive) {
		return nil, nil
	}
	if err := sm.TransitionTo(RecurringExpenseStatusActive); err != nil {
		return nil, err
	}
	if err := s.recurringRepo.UpdateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) CancelRecurringExpense(ctx context.Context, storeID, recurringExpenseID string) (*RecurringExpense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.recurringRepo.FindOne(ctx, s.recurringRepo.ScopeStoreID(storeID), s.recurringRepo.ScopeID(recurringExpenseID))
	if err != nil {
		return nil, err
	}
	sm := NewRecurringExpenseStateMachine(expense)
	if !sm.CanTransitionTo(RecurringExpenseStatusCanceled) {
		return nil, nil
	}
	if err := sm.TransitionTo(RecurringExpenseStatusCanceled); err != nil {
		return nil, err
	}
	if err := s.recurringRepo.UpdateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) EndRecurringExpense(ctx context.Context, storeID, recurringExpenseID string) (*RecurringExpense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.recurringRepo.FindOne(ctx, s.recurringRepo.ScopeStoreID(storeID), s.recurringRepo.ScopeID(recurringExpenseID))
	if err != nil {
		return nil, err
	}
	sm := NewRecurringExpenseStateMachine(expense)
	if !sm.CanTransitionTo(RecurringExpenseStatusEnded) {
		return nil, nil
	}
	if err := sm.TransitionTo(RecurringExpenseStatusEnded); err != nil {
		return nil, err
	}
	if err := s.recurringRepo.UpdateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}
