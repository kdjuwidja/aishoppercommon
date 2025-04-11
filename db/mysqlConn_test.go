package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestModel is a simple model for testing
type TestModel struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"size:255;not null"`
}

// TestUser is a model for testing user operations
type TestUser struct {
	gorm.Model
	ID       string `json:"id" gorm:"type:varchar(32);primaryKey"`
	Email    string `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
	IsActive bool   `json:"is_active" gorm:"type:tinyint(1);not null;default:1"`
}

// TestAPIClient is a model for testing API client operations
type TestAPIClient struct {
	gorm.Model
	ID          string `json:"id" gorm:"type:varchar(45);primaryKey"`
	Secret      string `json:"secret" gorm:"type:varchar(45);not null"`
	Domain      string `json:"domain" gorm:"type:varchar(255);not null"`
	IsPublic    bool   `json:"is_public" gorm:"type:tinyint(1);not null;default:0"`
	Description string `json:"description" gorm:"type:varchar(255);"`
}

// setupTestPool creates and configures a test connection pool
func setupTestPool(t *testing.T) *MySQLConnectionPool {
	pool := &MySQLConnectionPool{}
	pool.Configure("ai_shopper_dev", "password", "localhost", "4306", "test_db", 10, 5, []interface{}{&TestModel{}})
	return pool
}

func TestMySQLConnectionPool_Configure(t *testing.T) {
	pool := &MySQLConnectionPool{}

	// Test valid configuration
	pool.Configure("ai_shopper_dev", "password", "localhost", "4306", "test_db", 10, 5, []interface{}{&TestModel{}})

	assert.Equal(t, "ai_shopper_dev", pool.User)
	assert.Equal(t, "password", pool.Password)
	assert.Equal(t, "localhost", pool.Host)
	assert.Equal(t, "4306", pool.Port)
	assert.Equal(t, "test_db", pool.DBName)
	assert.Equal(t, 10, pool.MaxOpenConns)
	assert.Equal(t, 5, pool.MaxIdleConns)
	assert.Len(t, pool.models, 1)
}

func TestMySQLConnectionPool_Initialize_Validation(t *testing.T) {
	// Test missing configuration
	pool := &MySQLConnectionPool{}
	err := pool.Initialize()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required configuration parameters")

	// Test with empty values
	pool = &MySQLConnectionPool{}
	pool.Configure("", "", "", "", "", 0, 0, nil)
	err = pool.Initialize()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required configuration parameters")
}

func TestMySQLConnectionPool_Integration(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	pool := &MySQLConnectionPool{}

	// Configure with real database credentials and actual models
	pool.Configure(
		"ai_shopper_dev",
		"password",
		"localhost",
		"4306",
		"test_db",
		10,
		5,
		[]interface{}{
			&TestUser{},
			&TestAPIClient{},
		},
	)

	// Test initialization
	err := pool.Initialize()
	assert.NoError(t, err)

	// Test GetDB
	dbConn := pool.GetDB()
	assert.NotNil(t, dbConn)

	// Test AutoMigrate
	err = pool.AutoMigrate()
	assert.NoError(t, err)

	// Test database operations with User model
	user := &TestUser{
		ID:       "test-user-1",
		Email:    "test@example.com",
		Password: "hashedpassword",
		IsActive: true,
	}
	result := dbConn.Create(user)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, user.ID)

	// Verify user was created
	var retrievedUser TestUser
	result = dbConn.First(&retrievedUser, "id = ?", user.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "test@example.com", retrievedUser.Email)
	assert.Equal(t, "hashedpassword", retrievedUser.Password)
	assert.True(t, retrievedUser.IsActive)

	// Test database operations with APIClient model
	apiClient := &TestAPIClient{
		ID:          "test-client-1",
		Secret:      "test-secret",
		Domain:      "http://localhost:3000",
		IsPublic:    true,
		Description: "Test API Client",
	}
	result = dbConn.Create(apiClient)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, apiClient.ID)

	// Verify API client was created
	var retrievedClient TestAPIClient
	result = dbConn.First(&retrievedClient, "id = ?", apiClient.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "test-client-1", retrievedClient.ID)
	assert.Equal(t, "test-secret", retrievedClient.Secret)
	assert.Equal(t, "http://localhost:3000", retrievedClient.Domain)
	assert.True(t, retrievedClient.IsPublic)
	assert.Equal(t, "Test API Client", retrievedClient.Description)

	// Test DropTables
	err = pool.DropTables()
	assert.NoError(t, err)

	// Verify tables were dropped
	var userCount int64
	result = dbConn.Model(&TestUser{}).Count(&userCount)
	assert.Error(t, result.Error) // Should error because table doesn't exist

	var clientCount int64
	result = dbConn.Model(&TestAPIClient{}).Count(&clientCount)
	assert.Error(t, result.Error) // Should error because table doesn't exist

	// Test Close
	err = pool.Close()
	assert.NoError(t, err)
}

func TestMySQLConnectionPool_GetDB(t *testing.T) {
	// Test before initialization
	pool := &MySQLConnectionPool{}
	assert.Nil(t, pool.GetDB())

	// Test after initialization
	pool = setupTestPool(t)
	err := pool.Initialize()
	assert.NoError(t, err)

	db := pool.GetDB()
	assert.NotNil(t, db)

	// Clean up
	pool.Close()
}

func TestMySQLConnectionPool_Close(t *testing.T) {
	// Test closing uninitialized pool
	pool := &MySQLConnectionPool{}
	err := pool.Close()
	assert.NoError(t, err)

	// Test closing initialized pool
	pool = setupTestPool(t)
	err = pool.Initialize()
	assert.NoError(t, err)

	err = pool.Close()
	assert.NoError(t, err)
}

func TestMySQLConnectionPool_AutoMigrate(t *testing.T) {
	// Test auto-migrate before initialization
	pool := &MySQLConnectionPool{}
	err := pool.AutoMigrate()
	assert.Error(t, err)

	// Test auto-migrate after initialization
	pool = setupTestPool(t)
	err = pool.Initialize()
	assert.NoError(t, err)

	err = pool.AutoMigrate()
	assert.NoError(t, err)

	// Clean up
	pool.DropTables()
	pool.Close()
}

func TestMySQLConnectionPool_DropTables(t *testing.T) {
	// Test dropping tables before initialization
	pool := &MySQLConnectionPool{}
	err := pool.DropTables()
	assert.Error(t, err)

	// Test dropping tables after initialization and migration
	pool = setupTestPool(t)
	err = pool.Initialize()
	assert.NoError(t, err)

	err = pool.AutoMigrate()
	assert.NoError(t, err)

	err = pool.DropTables()
	assert.NoError(t, err)

	// Clean up
	pool.Close()
}
