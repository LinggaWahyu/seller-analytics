package repository

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/pkg/db/yugabyte"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/buyer/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db       *sql.DB
	mock     sqlmock.Sqlmock
	setuperr error
	gormdb   *gorm.DB
)

func setupTest() {
	db, mock, setuperr = sqlmock.New()
	if setuperr != nil {
		os.Exit(1)
	}

	gormdb, setuperr = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if setuperr != nil {
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	setupTest()
	m.Run()
}

func Test_repository_Get(t *testing.T) {
	tests := []struct {
		name    string
		buyer   domain.Buyer
		want    *domain.Buyer
		wantErr bool
		mock    func()
	}{
		{
			name: "error",
			buyer: domain.Buyer{
				Model: yugabyte.Model{
					ID: 1,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "buyers" WHERE "buyers"."id" = $1 AND "buyers"."deleted_at" IS NULL AND "buyers"."id" = $2 ORDER BY "buyers"."id" LIMIT 1`)).
					WithArgs(int64(1), int64(1)).WillReturnError(errors.New("mock error"))
			},
		},
		{
			name: "success",
			buyer: domain.Buyer{
				Model: yugabyte.Model{
					ID: 1,
				},
			},
			want: &domain.Buyer{
				Model: yugabyte.Model{
					ID: 1,
				},
				Username: "testuser",
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "buyers" WHERE "buyers"."id" = $1 AND "buyers"."deleted_at" IS NULL AND "buyers"."id" = $2 ORDER BY "buyers"."id" LIMIT 1`)).
					WithArgs(int64(1), int64(1)).WillReturnRows(sqlmock.NewRows([]string{"username"}).
					AddRow("testuser"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			sut := NewBuyerRepository(gormdb)
			res, err := sut.Get(context.TODO(), tt.buyer)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, res)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_GetByUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     *domain.Buyer
		wantErr  bool
		mock     func()
	}{
		{
			name:     "success",
			username: "testuser",
			want: &domain.Buyer{
				Username: "testuser",
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "buyers" WHERE username = $1 AND "buyers"."deleted_at" IS NULL ORDER BY "buyers"."id" LIMIT 1`)).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"username"}).
						AddRow("testuser"))
			},
		},
		{
			name:     "error",
			username: "testuser",
			want:     nil,
			wantErr:  true,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "buyers" WHERE username = $1 AND "buyers"."deleted_at" IS NULL ORDER BY "buyers"."id" LIMIT 1`)).
					WithArgs("testuser").WillReturnError(errors.New("mock error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			sut := NewBuyerRepository(gormdb)
			res, err := sut.GetByUsername(context.TODO(), tt.username)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, res)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_Create(t *testing.T) {
	tests := []struct {
		name    string
		buyer   domain.Buyer
		want    *domain.Buyer
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			buyer: domain.Buyer{
				Username: "testuser",
			},
			want: &domain.Buyer{
				Username: "testuser",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "buyers" ("created_at","updated_at","deleted_at","username") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser").
					WillReturnRows(sqlmock.NewRows([]string{"username"}).
						AddRow("testuser"))
				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			buyer: domain.Buyer{
				Username: "testuser",
			},
			want:    nil,
			wantErr: true,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "buyers" ("created_at","updated_at","deleted_at","username") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser").
					WillReturnError(errors.New("mock error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			sut := NewBuyerRepository(gormdb)
			res, err := sut.Create(context.TODO(), tt.buyer)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if tt.want == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tt.want.Username, res.Username)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
