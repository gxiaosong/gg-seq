package data

import (
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gouez/gg-seq/app/internal/biz"
)

func Test_seqRepo_GetNextSegment(t *testing.T) {
	type fields struct {
		data *Data
		log  *log.Helper
	}
	type args struct {
		bizType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *biz.Segment
		wantErr bool
	}{
		{
			args:   args{bizType: "test"},
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &seqRepo{
				data: tt.fields.data,
				log:  tt.fields.log,
			}
			got, err := repo.GetNextSegment(tt.args.bizType)
			if (err != nil) != tt.wantErr {
				t.Errorf("seqRepo.GetNextSegment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("seqRepo.GetNextSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}
