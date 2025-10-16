package expense

import (
	"context"
	"database/sql"
	"time"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/shopspring/decimal"
)

type ExpenseService struct {
	expenseRepo   *expenseRepository
	recurringRepo *recurringExpenseRepository
	storeService  *store.StoreService
	atomicProcess *db.AtomicProcess
}

func NewExpenseService(expenseRepo *expenseRepository, recurringRepo *recurringExpenseRepository, storeService *store.StoreService, atomicProcess *db.AtomicProcess) *ExpenseService {
	return &ExpenseService{
		expenseRepo:   expenseRepo,
		recurringRepo: recurringRepo,
		storeService:  storeService,
		atomicProcess: atomicProcess,
	}
}

func (s *ExpenseService) GetExpenseByID(ctx context.Context, storeID, id string) (*Expense, error) {
	exp, err := s.expenseRepo.findOne(ctx, s.expenseRepo.scopeStoreID(storeID), s.expenseRepo.scopeID(id))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func (s *ExpenseService) ListExpenses(ctx context.Context, storeID string, page, pageSize int, orderBy string, ascending bool) ([]*Expense, error) {
	exps, err := s.expenseRepo.list(ctx,
		s.expenseRepo.scopeStoreID(storeID),
		db.WithPagination(page, pageSize),
		db.WithSorting(orderBy, ascending),
	)
	if err != nil {
		return nil, err
	}
	return exps, nil
}

func (s *ExpenseService) CountExpenses(ctx context.Context, storeID string) (int64, error) {
	count, err := s.expenseRepo.count(ctx,
		s.expenseRepo.scopeStoreID(storeID),
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
	st, err := s.storeService.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	expense := &Expense{
		StoreID:  storeID,
		Amount:   req.Amount,
		Currency: st.Currency,
		Category: req.Category,
		Type:     req.Type,
		Note:     req.Note,
	}
	if req.OccurredOn != nil && !req.OccurredOn.IsZero() {
		expense.OccurredOn = startOfDay((*req.OccurredOn).In(time.UTC))
	} else {
		expense.OccurredOn = startOfDay(time.Now().In(time.UTC))
	}
	if err := s.expenseRepo.createOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, storeID, expenseID string, updates *UpdateExpenseRequest) (*Expense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.expenseRepo.findOne(ctx, s.expenseRepo.scopeStoreID(storeID), s.expenseRepo.scopeID(expenseID))
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
	if err := s.expenseRepo.updateOne(ctx, expense); err != nil {
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
	if err := s.expenseRepo.deleteOne(ctx, s.expenseRepo.scopeStoreID(storeID), s.expenseRepo.scopeID(id)); err != nil {
		return err
	}
	return nil
}

func (s *ExpenseService) GetRecurringExpenseByID(ctx context.Context, storeID, id string) (*RecurringExpense, error) {
	exp, err := s.recurringRepo.findOne(ctx, s.recurringRepo.scopeStoreID(storeID), s.recurringRepo.scopeID(id))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func (s *ExpenseService) ListRecurringExpenses(ctx context.Context, storeID string, page, pageSize int, orderBy string, ascending bool) ([]*RecurringExpense, error) {
	exps, err := s.recurringRepo.list(ctx,
		s.recurringRepo.scopeStoreID(storeID),
		db.WithPagination(page, pageSize),
		db.WithSorting(orderBy, ascending),
	)
	if err != nil {
		return nil, err
	}
	return exps, nil
}

func (s *ExpenseService) CountRecurringExpenses(ctx context.Context, storeID string) (int64, error) {
	count, err := s.recurringRepo.count(ctx,
		s.recurringRepo.scopeStoreID(storeID),
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
	st, err := s.storeService.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	expense := &RecurringExpense{
		StoreID:            storeID,
		Frequency:          req.Frequency,
		RecurringEndDate:   req.RecurringEndDate,
		RecurringStartDate: req.RecurringStartDate,
		Amount:             req.Amount,
		Currency:           st.Currency,
		Category:           req.Category,
		Note:               req.Note,
	}
	today := startOfDay(time.Now().In(time.UTC))
	start := startOfDay(req.RecurringStartDate.In(time.UTC))
	expense.NextRecurringDate = s.initNextRecurringDateForCreate(start, req.Frequency, today)
	if err := s.atomicProcess.Exec(ctx, func(txCtx context.Context) error {
		if err := s.recurringRepo.createOne(txCtx, expense); err != nil {
			return err
		}
		return s.backfillPastOccurrencesForCreate(txCtx, st, expense, start, today)
	}); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) UpdateRecurringExpense(ctx context.Context, storeID, expenseID string, updates *UpdateRecurringExpenseRequest) (*RecurringExpense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.recurringRepo.findOne(ctx, s.recurringRepo.scopeStoreID(storeID), s.recurringRepo.scopeID(expenseID))
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
	if err := s.recurringRepo.updateOne(ctx, expense); err != nil {
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
	if err := s.recurringRepo.deleteOne(ctx, s.recurringRepo.scopeStoreID(storeID), s.recurringRepo.scopeID(id)); err != nil {
		return err
	}
	return nil
}

func (s *ExpenseService) GetRecurringExpenseHistory(ctx context.Context, storeID, recurringExpenseID string) ([]*Expense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	history, err := s.expenseRepo.list(ctx,
		s.expenseRepo.scopeStoreID(storeID),
		s.expenseRepo.scopeRecurringExpenseID(recurringExpenseID),
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
	expense, err := s.recurringRepo.findOne(ctx, s.recurringRepo.scopeStoreID(storeID), s.recurringRepo.scopeID(recurringExpenseID))
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
	if err := s.recurringRepo.updateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) ResumeRecurringExpense(ctx context.Context, storeID, recurringExpenseID string) (*RecurringExpense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.recurringRepo.findOne(ctx, s.recurringRepo.scopeStoreID(storeID), s.recurringRepo.scopeID(recurringExpenseID))
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
	if err := s.recurringRepo.updateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) CancelRecurringExpense(ctx context.Context, storeID, recurringExpenseID string) (*RecurringExpense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.recurringRepo.findOne(ctx, s.recurringRepo.scopeStoreID(storeID), s.recurringRepo.scopeID(recurringExpenseID))
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
	if err := s.recurringRepo.updateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) EndRecurringExpense(ctx context.Context, storeID, recurringExpenseID string) (*RecurringExpense, error) {
	if err := s.storeService.ValidateStoreID(ctx, storeID); err != nil {
		return nil, err
	}
	expense, err := s.recurringRepo.findOne(ctx, s.recurringRepo.scopeStoreID(storeID), s.recurringRepo.scopeID(recurringExpenseID))
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
	if err := s.recurringRepo.updateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

// backfillPastOccurrences creates Expense rows for each past occurrence from RecurringStartDate
// up to but not including the first future occurrence relative to now, honoring RecurringEndDate.
// backfillPastOccurrences removed due to new NextRecurringDate approach

// ProcessRecurringExpensesDaily checks all active recurring expenses and creates an expense
// when an occurrence is due today (including first-time when start date passes).
func (s *ExpenseService) ProcessRecurringExpensesDaily(ctx context.Context) error {
	today := startOfDay(time.Now().In(time.UTC))
	end := today.Add(24 * time.Hour)
	return s.atomicProcess.Exec(ctx, func(txCtx context.Context) error {
		recurs, err := s.fetchActiveRecurringForWindow(txCtx, today)
		if err != nil || len(recurs) == 0 {
			return err
		}
		storeCache := map[string]*store.Store{}
		toCreate := make([]*Expense, 0)
		for _, re := range recurs {
			created, err := s.processOneRecurringForDay(txCtx, storeCache, re, today, end)
			if err != nil {
				return err
			}
			if created != nil {
				toCreate = append(toCreate, created)
			}
		}
		if len(toCreate) == 0 {
			return nil
		}
		return s.expenseRepo.createMany(txCtx, toCreate)
	})
}

func (s *ExpenseService) fetchActiveRecurringForWindow(ctx context.Context, today time.Time) ([]*RecurringExpense, error) {
	return s.recurringRepo.list(ctx,
		s.recurringRepo.scopeStatus(RecurringExpenseStatusActive),
		s.recurringRepo.scopeStartDateLte(today),
		s.recurringRepo.scopeEndDateGteOrNull(today),
	)
}

func (s *ExpenseService) processOneRecurringForDay(ctx context.Context, storeCache map[string]*store.Store, re *RecurringExpense, today, end time.Time) (*Expense, error) {
	if err := s.ensureNextRecurringInitialized(ctx, re, today); err != nil {
		return nil, err
	}
	if !re.NextRecurringDate.Before(end) {
		return nil, nil
	}
	exp, err := s.buildExpenseIfNotExists(ctx, storeCache, re, re.NextRecurringDate)
	if err != nil {
		return nil, err
	}
	if err := s.advanceOrEndRecurring(ctx, re); err != nil {
		return nil, err
	}
	return exp, nil
}

func (s *ExpenseService) ensureNextRecurringInitialized(ctx context.Context, re *RecurringExpense, today time.Time) error {
	if !re.NextRecurringDate.IsZero() {
		return nil
	}
	start := startOfDay(re.RecurringStartDate.In(time.UTC))
	if !start.After(today) {
		re.NextRecurringDate = nextOccurrenceOnOrAfter(start, re.Frequency, today)
	} else {
		re.NextRecurringDate = start
	}
	return s.recurringRepo.updateOne(ctx, re)
}

func (s *ExpenseService) buildExpenseIfNotExists(ctx context.Context, storeCache map[string]*store.Store, re *RecurringExpense, forDate time.Time) (*Expense, error) {
	exists, err := s.hasExpenseForRecurringOnDate(ctx, re, forDate)
	if err != nil || exists {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	st, ok := storeCache[re.StoreID]
	if !ok {
		st, err = s.storeService.GetStoreByID(ctx, re.StoreID)
		if err != nil {
			return nil, err
		}
		storeCache[re.StoreID] = st
	}
	return &Expense{
		StoreID:            re.StoreID,
		RecurringExpenseID: sql.NullString{String: re.ID, Valid: true},
		Amount:             re.Amount,
		Currency:           st.Currency,
		Category:           re.Category,
		Type:               ExpenseTypeRecurring,
		Note:               re.Note,
		OccurredOn:         startOfDay(forDate),
	}, nil
}

func (s *ExpenseService) advanceOrEndRecurring(ctx context.Context, re *RecurringExpense) error {
	next := addFrequency(startOfDay(re.NextRecurringDate), re.Frequency)
	if re.RecurringEndDate.Valid && next.After(re.RecurringEndDate.Time) {
		sm := NewRecurringExpenseStateMachine(re)
		if sm.CanTransitionTo(RecurringExpenseStatusEnded) {
			_ = sm.TransitionTo(RecurringExpenseStatusEnded)
		}
	} else {
		re.NextRecurringDate = next
	}
	return s.recurringRepo.updateOne(ctx, re)
}

// fetchActiveRecurringForDate and buildPendingExpensesForDate removed; replaced by NextRecurringDate-based processing

// Helpers for frequency calculations
func addFrequency(t time.Time, freq RecurringExpenseFrequency) time.Time {
	switch freq.String {
	case "daily":
		return t.AddDate(0, 0, 1)
	case "weekly":
		return t.AddDate(0, 0, 7)
	case "monthly":
		return t.AddDate(0, 1, 0)
	case "yearly":
		return t.AddDate(1, 0, 0)
	default:
		return t
	}
}

// nextOccurrenceOnOrAfter returns the first occurrence that is on or after 'after'.
func nextOccurrenceOnOrAfter(start time.Time, freq RecurringExpenseFrequency, after time.Time) time.Time {
	d := start
	for d.Before(after) {
		d = addFrequency(d, freq)
	}
	return d
}

// isOccurrenceDueToday returns true if there's an occurrence that should be recorded today.
// We treat the recurrence to fire when the occurrence date's date equals today's date in UTC.
// isOccurrenceOnDate no longer needed with NextRecurringDate

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func (s *ExpenseService) hasExpenseForRecurringOnDate(ctx context.Context, re *RecurringExpense, day time.Time) (bool, error) {
	day = startOfDay(day)
	count, err := s.expenseRepo.count(ctx,
		s.expenseRepo.scopeStoreID(re.StoreID),
		s.expenseRepo.scopeRecurringExpenseID(re.ID),
		s.expenseRepo.scopeOccurredOn(day),
	)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// initNextRecurringDateForCreate decides initial NextRecurringDate for a newly created recurring expense
func (s *ExpenseService) initNextRecurringDateForCreate(start time.Time, freq RecurringExpenseFrequency, today time.Time) time.Time {
	if start.After(today) {
		return start
	}
	return nextOccurrenceOnOrAfter(start, freq, today)
}

// backfillPastOccurrencesForCreate creates expenses for each due date prior to today during creation
func (s *ExpenseService) backfillPastOccurrencesForCreate(ctx context.Context, st *store.Store, re *RecurringExpense, start, today time.Time) error {
	if !start.Before(today) {
		return nil
	}
	var toCreate []*Expense
	for d := start; d.Before(today); d = addFrequency(d, re.Frequency) {
		if re.RecurringEndDate.Valid && d.After(re.RecurringEndDate.Time) {
			break
		}
		exist, err := s.hasExpenseForRecurringOnDate(ctx, re, d)
		if err != nil {
			return err
		}
		if exist {
			continue
		}
		toCreate = append(toCreate, &Expense{
			StoreID:            re.StoreID,
			RecurringExpenseID: sql.NullString{String: re.ID, Valid: true},
			Amount:             re.Amount,
			Currency:           st.Currency,
			Category:           re.Category,
			Type:               ExpenseTypeRecurring,
			Note:               re.Note,
			OccurredOn:         d,
		})
	}
	if len(toCreate) == 0 {
		return nil
	}
	return s.expenseRepo.createMany(ctx, toCreate)
}

// ---- Analytics wrappers ----

// ExpenseTotals returns total amount and count of expenses in range.
func (s *ExpenseService) ExpenseTotals(ctx context.Context, storeID string, from, to time.Time) (total decimal.Decimal, count int64, err error) {
	total, err = s.expenseRepo.sumAmount(ctx, s.expenseRepo.scopeStoreID(storeID), s.expenseRepo.scopeOccurredOnRange(from, to))
	if err != nil {
		return
	}
	count, err = s.expenseRepo.count(ctx, s.expenseRepo.scopeStoreID(storeID), s.expenseRepo.scopeOccurredOnRange(from, to))
	return
}

func (s *ExpenseService) ExpenseBreakdownByCategory(ctx context.Context, storeID string, from, to time.Time) ([]types.KeyValue, error) {
	return s.expenseRepo.breakdownByCategory(ctx, s.expenseRepo.scopeStoreID(storeID), s.expenseRepo.scopeOccurredOnRange(from, to))
}

func (s *ExpenseService) ExpenseAmountTimeSeries(ctx context.Context, storeID string, from, to time.Time, bucket string) ([]types.TimeSeriesRow, error) {
	return s.expenseRepo.amountTimeSeries(ctx, bucket, from, to, s.expenseRepo.scopeStoreID(storeID), s.expenseRepo.scopeOccurredOnRange(from, to))
}

func (s *ExpenseService) MarketingExpensesInRange(ctx context.Context, storeID string, from, to time.Time) (decimal.Decimal, error) {
	return s.expenseRepo.sumAmount(ctx, s.expenseRepo.scopeCategory(ExpenseCategoryMarketing), s.expenseRepo.scopeStoreID(storeID), s.expenseRepo.scopeOccurredOnRange(from, to))
}
