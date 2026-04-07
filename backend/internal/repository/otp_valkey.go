package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/valkey-io/valkey-go"
)

type otpValkeyRepository struct {
	client valkey.Client
}

// NewOTPValkeyRepository creates a new Valkey-based OTP repository.
func NewOTPValkeyRepository(client valkey.Client) domain.OTPRepository {
	return &otpValkeyRepository{client: client}
}

// hashKey returns the valkey key for storing the OTP hash
func hashKey(userID int) string {
	return "otp:" + strconv.Itoa(userID) + ":hash"
}

// attemptsKey returns the valkey key for storing the OTP attempts
func attemptsKey(userID int) string {
	return "otp:" + strconv.Itoa(userID) + ":attempts"
}

func (r *otpValkeyRepository) Store(ctx context.Context, userID int, codeHash string) error {
	hk := hashKey(userID)
	ak := attemptsKey(userID)
	ttl := time.Duration(domain.OTPTTLSeconds) * time.Second

	// We can use a pipeline or multiple commands.
	err := r.client.Do(ctx, r.client.B().Set().Key(hk).Value(codeHash).Ex(ttl).Build()).Error()
	if err != nil {
		return err
	}

	err = r.client.Do(ctx, r.client.B().Set().Key(ak).Value("0").Ex(ttl).Build()).Error()
	return err
}

func (r *otpValkeyRepository) Get(ctx context.Context, userID int) (*domain.OTPCode, error) {
	hk := hashKey(userID)
	ak := attemptsKey(userID)

	hashVal, err := r.client.Do(ctx, r.client.B().Get().Key(hk).Build()).ToString()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return nil, domain.ErrOTPNotFound
		}
		return nil, err
	}

	attemptsStr, err := r.client.Do(ctx, r.client.B().Get().Key(ak).Build()).ToString()
	if err != nil && !errors.Is(err, valkey.Nil) {
		return nil, err
	}

	attempts, _ := strconv.Atoi(attemptsStr) // ignore error, default to 0

	return &domain.OTPCode{
		UserID:   userID,
		CodeHash: hashVal,
		Attempts: attempts,
	}, nil
}

func (r *otpValkeyRepository) IncrAttempts(ctx context.Context, userID int) error {
	ak := attemptsKey(userID)
	// INCR doesn't fail if key doesn't exist, it creates it, but we only increment if we verified existance
	return r.client.Do(ctx, r.client.B().Incr().Key(ak).Build()).Error()
}

func (r *otpValkeyRepository) Delete(ctx context.Context, userID int) error {
	hk := hashKey(userID)
	ak := attemptsKey(userID)
	return r.client.Do(ctx, r.client.B().Del().Key(hk).Key(ak).Build()).Error()
}
