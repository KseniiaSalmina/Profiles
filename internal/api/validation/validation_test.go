package validation

import (
	"testing"

	"github.com/KseniiaSalmina/Profiles/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestAuthString(t1 *testing.T) {
	type args struct {
		authString string
	}
	type res struct {
		username string
		password string
		wantErr  bool
		err      string
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "standard case", args: args{authString: "Basic dXNlcm5hbWU6cGFzc3dvcmQ="}, want: res{username: "username", password: "password", wantErr: false, err: ""}},
		{name: "does not have prefix", args: args{authString: "dXNlcm5hbWU6cGFzc3dvcmQ="}, want: res{username: "", password: "", wantErr: true, err: ErrIncorrectAuth.Error()}},
		{name: "not encoded auth string", args: args{authString: "Basic username:password"}, want: res{username: "", password: "", wantErr: true, err: "failed to decode authorization string: illegal base64 data at input byte 8"}},
		{name: "have \":\" in username", args: args{authString: "Basic dXNlcjpuYW1lOnBhc3N3b3Jk"}, want: res{username: "", password: "", wantErr: true, err: ErrIncorrectAuth.Error()}},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			username, password, err := AuthString(tt.args.authString)
			if !tt.want.wantErr {
				assert.Equal(t1, tt.want.username, username)
				assert.Equal(t1, tt.want.password, password)
			} else {
				assert.Equal(t1, tt.want.err, err.Error())
			}
		})
	}
}

func TestUser(t1 *testing.T) {
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
		}}, want: res{wantErr: true, err: ErrIncorrectUserData}},
		{name: "incorrect username", args: args{username: "newUsername", password: "password", user: database.User{
			Username: "username",
			PassHash: "$2a$10$KIsJbN5.Jvtg1rvB4umGu.mbZGfN6..kOyPcEJ4u/GLNU.thjfeyO",
		}}, want: res{wantErr: true, err: ErrIncorrectUserData}},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			err := User(tt.args.username, tt.args.password, tt.args.user)
			if !tt.want.wantErr {
				assert.NoError(t1, err)
			}
			assert.Equal(t1, tt.want.err, err)
		})
	}
}
