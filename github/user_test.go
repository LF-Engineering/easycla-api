package github

import (
	"testing"
)

func TestGetUsernameFromID(t *testing.T) {
	Init()
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid id",
			args: args{
				id: 1476068,
			},
			want:    "prasannamahajan",
			wantErr: false,
		},
		{
			name: "invalid id",
			args: args{
				id: 1476068000000,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUsernameFromID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsernameFromID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUsernameFromID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIDFromUsername(t *testing.T) {
	Init()
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "valid username",
			args: args{
				username: "prasannamahajan",
			},
			want:    1476068,
			wantErr: false,
		},
		{
			name: "invalid username",
			args: args{
				username: "gaZMZR8RUeXVQHDv4XpI",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIDFromUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIDFromUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetIDFromUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
