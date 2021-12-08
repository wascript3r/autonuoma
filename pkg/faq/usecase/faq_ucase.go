package usecase

import (
	"context"
	"time"

	"github.com/wascript3r/autonuoma/pkg/faq"
)

type Usecase struct {
	faqRepo    faq.Repository
	ctxTimeout time.Duration
}

func New(fr faq.Repository, t time.Duration) *Usecase {
	return &Usecase{
		faqRepo:    fr,
		ctxTimeout: t,
	}
}

func (u *Usecase) GetAll(ctx context.Context) (*faq.GetAllRes, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	fs, err := u.faqRepo.GetAll(c)
	if err != nil {
		return nil, err
	}

	faqs := make([]*faq.FAQListInfo, len(fs))
	for i, f := range fs {
		faqs[i] = &faq.FAQListInfo{
			ID:         f.ID,
			CategoryID: f.CategoryID,
			Question:   f.Question,
			Answer:     f.Answer,
		}
	}

	return &faq.GetAllRes{
		FAQ: faqs,
	}, nil
}
