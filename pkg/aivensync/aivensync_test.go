package aivensync

import (
	"testing"
)

func TestFindLastAckedEvent(t *testing.T) {
	events := []*AivenEvent{
		{Actor: "a", EventDesc: "a", EventType: "a", ServiceName: "a", Time: "a"},
		{Actor: "b", EventDesc: "b", EventType: "b", ServiceName: "b", Time: "b"},
		{Actor: "c", EventDesc: "c", EventType: "c", ServiceName: "c", Time: "c"},
	}

	type args struct {
		events         []*AivenEvent
		lastAckedEvent *AivenEvent
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "find a event",
			want: 0,
			args: args{
				lastAckedEvent: &AivenEvent{Actor: "a", EventDesc: "a", EventType: "a", ServiceName: "a", Time: "a"},
				events:         events,
			},
		},
		{
			name: "find b event",
			want: 1,
			args: args{
				lastAckedEvent: &AivenEvent{Actor: "b", EventDesc: "b", EventType: "b", ServiceName: "b", Time: "b"},
				events:         events,
			},
		},
		{
			name: "find c event",
			want: 2,
			args: args{
				lastAckedEvent: &AivenEvent{Actor: "c", EventDesc: "c", EventType: "c", ServiceName: "c", Time: "c"},
				events:         events,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindLastAckedEvent(tt.args.events, tt.args.lastAckedEvent); got != tt.want {
				t.Errorf("FindLastAckedEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
