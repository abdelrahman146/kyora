package accounting

import "slices"

type RecurringExpenseStateMachine struct {
	recurringExpense *RecurringExpense
}

func NewRecurringExpenseStateMachine(rexp *RecurringExpense) *RecurringExpenseStateMachine {
	return &RecurringExpenseStateMachine{recurringExpense: rexp}
}

func (sm *RecurringExpenseStateMachine) CanTransitionTo(newStatus RecurringExpenseStatus) bool {
	allowedTransitions := map[RecurringExpenseStatus][]RecurringExpenseStatus{
		RecurringExpenseStatusActive:   {RecurringExpenseStatusPaused, RecurringExpenseStatusEnded, RecurringExpenseStatusCanceled},
		RecurringExpenseStatusPaused:   {RecurringExpenseStatusActive, RecurringExpenseStatusEnded, RecurringExpenseStatusCanceled},
		RecurringExpenseStatusEnded:    {RecurringExpenseStatusActive, RecurringExpenseStatusCanceled},
		RecurringExpenseStatusCanceled: {RecurringExpenseStatusActive},
	}

	if transitions, ok := allowedTransitions[sm.recurringExpense.Status]; ok {
		if slices.Contains(transitions, newStatus) {
			return true
		}
	}
	return false
}

func (sm *RecurringExpenseStateMachine) TransitionTo(newStatus RecurringExpenseStatus) error {
	if !sm.CanTransitionTo(newStatus) {
		return ErrRecurringExpenseInvalidTransition(string(sm.recurringExpense.Status), string(newStatus))
	}
	sm.recurringExpense.Status = newStatus
	return nil
}

func (sm *RecurringExpenseStateMachine) RecurringExpense() *RecurringExpense {
	return sm.recurringExpense
}
