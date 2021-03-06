package license

import (
	"context"
)

type Usecase interface {
	Confirm(ctx context.Context, req *ChangeStatusReq) error
	Reject(ctx context.Context, req *ChangeStatusReq) error
	GetAllUnconfirmed(ctx context.Context) (*GetAllRes, error)
	GetPhotos(ctx context.Context, req *GetPhotosReq) (*GetPhotosRes, error)
	Upload(ctx context.Context, req *UploadReq) (*UploadRes, error)
}
