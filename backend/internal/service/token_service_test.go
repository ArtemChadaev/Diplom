package service_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestTokenService is a helper that creates a TokenService with known settings.
func newTestTokenService() domain.TokenService {
	return service.NewTokenService("test-secret", 15*time.Minute, 15*24*time.Hour)
}

func testUser(id int, role domain.UserRole) *domain.User {
	return &domain.User{
		ID:    id,
		Email: "test@example.com",
		Role:  role,
	}
}

// ---------------------------------------------------------------------------
// GenerateAccessToken + ParseAccessToken
// ---------------------------------------------------------------------------

// Scenario: JWT round-trip
//
//	Given:  A user with id=42, role=admin, email=test@example.com and a sessionID
//	When:   GenerateAccessToken is called, then ParseAccessToken is called on the result
//	Then:   Claims match the original inputs exactly
func TestTokenService_GenerateAndParseAccessToken_HappyPath(t *testing.T) {
	svc := newTestTokenService()
	user := testUser(42, domain.RoleAdmin)
	sessionID := uuid.New()

	raw, err := svc.GenerateAccessToken(user, sessionID)
	require.NoError(t, err)
	assert.NotEmpty(t, raw)

	claims, err := svc.ParseAccessToken(raw)
	require.NoError(t, err)
	require.NotNil(t, claims)

	assert.Equal(t, 42, claims.UserID)
	assert.Equal(t, domain.RoleAdmin, claims.Role)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, sessionID, claims.SessionID)
}

// Scenario: Parse token signed with wrong secret
//
//	Given:  A token generated with secret "test-secret"
//	When:   ParseAccessToken is called by a service using "wrong-secret"
//	Then:   Error is returned, claims are nil
func TestTokenService_ParseAccessToken_WrongSecret(t *testing.T) {
	svcA := newTestTokenService()
	svcB := service.NewTokenService("wrong-secret", 15*time.Minute, 15*24*time.Hour)

	user := testUser(1, domain.RolePharmacist)
	raw, err := svcA.GenerateAccessToken(user, uuid.New())
	require.NoError(t, err)

	claims, err := svcB.ParseAccessToken(raw)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// Scenario: Parse expired token
//
//	Given:  A token with TTL = -1 second (already expired at generation time)
//	When:   ParseAccessToken is called
//	Then:   Error is returned
func TestTokenService_ParseAccessToken_Expired(t *testing.T) {
	// TTL of -1 second → token is expired the moment it's created
	svc := service.NewTokenService("test-secret", -1*time.Second, 15*24*time.Hour)
	raw, err := svc.GenerateAccessToken(testUser(1, domain.RolePharmacist), uuid.New())
	require.NoError(t, err)

	claims, err := svc.ParseAccessToken(raw)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// Scenario: Tampered token payload
//
//	Given:  A valid JWT whose payload has been modified after signing
//	When:   ParseAccessToken is called
//	Then:   Signature verification fails → error
func TestTokenService_ParseAccessToken_Tampered(t *testing.T) {
	svc := newTestTokenService()
	raw, err := svc.GenerateAccessToken(testUser(1, domain.RolePharmacist), uuid.New())
	require.NoError(t, err)

	// Flip the last character of the signature segment
	tampered := raw[:len(raw)-1] + "X"
	claims, err := svc.ParseAccessToken(tampered)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// ---------------------------------------------------------------------------
// HashToken
// ---------------------------------------------------------------------------

// Scenario: SHA-256 hash is deterministic and collision-resistant
//
//	Given:  Two distinct inputs "abc123" and "abc124"
//	Then:   Same input → same output (deterministic)
//	And:    Different inputs → different outputs
func TestTokenService_HashToken(t *testing.T) {
	svc := newTestTokenService()

	h1 := svc.HashToken("abc123")
	h2 := svc.HashToken("abc123")
	h3 := svc.HashToken("abc124")

	assert.Equal(t, h1, h2, "hash must be deterministic")
	assert.NotEqual(t, h1, h3, "different inputs must produce different hashes")
	assert.Len(t, h1, 64, "SHA-256 hex digest must be 64 characters")
}

// ---------------------------------------------------------------------------
// GenerateRefreshToken
// ---------------------------------------------------------------------------

// Scenario: Refresh tokens are unique and hash is consistent
//
//	Given:  Two calls to GenerateRefreshToken
//	Then:   Raw values are different (collision-free)
//	And:    HashToken(raw) == hash returned by GenerateRefreshToken
func TestTokenService_GenerateRefreshToken(t *testing.T) {
	svc := newTestTokenService()

	raw1, hash1, err := svc.GenerateRefreshToken()
	require.NoError(t, err)

	raw2, hash2, err := svc.GenerateRefreshToken()
	require.NoError(t, err)

	assert.NotEmpty(t, raw1)
	assert.NotEmpty(t, raw2)
	assert.NotEqual(t, raw1, raw2, "tokens must be unique across calls")
	assert.NotEqual(t, hash1, hash2, "hashes must differ for different tokens")

	// Hash consistency: manually hashing raw must match the returned hash
	assert.Equal(t, svc.HashToken(raw1), hash1, "manual hash must equal returned hash")
	assert.Equal(t, svc.HashToken(raw2), hash2, "manual hash must equal returned hash")
}
