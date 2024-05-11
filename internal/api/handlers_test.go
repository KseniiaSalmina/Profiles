package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/KseniiaSalmina/Profiles/internal/api/models"
	"github.com/KseniiaSalmina/Profiles/internal/config"
	"github.com/KseniiaSalmina/Profiles/internal/database"
	"github.com/KseniiaSalmina/Profiles/internal/logger"
	"github.com/KseniiaSalmina/Profiles/internal/service"
)

var serverCfg = config.Server{
	Listen:       ":8080",
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 5 * time.Second,
	IdleTimeout:  30 * time.Second,
}

var serviceCfg = config.Service{
	Salt:          "",
	AdminUsername: "username",
	AdminPassword: "password",
	AdminEmail:    "test@email.com",
}

var loggercfg = config.Logger{
	LogLevel: "debug",
}

var testUsers = []database.User{
	{
		ID:       "a7073076-8602-4b95-8c19-0cd24aa511c9",
		Email:    "test@email.com",
		Username: "testUser",
		PassHash: "super hash",
		Admin:    false},
	{
		ID:       "28ceb514-ea0d-4ca7-a330-9763b8bd7fc4",
		Email:    "test2@email.com",
		Username: "testUser2",
		PassHash: "super hash2",
		Admin:    false},
	{
		ID:       "db783cb2-8037-4b75-8c01-ab9065e568e3",
		Email:    "test3@email.com",
		Username: "testUser3",
		PassHash: "$2a$10$KIsJbN5.Jvtg1rvB4umGu.mbZGfN6..kOyPcEJ4u/GLNU.thjfeyO",
		Admin:    false},
}

func prepareServer() *Server {
	db := database.NewDatabase()
	service, err := service.NewService(serviceCfg, db)
	if err != nil {
		log.Fatal("failed to prepare service")
	}

	prepareDB(db)

	logger, err := logger.NewLogger(loggercfg)
	if err != nil {
		log.Fatal("failed to prepare logger")
	}

	return NewServer(serverCfg, service, logger)
}

func prepareDB(db *database.Database) {
	for _, user := range testUsers {
		err := db.AddUser(user)
		if err != nil {
			log.Fatalf("failed to add user: %v, %s", user.ID, err.Error())
		}
	}
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
				{ID: "28ceb514-ea0d-4ca7-a330-9763b8bd7fc4",
					Email:    "test2@email.com",
					Username: "testUser2"},
				{ID: "db783cb2-8037-4b75-8c01-ab9065e568e3",
					Email:    "test3@email.com",
					Username: "testUser3"},
			},
			PageNo:      2,
			Limit:       2,
			PagesAmount: 2,
		}}},
		{name: "unauthorized case", args: args{w: httptest.NewRecorder(), r: requests[1]}, want: res{statusCode: http.StatusUnauthorized}},
	}

	server := prepareServer()

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			server.httpServer.Handler.ServeHTTP(tt.args.w, tt.args.r)
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
	requests := make([]*http.Request, 0, 2)

	//standard case
	req, err := http.NewRequest("GET", "/user?limit=2&page=2", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("username", "password")
	requests = append(requests, req)

	//unauthorized case
	req2, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		log.Fatal(err)
	}
	requests = append(requests, req2)

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

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			server.httpServer.Handler.ServeHTTP(tt.args.w, tt.args.r)
			assert.Equal(t1, tt.want.statusCode, tt.args.w.Code)
		})
	}
}

func postUserPrepareReq() []*http.Request {
	requests := make([]*http.Request, 0, 3)

	correctUser1 := models.UserAdd{
		Email:    "test@gmail.com",
		Username: "newUser1",
		Password: "super",
		Admin:    false,
	}

	correctUser2 := models.UserAdd{
		Email:    "test3@gmail.com",
		Username: "newUser2",
		Password: "super3",
		Admin:    false,
	}

	incorrectUser := models.UserAdd{
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
	req.SetBasicAuth("username", "password")
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
	req2.SetBasicAuth("testUser3", "password")
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
	req3.SetBasicAuth("username", "password")
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
		{name: "standard case", args: args{w: httptest.NewRecorder(), r: requests[0]}, want: res{statusCode: http.StatusOK, user: models.UserResponse{
			ID:       "28ceb514-ea0d-4ca7-a330-9763b8bd7fc4",
			Email:    "test2@email.com",
			Username: "testUser2",
		}}},
		{name: "no user id", args: args{w: httptest.NewRecorder(), r: requests[1]}, want: res{statusCode: http.StatusBadRequest}},
		{name: "user with this id does not exist", args: args{w: httptest.NewRecorder(), r: requests[2]}, want: res{statusCode: http.StatusBadRequest}},
	}

	server := prepareServer()

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			server.httpServer.Handler.ServeHTTP(tt.args.w, tt.args.r)
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
	req, err := http.NewRequest("GET", "/user/28ceb514-ea0d-4ca7-a330-9763b8bd7fc4", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("username", "password")
	requests = append(requests, req)

	//incorrect path
	req2, err := http.NewRequest("GET", "/user/ggthyh", nil)
	if err != nil {
		log.Fatal(err)
	}
	req2.SetBasicAuth("username", "password")
	requests = append(requests, req2)

	//user with id does not exist
	req3, err := http.NewRequest("GET", "/user/472b62a5-b5a3-461b-83ef-158839bfe79f", nil)
	if err != nil {
		log.Fatal(err)
	}
	req3.SetBasicAuth("username", "password")
	requests = append(requests, req3)

	return requests
}

func TestServer_patchUser(t1 *testing.T) {
	requests := patchUserPrepareReq()

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
		{name: "not admin", args: args{w: httptest.NewRecorder(), r: requests[1]}, want: res{statusCode: http.StatusForbidden}},
		{name: "incorrect id", args: args{w: httptest.NewRecorder(), r: requests[2]}, want: res{statusCode: http.StatusBadRequest}},
		{name: "user with this id does not exist", args: args{w: httptest.NewRecorder(), r: requests[3]}, want: res{statusCode: http.StatusBadRequest}},
	}

	server := prepareServer()

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			server.httpServer.Handler.ServeHTTP(tt.args.w, tt.args.r)
			assert.Equal(t1, tt.want.statusCode, tt.args.w.Code)
		})
	}
}

func patchUserPrepareReq() []*http.Request {
	requests := make([]*http.Request, 0, 4)

	user1 := models.UserAdd{
		Email:    "test@gmail.com",
		Username: "updatedUser1",
		Password: "super",
		Admin:    false,
	}

	user2 := models.UserAdd{
		Email:    "test22@gmail.com",
		Username: "updatedUser2",
		Password: "super2",
		Admin:    false,
	}

	user3 := models.UserAdd{
		Email:    "test3@gmail.com",
		Username: "updatedUser3",
		Password: "super3",
		Admin:    false,
	}

	user4 := models.UserAdd{
		Email:    "test4@gmail.com",
		Username: "updatedUser4",
		Password: "super4",
		Admin:    false,
	}

	//standard case
	body, err := json.Marshal(user1)
	if err != nil {
		log.Fatal("can not marshal correct user 1")
	}
	req, err := http.NewRequest("PATCH", "/user/28ceb514-ea0d-4ca7-a330-9763b8bd7fc4", strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("username", "password")
	requests = append(requests, req)

	//not admin
	body, err = json.Marshal(user2)
	if err != nil {
		log.Fatal("can not marshal correct user 1")
	}
	req2, err := http.NewRequest("PATCH", "/user/28ceb514-ea0d-4ca7-a330-9763b8bd7fc4", strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	req2.SetBasicAuth("testUser3", "password")
	requests = append(requests, req2)

	//incorrect id
	body, err = json.Marshal(user3)
	if err != nil {
		log.Fatal("can not marshal correct user 1")
	}
	req3, err := http.NewRequest("PATCH", "/user/1000", strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	req3.SetBasicAuth("username", "password")
	requests = append(requests, req3)

	//change not existing user
	body, err = json.Marshal(user4)
	if err != nil {
		log.Fatal("can not marshal correct user 1")
	}
	req4, err := http.NewRequest("PATCH", "/user/34775464-a73b-4445-8866-1e6061c3b70b", strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	req4.SetBasicAuth("username", "password")
	requests = append(requests, req4)

	return requests
}

func TestServer_deleteUser(t1 *testing.T) {
	requests := deleteUserPrepareReq()

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
		{name: "not admin", args: args{w: httptest.NewRecorder(), r: requests[1]}, want: res{statusCode: http.StatusForbidden}},
		{name: "incorrect id", args: args{w: httptest.NewRecorder(), r: requests[2]}, want: res{statusCode: http.StatusBadRequest}},
		{name: "user with this id does not exist", args: args{w: httptest.NewRecorder(), r: requests[3]}, want: res{statusCode: http.StatusBadRequest}},
	}

	server := prepareServer()

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			server.httpServer.Handler.ServeHTTP(tt.args.w, tt.args.r)
			assert.Equal(t1, tt.want.statusCode, tt.args.w.Code)
		})
	}
}

func deleteUserPrepareReq() []*http.Request {
	requests := make([]*http.Request, 0, 4)

	//standard case
	req, err := http.NewRequest("DELETE", "/user/28ceb514-ea0d-4ca7-a330-9763b8bd7fc4", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("username", "password")
	requests = append(requests, req)

	//not admin
	req2, err := http.NewRequest("DELETE", "/user/a7073076-8602-4b95-8c19-0cd24aa511c9", nil)
	if err != nil {
		log.Fatal(err)
	}
	req2.SetBasicAuth("testUser3", "password")
	requests = append(requests, req2)

	//incorrect id
	req3, err := http.NewRequest("DELETE", "/user/1000", nil)
	if err != nil {
		log.Fatal(err)
	}
	req3.SetBasicAuth("username", "password")
	requests = append(requests, req3)

	//change not existing user
	req4, err := http.NewRequest("DELETE", "/user/34775464-a73b-4445-8866-1e6061c3b70b", nil)
	if err != nil {
		log.Fatal(err)
	}
	req4.SetBasicAuth("username", "password")
	requests = append(requests, req4)

	return requests
}
