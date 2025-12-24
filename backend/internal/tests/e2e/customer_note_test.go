package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

type CustomerNoteSuite struct {
	suite.Suite
	accountHelper  *AccountTestHelper
	customerHelper *CustomerTestHelper
}

func (s *CustomerNoteSuite) SetupSuite() {
	s.accountHelper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
	s.customerHelper = NewCustomerTestHelper(testEnv.Database, e2eBaseURL)
}

func (s *CustomerNoteSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "customers", "customer_notes", "subscriptions")
	s.NoError(err)
}

func (s *CustomerNoteSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "businesses", "customers", "customer_notes", "subscriptions")
	s.NoError(err)
}

func (s *CustomerNoteSuite) TestCreateNote_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	payload := map[string]interface{}{
		"content": "This is a test note about the customer",
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/notes", customer.ID), payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.NotEmpty(result["id"])
	s.Equal(customer.ID, result["customerId"])
	s.Equal("This is a test note about the customer", result["content"])
	s.Contains(result, "createdAt")

	// Verify DB state
	count, err := s.customerHelper.CountNotes(ctx, customer.ID)
	s.NoError(err)
	s.Equal(int64(1), count)
}

func (s *CustomerNoteSuite) TestCreateNote_ValidationErrors() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	tests := []struct {
		name    string
		payload map[string]interface{}
	}{
		{"missing content", map[string]interface{}{}},
		{"empty content", map[string]interface{}{"content": ""}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.customerHelper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/notes", customer.ID), tt.payload, token)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func (s *CustomerNoteSuite) TestCreateNote_CustomerNotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	payload := map[string]interface{}{
		"content": "Test note",
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("POST", "/v1/businesses/test-biz/customers/cus_nonexistent/notes", payload, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerNoteSuite) TestListNotes_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	// Create multiple notes
	for i := 0; i < 3; i++ {
		_, err := s.customerHelper.CreateTestNote(ctx, customer.ID, fmt.Sprintf("Note %d", i))
		s.NoError(err)
	}

	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/notes", customer.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result []interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Len(result, 3)

	// Verify all notes belong to the customer
	for _, noteInterface := range result {
		note := noteInterface.(map[string]interface{})
		s.Equal(customer.ID, note["customerId"])
		s.Contains(note, "content")
		s.Contains(note, "id")
	}
}

func (s *CustomerNoteSuite) TestListNotes_CustomerNotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", "/v1/businesses/test-biz/customers/cus_nonexistent/notes", nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerNoteSuite) TestDeleteNote_Success() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	note, err := s.customerHelper.CreateTestNote(ctx, customer.ID, "Test note")
	s.NoError(err)

	resp, err := s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/notes/%s", customer.ID, note.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify note is soft deleted
	count, err := s.customerHelper.CountNotes(ctx, customer.ID)
	s.NoError(err)
	s.Equal(int64(0), count)
}

func (s *CustomerNoteSuite) TestDeleteNote_NoteNotFound() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	resp, err := s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/notes/cnote_nonexistent", customer.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *CustomerNoteSuite) TestDeleteNote_WrongCustomer() {
	ctx := context.Background()
	_, ws, token, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer1, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "customer1@example.com", "Customer 1")
	s.NoError(err)

	customer2, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "customer2@example.com", "Customer 2")
	s.NoError(err)

	note1, err := s.customerHelper.CreateTestNote(ctx, customer1.ID, "Customer 1 note")
	s.NoError(err)

	// Try to delete customer1's note through customer2's endpoint
	resp, err := s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/notes/%s", customer2.ID, note1.ID), nil, token)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// Verify note still exists
	count, err := s.customerHelper.CountNotes(ctx, customer1.ID)
	s.NoError(err)
	s.Equal(int64(1), count)
}

func (s *CustomerNoteSuite) TestNoteOperations_CrossWorkspaceIsolation() {
	ctx := context.Background()

	// Workspace 1
	_, ws1, token1, err := s.accountHelper.CreateTestUser(ctx, "admin1@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws1.ID))
	biz1, err := s.customerHelper.CreateTestBusiness(ctx, ws1.ID, "biz-1")
	s.NoError(err)
	customer1, err := s.customerHelper.CreateTestCustomer(ctx, biz1.ID, "customer1@example.com", "Customer 1")
	s.NoError(err)
	note1, err := s.customerHelper.CreateTestNote(ctx, customer1.ID, "WS1 note")
	s.NoError(err)

	// Workspace 2
	_, ws2, token2, err := s.accountHelper.CreateTestUser(ctx, "admin2@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws2.ID))
	_, err = s.customerHelper.CreateTestBusiness(ctx, ws2.ID, "biz-2")
	s.NoError(err)

	// WS2 should not list WS1's customer notes
	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/biz-1/customers/%s/notes", customer1.ID), nil, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// WS2 should not create note for WS1's customer
	payload := map[string]interface{}{"content": "Hacked note"}
	resp, err = s.customerHelper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/biz-1/customers/%s/notes", customer1.ID), payload, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// WS2 should not delete WS1's note
	resp, err = s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/biz-1/customers/%s/notes/%s", customer1.ID, note1.ID), nil, token2)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// Verify WS1 note unchanged and exists
	count, err := s.customerHelper.CountNotes(ctx, customer1.ID)
	s.NoError(err)
	s.Equal(int64(1), count)

	// WS1 can still access their own notes
	resp, err = s.customerHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/biz-1/customers/%s/notes", customer1.ID), nil, token1)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)
}

func (s *CustomerNoteSuite) TestNoteOperations_UnauthorizedUser() {
	ctx := context.Background()

	// Admin creates customer
	_, ws, _, err := s.accountHelper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)
	s.NoError(s.accountHelper.CreateTestSubscription(ctx, ws.ID))

	biz, err := s.customerHelper.CreateTestBusiness(ctx, ws.ID, "test-biz")
	s.NoError(err)

	customer, err := s.customerHelper.CreateTestCustomer(ctx, biz.ID, "john@example.com", "John Doe")
	s.NoError(err)

	note, err := s.customerHelper.CreateTestNote(ctx, customer.ID, "Test note")
	s.NoError(err)

	// Regular user (non-admin) in the SAME workspace tries to perform operations
	userPassword, err := hash.Password("Password123!")
	s.NoError(err)
	member := &account.User{
		WorkspaceID:     ws.ID,
		Role:            role.RoleUser,
		FirstName:       "User",
		LastName:        "User",
		Email:           "user@example.com",
		Password:        userPassword,
		IsEmailVerified: true,
	}
	userRepo := database.NewRepository[account.User](testEnv.Database)
	s.NoError(userRepo.CreateOne(ctx, member))
	userToken, err := auth.NewJwtToken(member.ID, member.WorkspaceID, member.AuthVersion)
	s.NoError(err)

	// User can view notes (has view permission)
	resp, err := s.customerHelper.Client.AuthenticatedRequest("GET", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/notes", customer.ID), nil, userToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	// User cannot create notes (needs manage permission)
	payload := map[string]interface{}{"content": "User note"}
	resp, err = s.customerHelper.Client.AuthenticatedRequest("POST", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/notes", customer.ID), payload, userToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)

	// User cannot delete notes (needs manage permission)
	resp, err = s.customerHelper.Client.AuthenticatedRequest("DELETE", fmt.Sprintf("/v1/businesses/test-biz/customers/%s/notes/%s", customer.ID, note.ID), nil, userToken)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusForbidden, resp.StatusCode)
}

func TestCustomerNoteSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(CustomerNoteSuite))
}
