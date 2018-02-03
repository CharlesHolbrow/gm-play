package main

import (
	"reflect"
	"testing"
)

func TestNoteGroup_Dedupe(t *testing.T) {
	tests := []struct {
		name       string
		notes      NoteGroup
		wantResult NoteGroup
	}{
		{
			notes:      Group(3, 2, 1, 1, 2, 3),
			wantResult: Group(3, 2, 1),
		},
		{
			notes:      Group(3, 3, 3, 2, 2, 2, 1, 1, 1),
			wantResult: Group(3, 2, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := tt.notes.Dedupe(); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("NoteGroup.Dedupe() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestNoteGroup_Interleave(t *testing.T) {
	type args struct {
		others []NoteGroup
	}
	tests := []struct {
		name       string
		notes      NoteGroup
		args       args
		wantResult NoteGroup
	}{
		{
			name:  "different lengths",
			notes: Group(A, A),
			args: args{
				others: []NoteGroup{
					Group(B, B),
					Group(C, D, E),
				},
			},
			wantResult: Group(A, B, C, A, B, D),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := tt.notes.Interleave(tt.args.others...); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("NoteGroup.Interleave() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
