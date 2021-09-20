package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/session"
	"github.com/wascript3r/gocipher/encoder"
)

const (
	ctxKey = "session"
)

type Usecase struct {
	sessionRepo session.Repository
	ctxTimeout  time.Duration

	generator session.Generator
	cipher    session.Cipher
	opts      *Options
}

func New(sr session.Repository, t time.Duration, g session.Generator, c session.Cipher, opt ...Option) *Usecase {
	return &Usecase{
		sessionRepo: sr,
		ctxTimeout:  t,

		generator: g,
		cipher:    c,
		opts:      newOptions(opt...),
	}
}

func (u *Usecase) Create(ctx context.Context, userID int) (*domain.Session, error) {
	id, err := u.generator.GenerateID()
	if err != nil {
		return nil, err
	}

	ss := &domain.Session{
		ID:         id,
		UserID:     userID,
		Expiration: time.Now().Add(u.opts.SessionLifetime),
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	err = u.sessionRepo.Insert(c, ss)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func (u *Usecase) IsExpired(ss *domain.Session) bool {
	return time.Now().After(ss.Expiration)
}

func (u *Usecase) Validate(ctx context.Context, id string) (*domain.Session, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	s, err := u.sessionRepo.Get(c, id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, session.NotAuthenticatedError
		}
		return nil, err
	}

	if u.IsExpired(s) {
		u.sessionRepo.Delete(c, id)
		return nil, session.SessionExpiredError
	}

	return s, nil
}

func (u *Usecase) Delete(ctx context.Context, id string) error {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.sessionRepo.Delete(c, id)
}

func (u *Usecase) GenTempToken(ss *domain.Session) (string, error) {
	exp := time.Now().Add(u.opts.TokenLifetime)
	if exp.After(ss.Expiration) {
		exp = ss.Expiration
	}

	t := &session.TempToken{
		SessionID:  ss.ID,
		Expiration: exp,
	}

	encoded, err := t.GobEncode()
	if err != nil {
		return "", err
	}

	encrypted, err := u.cipher.Encrypt(encoded)
	if err != nil {
		return "", err
	}

	return string(encoder.Base64Encode(encrypted)), nil
}

func (u *Usecase) ValidateTempToken(ctx context.Context, token string) (*domain.Session, error) {
	encrypted, err := encoder.Base64Decode([]byte(token))
	if err != nil {
		return nil, err
	}

	encoded, err := u.cipher.Decrypt(encrypted)
	if err != nil {
		return nil, err
	}

	t := &session.TempToken{}
	err = t.GobDecode(encoded)
	if err != nil {
		return nil, err
	}

	if time.Now().After(t.Expiration) {
		return nil, session.TokenExpiredError
	}

	return u.Validate(ctx, t.SessionID)
}

func (u *Usecase) StoreCtx(ctx context.Context, ss *domain.Session) context.Context {
	return context.WithValue(ctx, ctxKey, ss)
}

func (u *Usecase) LoadCtx(ctx context.Context) (*domain.Session, error) {
	s, ok := ctx.Value(ctxKey).(*domain.Session)
	if !ok {
		return nil, session.ErrCannotLoadSession
	}
	return s, nil
}
