package models

import (
	"database/sql/driver"
	"reflect"
	"testing"
	"time"
)

func TestWeekdays_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		w       []byte
		want    Weekdays
		wantErr bool
	}{
		{
			name: "unordered days",
			w:    []byte(`"Sun Mon"`),
			want: Weekdays{
				time.Sunday: struct{}{},
				time.Monday: struct{}{},
			},
			wantErr: false,
		},
		{
			name:    "empty string",
			w:       []byte(`""`),
			want:    Weekdays{},
			wantErr: false,
		},
		{
			name:    "nil value",
			w:       nil,
			want:    Weekdays{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Weekdays
			err := got.UnmarshalJSON(tt.w)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestWeekdays_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		w       Weekdays
		want    []byte
		wantErr bool
	}{
		{
			name: "list of days always ordered, sunday is last",
			w: Weekdays{
				time.Sunday: struct{}{},
				time.Monday: struct{}{},
				time.Friday: struct{}{},
			},
			want:    []byte(`"Mon Fri Sun"`),
			wantErr: false,
		},
		{
			name:    "no days",
			w:       Weekdays{},
			want:    []byte(`""`),
			wantErr: false,
		},
		{
			name:    "not initialized",
			w:       nil,
			want:    []byte(`""`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestWeekdays_Scan(t *testing.T) {
	tests := []struct {
		name    string
		w       int64
		want    Weekdays
		wantErr bool
	}{
		{
			name: "sunday first",
			want: Weekdays{
				time.Sunday: struct{}{},
			},
			w:       int64(0b00000001),
			wantErr: false,
		},
		{
			name: "saturday last",
			want: Weekdays{
				time.Saturday: struct{}{},
			},
			w:       int64(0b01000000),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Weekdays
			err := got.Scan(tt.w)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scan() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestWeekdays_Value(t *testing.T) {
	tests := []struct {
		name    string
		w       Weekdays
		want    driver.Value
		wantErr bool
	}{
		{
			name: "sunday first",
			w: Weekdays{
				time.Sunday: struct{}{},
			},
			want:    int64(0b00000001),
			wantErr: false,
		},
		{
			name: "saturday last",
			w: Weekdays{
				time.Saturday: struct{}{},
			},
			want:    int64(0b01000000),
			wantErr: false,
		},
		{
			name:    "not initialized",
			w:       nil,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
		})
	}
}
