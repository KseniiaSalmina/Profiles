package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/KseniiaSalmina/Profiles/internal/api/models"
	"github.com/KseniiaSalmina/Profiles/internal/config"
	"github.com/KseniiaSalmina/Profiles/internal/database"
	"github.com/KseniiaSalmina/Profiles/internal/formatter"
	"github.com/stretchr/testify/assert"
)

var dbCfg = config.Database{
	AdminUsername: "username",
	AdminPassword: "password",
	AdminEmail:    "test@email",
}

var serverCfg = config.Server{
	Listen:       ":8080",
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 5 * time.Second,
	IdleTimeout:  30 * time.Second,
}

var formatterCfg = config.Formatter{
	Salt: "",
}

var testUsers = []database.User{
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
		PassHash: "$2a$10$KIsJbN5.Jvtg1rvB4umGu.mbZGfN6..kOyPcEJ4u/GLNU.thjfeyO",
		Admin:    false},
}

func prepareDB() *database.Database {
	db, err := database.NewDatabase(dbCfg, formatterCfg.Salt)
	if err != nil {
		log.Fatalf("failed to create db: %s", err.Error())
	}
	for _, user := range testUsers {
		err = db.AddUser(user)
		if err != nil {
			log.Fatalf("failed to add user: %v, %s", user.ID, err.Error())
		}
	}

	return db
}

func prepareServer() *Server {
	db := prepareDB()
	formatter := formatter.NewFormatter(formatterCfg, db)
	return NewServer(serverCfg, formatter)
}

func TestServer_getAllUsers(t1 *testing.T) {
	requests := getAllUsersPrepareReq()

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type res struct {
		statusCode int
		users      models.PageUsers
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "standard case", args: args{w: httptest.NewRecorder(), r: requests[0]}, want: res{statusCode: http.StatusOK, users: models.PageUsers{
			Users: []models.UserResponse{
				{ID: "2",
					Email:    "test2@email.com",
					Username: "testUser2"},
				{ID: "3",
					Email:    "test3@email.com",
					Username: "testUser3"},
			},
			PageNo:      2,
			Limit:       2,
			PagesAmount: 2,
		}}},
		{name: "unauthorized case", args: args{w: httptest.NewRecorder(), r: requests[1]}, want: res{statusCode: http.StatusUnauthorized}},
		{name: "incorrect encoded auth string", args: args{w: httptest.NewRecorder(), r: requests[2]}, want: res{statusCode: http.StatusBadRequest}},
	}

	server := prepareServer()
	handler := http.HandlerFunc(server.getAllUsers)

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			handler.ServeHTTP(tt.args.w, tt.args.r)
			assert.Equal(t1, tt.want.statusCode, tt.args.w.Code)

			if tt.want.statusCode == 200 {
				var res models.PageUsers
				err := json.NewDecoder(tt.args.w.Body).Decode(&res)
				if err != nil {
					t1.Fatalf("can not decode: %v", err.Error())
				}
				assert.Equal(t1, tt.want.users, res)
			}

		})
	}
}

func getAllUsersPrepareReq() []*http.Request {
	requests := make([]*http.Request, 0, 3)

	//standard case
	req, err := http.NewRequest("GET", "/user?limit=2&page=2", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	requests = append(requests, req)

	//unauthorized case
	req2, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		log.Fatal(err)
	}
	requests = append(requests, req2)

	//incorrect encoded auth string
	req3, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		log.Fatal(err)
	}
	req3.Header.Set("Authorization", "Basic incorrect string")
	requests = append(requests, req3)

	return requests
}

func TestServer_postUser(t1 *testing.T) {
	requests := postUserPrepareReq()

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type res struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "standard case", args: args{w: httptest.NewRecorder(), r: requests[0]}, want: res{statusCode: http.StatusOK}},
		{name: "username taken", args: args{w: httptest.NewRecorder(), r: requests[2]}, want: res{statusCode: http.StatusBadRequest}},
		{name: "not admin", args: args{w: httptest.NewRecorder(), r: requests[1]}, want: res{statusCode: http.StatusForbidden}},
	}

	server := prepareServer()
	handler := http.HandlerFunc(server.postUser)

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			handler.ServeHTTP(tt.args.w, tt.args.r)
			assert.Equal(t1, tt.want.statusCode, tt.args.w.Code)
		})
	}
}

func postUserPrepareReq() []*http.Request {
	requests := make([]*http.Request, 0, 3)

	correctUser1 := models.UserRequest{
		Email:    "test@gmail.com",
		Username: "newUser1",
		Password: "super",
		Admin:    false,
	}

	correctUser2 := models.UserRequest{
		Email:    "test3@gmail.com",
		Username: "newUser2",
		Password: "super3",
		Admin:    false,
	}

	incorrectUser := models.UserRequest{
		Email:    "test2@gmail.com",
		Username: "testUser",
		Password: "super2",
		Admin:    false,
	}

	//standard case: admin
	body, err := json.Marshal(correctUser1)
	if err != nil {
		log.Fatal("can not marshal correct user 1")
	}
	req, err := http.NewRequest("POST", "/user", strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	requests = append(requests, req)

	//not admin
	body, err = json.Marshal(correctUser2)
	if err != nil {
		log.Fatal("can not marshal correct user 1")
	}
	req2, err := http.NewRequest("POST", "/user", strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	req2.Header.Set("Authorization", "Basic dGVzdFVzZXIzOnBhc3N3b3Jk")
	requests = append(requests, req2)

	//username taken
	body, err = json.Marshal(incorrectUser)
	if err != nil {
		log.Fatal("can not marshal correct user 1")
	}
	req3, err := http.NewRequest("POST", "/user", strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	req3.Header.Set("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	requests = append(requests, req3)

	return requests
}

func TestServer_getUser(t1 *testing.T) {
	requests := getUserPrepareReq()

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type res struct {
		statusCode int
		user       models.UserResponse
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{name: "standard case", args: args{w: httptest.NewRecorder(), r: requests[0]}, want: res{statusCode: http.StatusOK}},
		{name: "no user id", args: args{w: httptest.NewRecorder(), r: requests[1]}, want: res{statusCode: http.StatusBadRequest}},
		{name: "user with this id does not exist", args: args{w: httptest.NewRecorder(), r: requests[2]}, want: res{statusCode: http.StatusBadRequest}},
	}

	server := prepareServer()
	handler := http.HandlerFunc(server.getUser)

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			handler.ServeHTTP(tt.args.w, tt.args.r)
			assert.Equal(t1, tt.want.statusCode, tt.args.w.Code)

			if tt.want.statusCode == 200 {
				var res models.UserResponse
				err := json.NewDecoder(tt.args.w.Body).Decode(&res)
				if err != nil {
					t1.Fatalf("can not decode: %v", err.Error())
				}
				assert.Equal(t1, tt.want.user, res)
			}
		})
	}
}

func getUserPrepareReq() []*http.Request {
	requests := make([]*http.Request, 0, 3)

	//standard case
	req, err := http.NewRequest("GET", "/user/2", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	requests = append(requests, req)

	//incorrect path
	req2, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		log.Fatal(err)
	}
	req2.Header.Set("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	requests = append(requests, req2)

	//user with id does not exist
	req3, err := http.NewRequest("GET", "/user/1000", nil)
	if err != nil {
		log.Fatal(err)
	}
	req3.Header.Set("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	requests = append(requests, req3)

	return requests
}
