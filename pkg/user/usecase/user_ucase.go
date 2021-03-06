package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/wascript3r/autonuoma/pkg/domain"
	"github.com/wascript3r/autonuoma/pkg/session"
	"github.com/wascript3r/autonuoma/pkg/user"
)

type Usecase struct {
	userRepo   user.Repository
	ctxTimeout time.Duration

	sessionUcase session.Usecase
	pwHasher     user.PwHasher
	validate     user.Validate
}

func New(ur user.Repository, t time.Duration, su session.Usecase, ph user.PwHasher, v user.Validate) *Usecase {
	return &Usecase{
		userRepo:   ur,
		ctxTimeout: t,

		sessionUcase: su,
		pwHasher:     ph,
		validate:     v,
	}
}

func (u *Usecase) Create(ctx context.Context, req *user.CreateReq) error {
	if err := u.validate.RawRequest(req); err != nil {
		return user.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	err := u.validate.EmailUniqueness(c, req.Email)
	if err != nil {
		if err == user.ErrEmailExists {
			return user.EmailAlreadyExistsError
		}
		return err
	}

	hash, err := u.pwHasher.Hash(req.Password)
	if err != nil {
		return err
	}

	us := &domain.User{
		Email:     req.Email,
		Password:  hash,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: time.Time(req.BirthDate),
		Balance:   0,
		PIN:       req.PIN,
		RoleID:    domain.ClientRole,
	}

	err = u.userRepo.InsertIfNotExists(c, us)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) Authenticate(ctx context.Context, req *user.AuthenticateReq) (*domain.Session, *user.AuthenticateRes, error) {
	if err := u.validate.RawRequest(req); err != nil {
		return nil, nil, user.InvalidInputError
	}

	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	credentials, err := u.userRepo.GetCredentials(c, req.Email)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, nil, user.InvalidCredentialsError
		}
		return nil, nil, err
	}

	err = u.pwHasher.Validate(credentials.Password, req.Password)
	if err != nil {
		return nil, nil, user.InvalidCredentialsError
	}

	s, err := u.sessionUcase.Create(ctx, credentials.ID)
	if err != nil {
		return nil, nil, err
	}

	res := &user.AuthenticateRes{
		UserID: credentials.ID,
		RoleID: credentials.RoleID,
	}

	return s, res, nil
}

func (u *Usecase) GetTempToken(ss *domain.Session) (*user.TempToken, error) {
	token, err := u.sessionUcase.GenTempToken(ss)
	if err != nil {
		return nil, err
	}
	return &user.TempToken{Token: token}, nil
}

func (u *Usecase) AuthenticateToken(ctx context.Context, req *user.TempToken) (*domain.Session, error) {
	return u.sessionUcase.ValidateTempToken(ctx, req.Token)
}

func (u *Usecase) Logout(ctx context.Context, ss *domain.Session) error {
	return u.sessionUcase.Delete(ctx, ss.ID)
}

func (u *Usecase) GetInfo(userID int, role domain.Role) *user.AuthenticateRes {
	return &user.AuthenticateRes{
		UserID: userID,
		RoleID: role,
	}
}

func (u *Usecase) GetData(ctx context.Context, uid int) (*user.UserProfile, error) {
	user, err := u.userRepo.GetData(ctx, uid)
	if err != nil {
		return nil, err
	}

	licenseStatus, err := u.userRepo.GetLicenseStatus(ctx, uid)
	if err != nil {
		return nil, err
	}

	user.LicenseStatus = licenseStatus

	return user, nil
}

func (u *Usecase) UpdateUser(ctx context.Context, uid int, data *user.UpdateReq) (*user.UpdateRes, error) {
	if len(data.Email) > 0 {
		exists, err := u.userRepo.EmailExists(ctx, data.Email)
		if err != nil || exists {
			return nil, err
		}

		if err := u.userRepo.UpdateEmail(ctx, uid, data.Email); err != nil {
			return nil, err
		}
	}

	if len(data.Password) > 0 {
		hash, err := u.pwHasher.Hash(data.Password)
		if err != nil {
			return nil, err
		}

		if err := u.userRepo.UpdatePassword(ctx, uid, hash); err != nil {
			return nil, err
		}
	}

	result, err := u.userRepo.GetData(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &user.UpdateRes{Email: result.Email}, nil
}

func (u *Usecase) GetTrips(ctx context.Context, uid int) ([]*user.TripsRes, error) {
	res, err := u.userRepo.GetTrips(ctx, uid)
	if err != nil {
		return nil, err
	}

	trips := make([]*user.TripsRes, 0)
	for _, t := range res {
		trips = append(trips, &user.TripsRes{
			ID:    t.ID,
			Begin: t.Begin.Format(user.TripDateTimeFormat),
			End:   t.End.Format(user.TripDateTimeFormat),
			From:  "",
			To:    "",
			Price: t.Price,
		})
	}

	return trips, nil
}

func (u *Usecase) CheckPayment(ctx context.Context, uid int) (*user.PaymentRes, error) {
	rand.Seed(time.Now().UnixNano())

	amount := rand.Intn(50) * 100
	var status string

	switch rand.Intn(2) {
	case 0:
		status = "successful"
		if err := u.userRepo.AddPayment(ctx, uid, int64(amount)); err != nil {
			return nil, err
		}

		if err := u.userRepo.AddBalance(ctx, uid, int64(amount)); err != nil {
			return nil, err
		}
	case 1:
		status = "cancelled"
		if err := u.userRepo.AddPayment(ctx, uid, 0); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown action")
	}

	data, err := u.GetData(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &user.PaymentRes{
		Status:  status,
		Balance: data.Balance,
	}, nil
}
