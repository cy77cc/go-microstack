package logic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileMetaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileMetaLogic {
	return &GetFileMetaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFileMetaLogic) GetFileMeta(in *pb.GetFileMetaReq) (*pb.GetFileMetaResp, error) {
	if in.FileId == "" {
		return nil, errors.New("fileId empty")
	}
	fi, err := l.svcCtx.FileModel.FindOneByFileId(l.ctx, in.FileId)
	if err != nil {
		return nil, err
	}
	return &pb.GetFileMetaResp{
		Meta: &pb.FileMeta{
			FileId:      fi.FileId,
			FileName:    fi.FileName,
			Bucket:      fi.Bucket,
			ObjectName:  fi.ObjectName,
			Size:        fi.Size,
			ContentType: fi.ContentType,
			Uploader:    fi.Uploader,
			UploadTime:  fi.UploadTime,
			Hash:        fi.Hash,
			Status:      fi.Status,
		},
	}, nil

}
