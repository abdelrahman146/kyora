package accounting

import (
	"context"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/list"
	"github.com/abdelrahman146/kyora/internal/platform/utils/transformer"
	"github.com/shopspring/decimal"
)

type Service struct {
	storage         *Storage
	atomicProcessor atomic.AtomicProcessor
	bus             *bus.Bus
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus) *Service {
	return &Service{
		storage:         storage,
		atomicProcessor: atomicProcessor,
		bus:             bus,
	}
}

func (s *Service) CreateAsset(ctx context.Context, actor *account.User, biz *business.Business, req *CreateAssetRequest) (*Asset, error) {
	asset := &Asset{
		BusinessID: biz.ID,
		Name:       req.Name,
		Type:       req.Type,
		Value:      req.Value,
		Currency:   biz.Currency,
	}
	if !req.PurchasedAt.IsZero() {
		asset.PurchasedAt = req.PurchasedAt
	}
	if req.Note != "" {
		asset.Note = req.Note
	}
	err := s.storage.asset.CreateOne(ctx, asset)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *Service) UpdateAsset(ctx context.Context, actor *account.User, biz *business.Business, id string, req *UpdateAssetRequest) (*Asset, error) {
	asset, err := s.storage.asset.FindOne(ctx,
		s.storage.asset.ScopeBusinessID(biz.ID),
		s.storage.asset.ScopeID(id),
	)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		asset.Name = req.Name
	}
	if req.Type != "" {
		asset.Type = req.Type
	}
	if !req.Value.IsZero() {
		asset.Value = req.Value
	}
	if !req.PurchasedAt.IsZero() {
		asset.PurchasedAt = req.PurchasedAt
	}
	if req.Note != "" {
		asset.Note = req.Note
	}
	err = s.storage.asset.UpdateOne(ctx, asset)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *Service) DeleteAsset(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	asset, err := s.GetAssetByID(ctx, actor, biz, id)
	if err != nil {
		return err
	}
	return s.storage.asset.DeleteOne(ctx, asset)
}

func (s *Service) GetAssetByID(ctx context.Context, actor *account.User, biz *business.Business, id string) (*Asset, error) {
	return s.storage.asset.FindOne(ctx,
		s.storage.asset.ScopeBusinessID(biz.ID),
		s.storage.asset.ScopeID(id),
	)
}

func (s *Service) ListAssets(ctx context.Context, actor *account.User, biz *business.Business, listReq *list.ListRequest) ([]*Asset, error) {
	return s.storage.asset.FindMany(ctx,
		s.storage.asset.ScopeBusinessID(biz.ID),
		s.storage.asset.WithPagination(listReq.Offset(), listReq.Limit()),
		s.storage.asset.WithOrderBy(listReq.ParsedOrderBy(AssetSchema)),
	)
}

func (s *Service) CountAssets(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.asset.Count(ctx,
		s.storage.asset.ScopeBusinessID(biz.ID),
	)
}

func (s *Service) SumAssetsValue(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.asset.Sum(ctx,
		AssetSchema.Value,
		s.storage.asset.ScopeBusinessID(biz.ID),
		s.storage.asset.ScopeTime(AssetSchema.PurchasedAt, from, to),
	)
}

func (s *Service) CreateInvestment(ctx context.Context, actor *account.User, biz *business.Business, req *CreateInvestmentRequest) (*Investment, error) {
	investment := &Investment{
		BusinessID: biz.ID,
		Amount:     req.Amount,
		Currency:   biz.Currency,
		InvestorID: req.InvestorID,
		Note:       req.Note,
	}
	if !req.InvestedAt.IsZero() {
		investment.InvestedAt = req.InvestedAt
	}
	err := s.storage.investment.CreateOne(ctx, investment)
	if err != nil {
		return nil, err
	}
	return investment, nil
}

func (s *Service) UpdateInvestment(ctx context.Context, actor *account.User, biz *business.Business, id string, req *UpdateInvestmentRequest) (*Investment, error) {
	investment, err := s.storage.investment.FindOne(ctx,
		s.storage.investment.ScopeBusinessID(biz.ID),
		s.storage.investment.ScopeID(id),
	)
	if err != nil {
		return nil, err
	}
	if !req.Amount.IsZero() {
		investment.Amount = req.Amount
	}
	if !req.InvestedAt.IsZero() {
		investment.InvestedAt = req.InvestedAt
	}
	if req.InvestorID != "" {
		investment.InvestorID = req.InvestorID
	}
	if req.Note != "" {
		investment.Note = req.Note
	}
	err = s.storage.investment.UpdateOne(ctx, investment)
	if err != nil {
		return nil, err
	}
	return investment, nil
}

func (s *Service) DeleteInvestment(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	investment, err := s.GetInvestmentByID(ctx, actor, biz, id)
	if err != nil {
		return err
	}
	return s.storage.investment.DeleteOne(ctx, investment)
}

func (s *Service) GetInvestmentByID(ctx context.Context, actor *account.User, biz *business.Business, id string) (*Investment, error) {
	return s.storage.investment.FindOne(ctx,
		s.storage.investment.ScopeBusinessID(biz.ID),
		s.storage.investment.ScopeID(id),
	)
}

func (s *Service) ListInvestments(ctx context.Context, actor *account.User, biz *business.Business, listReq *list.ListRequest) ([]*Investment, error) {
	return s.storage.investment.FindMany(ctx,
		s.storage.investment.ScopeBusinessID(biz.ID),
		s.storage.investment.WithPagination(listReq.Offset(), listReq.Limit()),
		s.storage.investment.WithOrderBy(listReq.ParsedOrderBy(InvestmentSchema)),
	)
}

func (s *Service) CountInvestments(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.investment.Count(ctx,
		s.storage.investment.ScopeBusinessID(biz.ID),
	)
}

func (s *Service) SumInvestmentsAmount(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.investment.Sum(ctx,
		InvestmentSchema.Amount,
		s.storage.investment.ScopeBusinessID(biz.ID),
		s.storage.investment.ScopeTime(InvestmentSchema.InvestedAt, from, to),
	)
}

func (s *Service) SumInvestmentAmountByInvestor(ctx context.Context, actor *account.User, biz *business.Business, investorID string, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.investment.Sum(ctx,
		InvestmentSchema.Amount,
		s.storage.investment.ScopeBusinessID(biz.ID),
		s.storage.investment.ScopeEquals(InvestmentSchema.InvestorID, investorID),
		s.storage.investment.ScopeTime(InvestmentSchema.InvestedAt, from, to),
	)
}

func (s *Service) CreateWithdrawal(ctx context.Context, actor *account.User, biz *business.Business, req *CreateWithdrawalRequest) (*Withdrawal, error) {
	withdrawal := &Withdrawal{
		BusinessID:   biz.ID,
		Amount:       req.Amount,
		Currency:     biz.Currency,
		WithdrawerID: req.WithdrawerID,
	}
	if !req.WithdrawnAt.IsZero() {
		withdrawal.WithdrawnAt = req.WithdrawnAt
	}
	if req.Note != "" {
		withdrawal.Note = req.Note
	}
	err := s.storage.withdrawal.CreateOne(ctx, withdrawal)
	if err != nil {
		return nil, err
	}
	return withdrawal, nil
}

func (s *Service) UpdateWithdrawal(ctx context.Context, actor *account.User, biz *business.Business, id string, req *UpdateWithdrawalRequest) (*Withdrawal, error) {
	withdrawal, err := s.storage.withdrawal.FindOne(ctx,
		s.storage.withdrawal.ScopeBusinessID(biz.ID),
		s.storage.withdrawal.ScopeID(id),
	)
	if err != nil {
		return nil, err
	}
	if !req.Amount.IsZero() {
		withdrawal.Amount = req.Amount
	}
	if !req.WithdrawnAt.IsZero() {
		withdrawal.WithdrawnAt = req.WithdrawnAt
	}
	if req.WithdrawerID != "" {
		withdrawal.WithdrawerID = req.WithdrawerID
	}
	if req.Note != "" {
		withdrawal.Note = req.Note
	}
	err = s.storage.withdrawal.UpdateOne(ctx, withdrawal)
	if err != nil {
		return nil, err
	}
	return withdrawal, nil
}

func (s *Service) DeleteWithdrawal(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	withdrawal, err := s.GetWithdrawalByID(ctx, actor, biz, id)
	if err != nil {
		return err
	}
	return s.storage.withdrawal.DeleteOne(ctx, withdrawal)
}

func (s *Service) GetWithdrawalByID(ctx context.Context, actor *account.User, biz *business.Business, id string) (*Withdrawal, error) {
	return s.storage.withdrawal.FindOne(ctx,
		s.storage.withdrawal.ScopeBusinessID(biz.ID),
		s.storage.withdrawal.ScopeID(id),
	)
}

func (s *Service) ListWithdrawals(ctx context.Context, actor *account.User, biz *business.Business, listReq *list.ListRequest) ([]*Withdrawal, error) {
	return s.storage.withdrawal.FindMany(ctx,
		s.storage.withdrawal.ScopeBusinessID(biz.ID),
		s.storage.withdrawal.WithPagination(listReq.Offset(), listReq.Limit()),
		s.storage.withdrawal.WithOrderBy(listReq.ParsedOrderBy(WithdrawalSchema)),
	)
}

func (s *Service) CountWithdrawals(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.withdrawal.Count(ctx,
		s.storage.withdrawal.ScopeBusinessID(biz.ID),
	)
}

func (s *Service) SumWithdrawalsAmount(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.withdrawal.Sum(ctx,
		WithdrawalSchema.Amount,
		s.storage.withdrawal.ScopeBusinessID(biz.ID),
		s.storage.withdrawal.ScopeTime(WithdrawalSchema.WithdrawnAt, from, to),
	)
}

func (s *Service) SumWithdrawalAmountByWithdrawer(ctx context.Context, actor *account.User, biz *business.Business, withdrawerID string, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.withdrawal.Sum(ctx,
		WithdrawalSchema.Amount,
		s.storage.withdrawal.ScopeBusinessID(biz.ID),
		s.storage.withdrawal.ScopeEquals(WithdrawalSchema.WithdrawerID, withdrawerID),
		s.storage.withdrawal.ScopeTime(WithdrawalSchema.WithdrawnAt, from, to),
	)
}

// ComputeSafeToDrawAmount computes how much money can be withdrawn safely without jeopardizing operations.
//
// The computation is intentionally deterministic:
// safeToDraw = totalIncome - totalCOGS - totalExpenses - totalWithdrawals - businessSafetyBuffer
//
// Important:
// - The calculation respects the provided date range by using expenses/withdrawals within [from,to].
// - SafetyBuffer is treated as an explicit business setting; a value of 0 means no buffer.
func (s *Service) ComputeSafeToDrawAmount(ctx context.Context, actor *account.User, biz *business.Business, totalIncome, totalCOGS decimal.Decimal, from, to time.Time) (decimal.Decimal, error) {
	totalWithdrawals, err := s.SumWithdrawalsAmount(ctx, actor, biz, from, to)
	if err != nil {
		return decimal.Zero, err
	}
	totalExpenses, err := s.SumExpensesAmount(ctx, actor, biz, from, to)
	if err != nil {
		return decimal.Zero, err
	}

	safetyBuffer := biz.SafetyBuffer
	if safetyBuffer.IsZero() {
		referenceEnd := to
		if referenceEnd.IsZero() {
			referenceEnd = time.Now().UTC()
		} else {
			referenceEnd = referenceEnd.UTC()
		}
		last30Days := referenceEnd.AddDate(0, 0, -30)
		expenseLast30Days, err := s.SumExpensesAmount(ctx, actor, biz, last30Days, referenceEnd)
		if err != nil {
			return decimal.Zero, err
		}
		safetyBuffer = expenseLast30Days
	}
	safeToDraw := totalIncome.Sub(totalCOGS).Sub(totalExpenses).Sub(totalWithdrawals).Sub(safetyBuffer)
	if safeToDraw.IsNegative() {
		return decimal.Zero, nil
	}
	return safeToDraw, nil
}

func (s *Service) GetExpenseByID(ctx context.Context, actor *account.User, biz *business.Business, id string) (*Expense, error) {
	return s.storage.expense.FindOne(ctx,
		s.storage.expense.ScopeID(id),
		s.storage.expense.ScopeBusinessID(biz.ID),
		s.storage.expense.WithPreload(RecurringExpenseStruct),
	)
}

func (s *Service) ListExpenses(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest) ([]*Expense, error) {
	return s.storage.expense.FindMany(ctx,
		s.storage.expense.ScopeBusinessID(biz.ID),
		s.storage.expense.WithPagination(req.Offset(), req.Limit()),
		s.storage.expense.WithOrderBy(req.ParsedOrderByWithDefault(ExpenseSchema, []string{"-occurredOn"})),
	)
}

func (s *Service) CountExpenses(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.expense.Count(ctx,
		s.storage.expense.ScopeBusinessID(biz.ID),
	)
}

func (s *Service) CreateExpense(ctx context.Context, actor *account.User, biz *business.Business, req *CreateExpenseRequest) (*Expense, error) {
	expense := &Expense{
		BusinessID:         biz.ID,
		Amount:             req.Amount,
		Currency:           biz.Currency,
		Category:           req.Category,
		Note:               transformer.ToNullString(req.Note),
		RecurringExpenseID: transformer.ToNullString(req.RecurringExpenseID),
		Type:               req.Type,
	}
	if req.OccurredOn != nil {
		expense.OccurredOn = *req.OccurredOn
	}
	if err := s.storage.expense.CreateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *Service) DeleteExpense(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	expense, err := s.GetExpenseByID(ctx, actor, biz, id)
	if err != nil {
		return err
	}
	return s.storage.expense.DeleteOne(ctx, expense)
}

func (s *Service) UpdateExpense(ctx context.Context, actor *account.User, biz *business.Business, id string, req *UpdateExpenseRequest) (*Expense, error) {
	expense, err := s.GetExpenseByID(ctx, actor, biz, id)
	if err != nil {
		return nil, err
	}
	if !req.Amount.IsZero() {
		expense.Amount = req.Amount
	}
	if req.Category != "" {
		expense.Category = req.Category
	}
	expense.Note = transformer.ToNullString(req.Note)
	if req.OccurredOn != nil {
		expense.OccurredOn = *req.OccurredOn
	}
	expense.RecurringExpenseID = transformer.ToNullString(req.RecurringExpenseID)
	if req.Type != "" {
		expense.Type = req.Type
	}
	if err := s.storage.expense.UpdateOne(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

// UpsertTransactionFeeExpenseForOrder creates (or updates) an idempotent transaction-fee expense linked to an order.
//
// Idempotency is enforced by a unique constraint on (business_id, order_id, category).
// This method is intended for internal automation (e.g., background events) and does not perform actor permission checks.
func (s *Service) UpsertTransactionFeeExpenseForOrder(
	ctx context.Context,
	businessID string,
	orderID string,
	amount decimal.Decimal,
	currency string,
	occurredOn time.Time,
	paymentMethod string,
) error {
	if businessID == "" || orderID == "" {
		return fmt.Errorf("businessID and orderID are required")
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil
	}
	if occurredOn.IsZero() {
		occurredOn = time.Now().UTC()
	}

	note := "Transaction fee"
	if paymentMethod != "" {
		note = fmt.Sprintf("Transaction fee (%s)", paymentMethod)
	}

	return s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		existing, err := s.storage.expense.FindOne(tctx,
			s.storage.expense.ScopeBusinessID(businessID),
			s.storage.expense.ScopeEquals(ExpenseSchema.OrderID, orderID),
			s.storage.expense.ScopeEquals(ExpenseSchema.Category, ExpenseCategoryTransactionFee),
			s.storage.expense.WithLockingStrength(database.LockingStrengthUpdate),
		)
		if err != nil && !database.IsRecordNotFound(err) {
			return err
		}
		if err == nil {
			existing.Amount = amount
			existing.Currency = currency
			existing.OccurredOn = occurredOn
			existing.Note = transformer.ToNullString(note)
			existing.Type = ExpenseTypeOneTime
			return s.storage.expense.UpdateOne(tctx, existing)
		}

		exp := &Expense{
			BusinessID: businessID,
			OrderID:    transformer.ToNullString(orderID),
			Amount:     amount,
			Currency:   currency,
			Category:   ExpenseCategoryTransactionFee,
			Note:       transformer.ToNullString(note),
			Type:       ExpenseTypeOneTime,
			OccurredOn: occurredOn,
		}
		if err := s.storage.expense.CreateOne(tctx, exp); err != nil {
			if database.IsUniqueViolation(err) {
				again, err2 := s.storage.expense.FindOne(tctx,
					s.storage.expense.ScopeBusinessID(businessID),
					s.storage.expense.ScopeEquals(ExpenseSchema.OrderID, orderID),
					s.storage.expense.ScopeEquals(ExpenseSchema.Category, ExpenseCategoryTransactionFee),
					s.storage.expense.WithLockingStrength(database.LockingStrengthUpdate),
				)
				if err2 != nil {
					return err2
				}
				again.Amount = amount
				again.Currency = currency
				again.OccurredOn = occurredOn
				again.Note = transformer.ToNullString(note)
				again.Type = ExpenseTypeOneTime
				return s.storage.expense.UpdateOne(tctx, again)
			}
			return err
		}
		return nil
	}, atomic.WithIsolationLevel(atomic.LevelSerializable), atomic.WithRetries(3))
}

func (s *Service) GetRecurringExpenseByID(ctx context.Context, actor *account.User, biz *business.Business, id string) (*RecurringExpense, error) {
	return s.storage.recurringExpense.FindOne(ctx,
		s.storage.recurringExpense.ScopeID(id),
		s.storage.recurringExpense.ScopeBusinessID(biz.ID),
		s.storage.recurringExpense.WithPreload("Expenses"),
	)
}

func (s *Service) ListRecurringExpenses(ctx context.Context, actor *account.User, biz *business.Business, req *list.ListRequest) ([]*RecurringExpense, error) {
	return s.storage.recurringExpense.FindMany(ctx,
		s.storage.recurringExpense.ScopeBusinessID(biz.ID),
		s.storage.recurringExpense.WithPagination(req.Offset(), req.Limit()),
		s.storage.recurringExpense.WithOrderBy(req.ParsedOrderBy(RecurringExpenseSchema)),
	)
}

func (s *Service) CountRecurringExpenses(ctx context.Context, actor *account.User, biz *business.Business) (int64, error) {
	return s.storage.recurringExpense.Count(ctx,
		s.storage.recurringExpense.ScopeBusinessID(biz.ID),
	)
}

func (s *Service) CreateRecurringExpense(ctx context.Context, actor *account.User, biz *business.Business, req *CreateRecurringExpenseRequest) (*RecurringExpense, error) {
	if !req.Amount.GreaterThan(decimal.Zero) {
		return nil, ErrRecurringExpenseInvalidAmount()
	}

	var recurringExpense *RecurringExpense
	if err := s.atomicProcessor.Exec(ctx, func(tctx context.Context) error {
		recurringExpense = &RecurringExpense{
			BusinessID:         biz.ID,
			Frequency:          req.Frequency,
			RecurringStartDate: req.RecurringStartDate,
			Amount:             req.Amount,
			Currency:           biz.Currency,
			Category:           req.Category,
			RecurringEndDate:   transformer.ToNullTime(req.RecurringEndDate),
			Note:               transformer.ToNullString(req.Note),
		}
		if err := s.storage.recurringExpense.CreateOne(tctx, recurringExpense); err != nil {
			return err
		}
		if req.AutoCreateHistoricalExpenses {
			err := s.backfillPastOccurrencesForCreate(tctx, biz, recurringExpense, time.Now())
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return recurringExpense, nil
}

func (s *Service) UpdateRecurringExpense(ctx context.Context, actor *account.User, biz *business.Business, id string, req *UpdateRecurringExpenseRequest) (*RecurringExpense, error) {
	recurringExpense, err := s.GetRecurringExpenseByID(ctx, actor, biz, id)
	if err != nil {
		return nil, err
	}
	if req.Frequency != "" {
		recurringExpense.Frequency = req.Frequency
	}
	if !req.RecurringStartDate.IsZero() {
		recurringExpense.RecurringStartDate = req.RecurringStartDate
	}
	if !req.Amount.IsZero() {
		if !req.Amount.GreaterThan(decimal.Zero) {
			return nil, ErrRecurringExpenseInvalidAmount()
		}
		recurringExpense.Amount = req.Amount
	}
	if req.Category != "" {
		recurringExpense.Category = req.Category
	}
	if !req.RecurringEndDate.IsZero() {
		recurringExpense.RecurringEndDate = transformer.ToNullTime(req.RecurringEndDate)
	}
	recurringExpense.Note = transformer.ToNullString(req.Note)
	if err := s.storage.recurringExpense.UpdateOne(ctx, recurringExpense); err != nil {
		return nil, err
	}
	return recurringExpense, nil
}

func (s *Service) backfillPastOccurrencesForCreate(ctx context.Context, biz *business.Business, re *RecurringExpense, today time.Time) error {
	current := re.RecurringStartDate
	expenses := make([]*Expense, 0)
	for current.Before(today) {
		expenses = append(expenses, &Expense{
			BusinessID:         biz.ID,
			Amount:             re.Amount,
			Currency:           re.Currency,
			Category:           re.Category,
			Note:               re.Note,
			RecurringExpenseID: transformer.ToNullString(re.ID),
			Type:               ExpenseTypeRecurring,
			OccurredOn:         current,
		})
		current = re.Frequency.GetNextRecurrenceDate(current)
	}
	return s.storage.expense.CreateMany(ctx, expenses)
}

func (s *Service) DeleteRecurringExpense(ctx context.Context, actor *account.User, biz *business.Business, id string) error {
	recurringExpense, err := s.GetRecurringExpenseByID(ctx, actor, biz, id)
	if err != nil {
		return err
	}
	return s.storage.recurringExpense.DeleteOne(ctx, recurringExpense)
}

func (s *Service) GetRecurringExpenseOccurrences(ctx context.Context, actor *account.User, biz *business.Business, rexpID string) ([]*Expense, error) {
	return s.storage.expense.FindMany(ctx,
		s.storage.expense.ScopeBusinessID(biz.ID),
		s.storage.expense.ScopeEquals(ExpenseSchema.RecurringExpenseID, rexpID),
		s.storage.expense.WithOrderBy([]string{"occurred_on ASC"}),
	)
}

func (s *Service) UpdateRecurringExpenseStatus(ctx context.Context, actor *account.User, biz *business.Business, rexpID string, newStatus RecurringExpenseStatus) (*RecurringExpense, error) {
	recurringExpense, err := s.GetRecurringExpenseByID(ctx, actor, biz, rexpID)
	if err != nil {
		return nil, err
	}
	sm := NewRecurringExpenseStateMachine(recurringExpense)
	if err := sm.TransitionTo(newStatus); err != nil {
		return nil, err
	}
	if err := s.storage.recurringExpense.UpdateOne(ctx, sm.RecurringExpense()); err != nil {
		return nil, err
	}
	return sm.RecurringExpense(), nil
}

func (s *Service) SumExpensesAmount(ctx context.Context, actor *account.User, biz *business.Business, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.expense.Sum(ctx,
		ExpenseSchema.Amount,
		s.storage.expense.ScopeBusinessID(biz.ID),
		s.storage.expense.ScopeTime(ExpenseSchema.OccurredOn, from, to),
	)
}

// SumExpensesAmountByCategory returns the total expense amount filtered by category within the given date range.
func (s *Service) SumExpensesAmountByCategory(ctx context.Context, actor *account.User, biz *business.Business, category ExpenseCategory, from, to time.Time) (decimal.Decimal, error) {
	return s.storage.expense.Sum(ctx,
		ExpenseSchema.Amount,
		s.storage.expense.ScopeBusinessID(biz.ID),
		s.storage.expense.ScopeEquals(ExpenseSchema.Category, category),
		s.storage.expense.ScopeTime(ExpenseSchema.OccurredOn, from, to),
	)
}

// this is not a multi-tenant scoped function it should be only used for internal processing
func (s *Service) ListAllActiveRecurringExpenses(ctx context.Context, from time.Time, to time.Time) ([]*RecurringExpense, error) {
	return s.storage.recurringExpense.FindMany(ctx,
		s.storage.recurringExpense.ScopeEquals(RecurringExpenseSchema.Status, RecurringExpenseStatusActive),
		s.storage.recurringExpense.ScopeTime(RecurringExpenseSchema.RecurringStartDate, from, to),
	)
}

func (s *Service) ListActiveRecurringExpenses(ctx context.Context, actor *account.User, business *business.Business, from time.Time, to time.Time) ([]*RecurringExpense, error) {
	return s.storage.recurringExpense.FindMany(ctx,
		s.storage.recurringExpense.ScopeEquals(RecurringExpenseSchema.Status, RecurringExpenseStatusActive),
		s.storage.recurringExpense.ScopeTime(RecurringExpenseSchema.RecurringStartDate, from, to),
	)
}

func (s *Service) GetLastRecurringExpenseOccurance(ctx context.Context, actor *account.User, business *business.Business, recurringExpenseID string) (*Expense, error) {
	return s.storage.expense.FindOne(ctx,
		s.storage.expense.ScopeEquals(ExpenseSchema.RecurringExpenseID, recurringExpenseID),
		s.storage.expense.WithOrderBy([]string{fmt.Sprintf("%s  DESC", ExpenseSchema.OccurredOn.Column())}),
	)
}

func (s *Service) CreateNewRecurringExpenseOccurrence(ctx context.Context, actor *account.User, business *business.Business, recurringExpense *RecurringExpense, occurrenceDate time.Time) (*Expense, error) {
	expense := &Expense{
		BusinessID:         business.ID,
		Amount:             recurringExpense.Amount,
		Currency:           recurringExpense.Currency,
		Category:           recurringExpense.Category,
		Note:               transformer.ToNullString("This is an auto-generated expense for recurring expense: " + recurringExpense.ID),
		RecurringExpenseID: transformer.ToNullString(recurringExpense.ID),
		Type:               ExpenseTypeRecurring,
		OccurredOn:         occurrenceDate,
	}
	if err := s.atomicProcessor.Exec(ctx, func(ctx context.Context) error {
		if err := s.storage.expense.CreateOne(ctx, expense); err != nil {
			return err
		}
		// update recurring expense next occurrence date
		recurringExpense.NextRecurringDate = recurringExpense.Frequency.GetNextRecurrenceDate(occurrenceDate)
		return s.storage.recurringExpense.UpdateOne(ctx, recurringExpense)
	}); err != nil {
		return nil, err
	}
	return expense, nil
}
