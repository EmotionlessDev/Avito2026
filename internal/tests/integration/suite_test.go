package integration

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/tests"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/tests/fixtures"
)

type IntegrationSuite struct {
	suite.Suite
	pg  *tests.PostgresContainer
	db  *sql.DB
	ctx context.Context
}

func (s *IntegrationSuite) SetupSuite() {
	s.T().Log("Setting up integration test suite...")

	s.ctx = context.Background()
	s.pg = tests.SetupPostgres(s.T())
	s.db = s.pg.DB

	migrationsPath := filepath.Join("..", "..", "..", "migrations")
	tests.RunMigrations(s.T(), s.db, migrationsPath)

	s.T().Log("Integration test suite setup complete")
}

func (s *IntegrationSuite) TearDownSuite() {
	s.T().Log("Tearing down integration test suite...")
	tests.TeardownPostgres(s.T(), s.pg)
	s.T().Log("Integration test suite teardown complete")
}

func (s *IntegrationSuite) SetupTest() {
	s.T().Log("Cleaning up tables before test...")
	tests.TruncateTables(s.T(), s.db)
}

func (s *IntegrationSuite) TearDownTest() {
}

func (s *IntegrationSuite) CreateUser(fixture fixtures.UserFixture) string {
	return fixtures.CreateUser(s.T(), s.db, fixture)
}

func (s *IntegrationSuite) CreateRoom(fixture fixtures.RoomFixture) string {
	return fixtures.CreateRoom(s.T(), s.db, fixture)
}

func (s *IntegrationSuite) CreateSlot(fixture fixtures.SlotFixture) string {
	return fixtures.CreateSlot(s.T(), s.db, fixture)
}

func (s *IntegrationSuite) CreateBooking(fixture fixtures.BookingFixture) string {
	return fixtures.CreateBooking(s.T(), s.db, fixture)
}

func (s *IntegrationSuite) GetTestUser(role string) string {
	return fixtures.GetTestUser(s.T(), s.db, role)
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
