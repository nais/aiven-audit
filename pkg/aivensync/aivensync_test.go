package aivensync

import (
	"testing"
)

func TestFindStartIndex(t *testing.T) {
	events := []*AivenEvent{
		{ID: "a", Actor: "a", EventDesc: "a", EventType: "a", ServiceName: "a", Time: "a"},
		{ID: "d", Actor: "a", EventDesc: "a", EventType: "a", ServiceName: "a", Time: "a"},
		{ID: "b", Actor: "b", EventDesc: "b", EventType: "b", ServiceName: "b", Time: "b"},
		{ID: "c", Actor: "c", EventDesc: "c", EventType: "c", ServiceName: "c", Time: "c"},
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
			name: "find where to start if a was last acked event",
			want: -1,
			args: args{
				lastAckedEvent: &AivenEvent{ID: "a", Actor: "a", EventDesc: "a", EventType: "a", ServiceName: "a", Time: "a"},
				events:         events,
			},
		},
		{
			name: "find where to start if b was last acked event",
			want: 1,
			args: args{
				lastAckedEvent: &AivenEvent{ID: "b", Actor: "b", EventDesc: "b", EventType: "b", ServiceName: "b", Time: "b"},
				events:         events,
			},
		},
		{
			name: "find where to start if c was last acked event",
			want: 2,
			args: args{
				lastAckedEvent: &AivenEvent{ID: "c", Actor: "c", EventDesc: "c", EventType: "c", ServiceName: "c", Time: "c"},
				events:         events,
			},
		},
		{
			name: "find where to start if last acked event not in set of events",
			want: 3,
			args: args{
				lastAckedEvent: &AivenEvent{ID: "z", Actor: "z", EventDesc: "z", EventType: "z", ServiceName: "z", Time: "z"},
				events:         events,
			},
		},
		{
			name: "find where to start if last acked event not set",
			want: 3,
			args: args{
				lastAckedEvent: nil,
				events:         events,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindStartIndex(tt.args.events, tt.args.lastAckedEvent); got != tt.want {
				t.Errorf("FindStartIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
