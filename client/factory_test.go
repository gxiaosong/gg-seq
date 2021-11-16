package client

import (
	"reflect"
	"testing"

	"github.com/gouez/gg-seq/comm"
)

func Test_clientIdGeneratorFactory_GetIdGenerator(t *testing.T) {
	type fields struct {
		cache   map[string]comm.IdGenerator
		service comm.SegmentService
	}
	type args struct {
		bizType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   comm.IdGenerator
	}{
		{
			fields: fields{
				service: NewHttpSegmentService("http://127.0.0.1:8888"),
				cache:   map[string]comm.IdGenerator{},
			},
			args: args{bizType: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &clientIdGeneratorFactory{
				cache:   tt.fields.cache,
				service: tt.fields.service,
			}
			if got := d.GetIdGenerator(tt.args.bizType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clientIdGeneratorFactory.GetIdGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}
