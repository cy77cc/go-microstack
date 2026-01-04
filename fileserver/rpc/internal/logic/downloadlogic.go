package logic

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/fileserver/model"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadLogic {
	return &DownloadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DownloadLogic) Download(in *pb.DownloadReq) (*pb.DownloadResp, error) {
	var bucket, objectName, fileName string
	if in.FileId != "" {
		fi, err := l.svcCtx.FileModel.FindOneByFileId(l.ctx, in.FileId)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return nil, xcode.NewErrCodeMsg(xcode.NotFound, "file not exist")
			}
			return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
		}
		bucket = fi.Bucket
		objectName = fi.ObjectName
		fileName = fi.FileName
	} else {
		bucket = in.Bucket
		objectName = in.ObjectName
		fileName = filepath.Base(objectName)
	}
	stor, err := l.svcCtx.Storage.Select(l.ctx, bucket)
	if err != nil {
		return nil, err
	}
	data, contentType, err := stor.GetObject(l.ctx, bucket, objectName)
	if err != nil {
		return nil, err
	}

	return &pb.DownloadResp{
		Data:        data,
		ContentType: contentType,
		Size:        int64(len(data)),
		FileName:    fileName,
	}, nil
}
