package validation

import (
	"testing"

	"github.com/KseniiaSalmina/Profiles/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestAuthData(t1 *testing.T) {
	type args struct {
		username string
		password string
		user     database.User
	}
	type res struct {
		wantErr bool
		err     error
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "standard case", args: args{username: "username", password: "password", user: database.User{
			Username: "username",
			PassHash: "$2a$10$KIsJbN5.Jvtg1rvB4umGu.mbZGfN6..kOyPcEJ4u/GLNU.thjfeyO",
		}}, want: res{wantErr: false, err: nil}},
		{name: "incorrect password", args: args{username: "username", password: "newPassword", user: database.User{
			Username: "username",
			PassHash: "$2a$10$KIsJbN5.Jvtg1rvB4umGu.mbZGfN6..kOyPcEJ4u/GLNU.thjfeyO",
		}}, want: res{wantErr: true, err: ErrIncorrectAuthData}},
		{name: "incorrect username", args: args{username: "newUsername", password: "password", user: database.User{
			Username: "username",
			PassHash: "$2a$10$KIsJbN5.Jvtg1rvB4umGu.mbZGfN6..kOyPcEJ4u/GLNU.thjfeyO",
		}}, want: res{wantErr: true, err: ErrIncorrectAuthData}},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			err := Auth(tt.args.username, tt.args.password, tt.args.user)
			if !tt.want.wantErr {
				assert.NoError(t1, err)
			}
			assert.Equal(t1, tt.want.err, err)
		})
	}
}
