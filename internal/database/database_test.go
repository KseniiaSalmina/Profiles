package database

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testUsers = []User{
	{
		ID:       "1",
		Email:    "test@email.com",
		Username: "testUser",
		PassHash: "super hash",
		Admin:    false},
	{
		ID:       "2",
		Email:    "test2@email.com",
		Username: "testUser2",
		PassHash: "super hash2",
		Admin:    false},
	{
		ID:       "3",
		Email:    "test3@email.com",
		Username: "testUser3",
		PassHash: "super hash3",
		Admin:    false},
}

func prepareDB(isFull bool) *Database {
	db := NewDatabase()

	if isFull {
		for _, user := range testUsers {
			err := db.AddUser(user)
			if err != nil {
				log.Fatalf("failed to add user: %v, %s", user.ID, err.Error())
			}
		}
	}

	return db
}

func TestDatabase_AddUser(t1 *testing.T) {
	type args struct {
		user User
	}
	type res struct {
		wantErr bool
		error   error
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "standard case", args: args{user: User{
			ID:       "1",
			Email:    "test@email.com",
			Username: "testUser",
			PassHash: "super hash",
			Admin:    true,
		}}, want: res{wantErr: false, error: nil}},
		{name: "repeating ID", args: args{user: User{
			ID:       "1",
			Email:    "test2@email.com",
			Username: "test2User",
			PassHash: "super hash2",
			Admin:    true,
		}}, want: res{wantErr: true, error: ErrUserAlreadyExist}},
		{name: "repeating username", args: args{user: User{
			ID:       "3",
			Email:    "test3@email.com",
			Username: "testUser",
			PassHash: "super hash3",
			Admin:    true,
		}}, want: res{wantErr: true, error: ErrNotUniqueUsername}},
	}

	db := prepareDB(false)

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			err := db.AddUser(tt.args.user)
			if !tt.want.wantErr {
				assert.NoError(t1, err)
				assert.Equal(t1, &tt.args.user, db.idIDX[tt.args.user.ID])
				assert.Equal(t1, &tt.args.user, db.usernameIDX[tt.args.user.Username])
				assert.Equal(t1, &tt.args.user, db.users[0]) //TODO: change if add new test cases
			}
			assert.Equal(t1, tt.want.error, err)
		})
	}
}

func TestDatabase_GetAllUsers(t1 *testing.T) {
	type args struct {
		offset int
		limit  int
	}
	type res struct {
		users []User
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "standard case", args: args{offset: 0, limit: 2}, want: res{testUsers[0:2]}},
		{name: "only one user in result", args: args{offset: 1, limit: 1}, want: res{testUsers[1:2]}},
		{name: "offset more than amount of users in db", args: args{offset: 5, limit: 2}, want: res{[]User{}}},
		{name: "offset is equal to amount of users in db", args: args{offset: 4, limit: 2}, want: res{[]User{}}},
		{name: "offset+limit is more than len of slice of users in db", args: args{offset: 0, limit: 5}, want: res{testUsers}},
	}

	db := prepareDB(true)

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			users := db.GetAllUsers(tt.args.offset, tt.args.limit)
			assert.Equal(t1, tt.want.users, users)
		})
	}
}

func TestDatabase_GetUserByID(t1 *testing.T) {
	type args struct {
		userID string
	}
	type res struct {
		user    User
		wantErr bool
		err     error
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "standard case", args: args{userID: "1"}, want: res{user: testUsers[0], wantErr: false, err: nil}},
		{name: "user does not exist", args: args{userID: "10"}, want: res{wantErr: true, err: ErrUserDoesNotExist}},
	}

	db := prepareDB(true)

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			user, err := db.GetUserByID(tt.args.userID)
			if tt.want.wantErr {
				assert.Equal(t1, tt.want.err, err)
			} else {
				assert.NoError(t1, err)
				assert.Equal(t1, tt.want.user, *user)
			}
		})
	}
}

func TestDatabase_GetUserByUsername(t1 *testing.T) {
	type args struct {
		username string
	}
	type res struct {
		user    User
		wantErr bool
		err     error
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "standard case", args: args{username: "testUser"}, want: res{user: testUsers[0], wantErr: false, err: nil}},
		{name: "user does not exist", args: args{username: "superUser2000"}, want: res{wantErr: true, err: ErrUserDoesNotExist}},
	}

	db := prepareDB(true)

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			user, err := db.GetUserByUsername(tt.args.username)
			if tt.want.wantErr {
				assert.Equal(t1, tt.want.err, err)
			} else {
				assert.NoError(t1, err)
				assert.Equal(t1, tt.want.user, *user)
			}
		})
	}
}

func TestDatabase_ChangeUser(t1 *testing.T) {
	type args struct {
		user UserUpdate
	}
	type res struct {
		wantErr bool
		err     error
	}

	email := "newTest@email.com"
	username1, username2, username3 := "testUser", "testUser2000", "testUser2"
	password := "super hash"
	admin := false

	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "standard case", args: args{user: UserUpdate{
			ID:       "1",
			Email:    &email,
			Username: &username1,
			PassHash: &password,
			Admin:    &admin,
		}}, want: res{wantErr: false, err: nil}},
		{name: "change username", args: args{user: UserUpdate{
			ID:       "1",
			Username: &username2,
		}}, want: res{wantErr: false, err: nil}},
		{name: "change username to already taken username", args: args{user: UserUpdate{
			ID:       "1",
			Username: &username3,
		}}, want: res{wantErr: true, err: ErrNotUniqueUsername}},
		{name: "change not existing user", args: args{user: UserUpdate{
			ID:       "1000",
			Email:    &email,
			Username: &username3,
			PassHash: &password,
			Admin:    &admin,
		}}, want: res{wantErr: true, err: ErrUserDoesNotExist}},
	}

	db := prepareDB(true)

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			err := db.ChangeUser(tt.args.user)
			if tt.want.wantErr {
				assert.Equal(t1, tt.want.err, err)
			} else {
				assert.NoError(t1, err)
			}
		})
	}
}

func TestDatabase_DeleteUser(t1 *testing.T) {
	type args struct {
		userID string
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
		{name: "standart case", args: args{userID: "1"}, want: res{wantErr: false, err: nil}},
		{name: "user does not exist", args: args{userID: "25"}, want: res{wantErr: true, err: ErrUserDoesNotExist}},
	}

	db := prepareDB(true)

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			err := db.DeleteUser(tt.args.userID)
			if tt.want.wantErr {
				assert.Equal(t1, tt.want.err, err)
			} else {
				assert.NoError(t1, err)
				_, ok := db.idIDX[tt.args.userID]
				assert.Equal(t1, false, ok)
				_, ok = db.usernameIDX["testUser"] //TODO: change if add new test cases
				assert.Equal(t1, false, ok)
				assert.NotEqual(t1, testUsers[0], *db.users[0]) //TODO: change if add new test cases
			}
		})
	}
}
