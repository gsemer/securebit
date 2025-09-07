package presentation

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"securebit/domain"
	"securebit/utils"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Get(username string) (domain.User, error) {
	args := m.Called(username)
	if user, ok := args.Get(0).(*domain.User); ok {
		return *user, args.Error(1)
	}
	return domain.User{}, args.Error(1)
}

func (m *MockUserService) Create(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserService) Delete(user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(1)
}

func TestLogin(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	err := redisClient.FlushDB(context.Background()).Err()
	assert.NoError(t, err, "should flush Redis before test")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	mockUser := &domain.User{Username: "username", HashedPassword: string(hashedPassword)}

	mockUserService := new(MockUserService)
	mockUserService.On("Get", "username").Return(mockUser, nil)

	utils.SignedToken = func(userID string, duration time.Time, secretKey string) (string, error) {
		return "token", nil
	}

	mockRedisClient := new(MockRedisClient)
	mockRedisClient.On("Set", mock.Anything, mock.AnythingOfType("string"), "token", 24*time.Hour).Return(nil)

	loginRequest := &domain.UserRequest{
		Username: "username",
		Password: "password",
	}
	loginRequestToJSON, err := json.Marshal(loginRequest)
	assert.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, "api/v1/login", bytes.NewBuffer(loginRequestToJSON))
	assert.NoError(t, err)

	response := httptest.NewRecorder()

	handler := AuthHandler{us: mockUserService, redisClient: redisClient}
	handler.Login(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, "user logged in", string(body))
}
