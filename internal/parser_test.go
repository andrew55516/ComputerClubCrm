package internal

import (
	"reflect"
	"testing"
)

func TestParseEvent(t *testing.T) {
	type args struct {
		event        string
		tablesAmount int
	}
	tests := []struct {
		name      string
		args      args
		wantEvent Event
		wantOk    bool
	}{
		{
			name: "OK1",
			args: args{
				event:        "08:48 1 client1",
				tablesAmount: 3,
			},
			wantEvent: Event{
				ID:    1,
				Time:  528,
				Name:  "client1",
				Table: 0,
			},
			wantOk: true,
		},
		{
			name: "OK2",
			args: args{
				event:        "08:48 2 client1 3",
				tablesAmount: 3,
			},
			wantEvent: Event{
				ID:    2,
				Time:  528,
				Name:  "client1",
				Table: 3,
			},
			wantOk: true,
		},
		{
			name: "wrong time1",
			args: args{
				event:        "38:48 1 client1",
				tablesAmount: 3,
			},
			wantOk: false,
		},
		{
			name: "wrong time2",
			args: args{
				event:        "08-48 1 client1",
				tablesAmount: 3,
			},
			wantOk: false,
		},
		{
			name: "wrong time3",
			args: args{
				event:        "8:48 1 client1",
				tablesAmount: 3,
			},
			wantOk: false,
		},
		{
			name: "wrong time4",
			args: args{
				event:        "28:48 1 client1",
				tablesAmount: 3,
			},
			wantOk: false,
		},
		{
			name: "wrong id1",
			args: args{
				event:        "08:48 8 client1 2",
				tablesAmount: 3,
			},
			wantOk: false,
		},
		{
			name: "wrong id2",
			args: args{
				event:        "08:48 a client1 2",
				tablesAmount: 3,
			},
			wantOk: false,
		},
		{
			name: "wrong id3",
			args: args{
				event:        "08:48 1 client1 2",
				tablesAmount: 3,
			},
			wantOk: false,
		},
		{
			name: "wrong name",
			args: args{
				event:        "08:48 1 Client1",
				tablesAmount: 3,
			},
			wantOk: false,
		},
		{
			name: "wrong table1",
			args: args{
				event:        "08:48 2 client1 0",
				tablesAmount: 3,
			},
			wantOk: false,
		},
		{
			name: "wrong table2",
			args: args{
				event:        "08:48 2 client1 4",
				tablesAmount: 3,
			},
			wantOk: false,
		},
		{
			name: "wrong event",
			args: args{
				event:        "08:48 2",
				tablesAmount: 3,
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEvent, gotOk := ParseEvent(tt.args.event, tt.args.tablesAmount)
			if !reflect.DeepEqual(gotOk, tt.wantOk) {
				t.Errorf("ParseEvent() gotOk = %v, wantOk %v", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotEvent, tt.wantEvent) && tt.wantOk {
				t.Errorf("ParseEvent() gotEvent = %v, want %v", gotEvent, tt.wantEvent)
			}
		})
	}
}
