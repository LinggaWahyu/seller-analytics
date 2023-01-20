package repository

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/statistic/domain"
	"gorm.io/datatypes"
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

func Test_statisticsRepository_GetByDate(t *testing.T) {
	date := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name    string
		date    time.Time
		want    *domain.Statistics
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			date: date,
			want: &domain.Statistics{
				TotalRevenue: 10000,
				DateStr:      "2022-01-01",
				Date:         datatypes.Date(date),
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "statistics" WHERE Date = $1 AND "statistics"."deleted_at" IS NULL ORDER BY "statistics"."id" LIMIT 1`)).
					WithArgs(date).
					WillReturnRows(sqlmock.NewRows([]string{"TotalRevenue", "Date"}).
						AddRow(10000, date))
			},
		},
		{
			name:    "error",
			date:    date,
			want:    nil,
			wantErr: true,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "statistics" WHERE Date = $1 AND "statistics"."deleted_at" IS NULL ORDER BY "statistics"."id" LIMIT 1`)).
					WithArgs(date).WillReturnError(errors.New("mock error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			sr := NewStatisticsRepository(gormdb, nil)
			res, err := sr.GetByDate(context.TODO(), tt.date)
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

func Test_statisticsRepository_Create(t *testing.T) {
	date := datatypes.Date(time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local))
	tests := []struct {
		name      string
		statistic domain.Statistics
		want      *domain.Statistics
		wantErr   bool
		mock      func()
	}{
		{
			name: "success",
			statistic: domain.Statistics{
				TotalRevenue:     10000,
				TotalProductSold: 2,
				CompletedOrder:   1,
				CancelledOrder:   0,
				TotalOrder:       1,
				Date:             date,
			},
			want: &domain.Statistics{
				TotalRevenue:     10000,
				TotalProductSold: 2,
				CompletedOrder:   1,
				CancelledOrder:   0,
				TotalOrder:       1,
				Date:             date,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "statistics" ("created_at","updated_at","deleted_at","total_revenue","total_product_sold","completed_order","cancelled_order","total_order","date") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(10000), int64(2), int64(1), int64(0), int64(1), date).
					WillReturnRows(sqlmock.NewRows([]string{"TotalRevenue", "TotalProductSold", "CompletedOrder", "CancelledOrder", "TotalOrder", "Date"}).
						AddRow(10000, 2, 1, 0, 1, date))
				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			statistic: domain.Statistics{
				TotalRevenue:     10000,
				TotalProductSold: 2,
				CompletedOrder:   1,
				CancelledOrder:   0,
				TotalOrder:       1,
				Date:             date,
			},
			want:    nil,
			wantErr: true,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "statistics" ("created_at","updated_at","deleted_at","total_revenue","total_product_sold","completed_order","cancelled_order","total_order","date") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(10000), int64(2), int64(1), int64(0), int64(1), date).
					WillReturnError(errors.New("mock error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ar := NewStatisticsRepository(gormdb, nil)
			res, err := ar.Create(context.TODO(), tt.statistic)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.want == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tt.want.TotalRevenue, res.TotalRevenue)
				assert.Equal(t, tt.want.TotalProductSold, res.TotalProductSold)
				assert.Equal(t, tt.want.CompletedOrder, res.CompletedOrder)
				assert.Equal(t, tt.want.CancelledOrder, res.CancelledOrder)
				assert.Equal(t, tt.want.TotalOrder, res.TotalOrder)
				assert.Equal(t, tt.want.Date, res.Date)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_statisticsRepository_Update(t *testing.T) {
	date := datatypes.Date(time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local))
	tests := []struct {
		name      string
		statistic domain.Statistics
		want      *domain.Statistics
		wantErr   bool
		mock      func()
	}{
		{
			name: "success",
			statistic: domain.Statistics{
				TotalRevenue:     10000,
				TotalProductSold: 2,
				CompletedOrder:   1,
				CancelledOrder:   0,
				TotalOrder:       1,
				Date:             date,
			},
			wantErr: false,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "statistics" WHERE Date = $1 AND "statistics"."deleted_at" IS NULL ORDER BY "statistics"."id" LIMIT 1`)).
					WithArgs(date).WillReturnRows(sqlmock.NewRows([]string{"TotalRevenue", "Date"}).
					AddRow(10000, date))
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "statistics" ("created_at","updated_at","deleted_at","total_revenue","total_product_sold","completed_order","cancelled_order","total_order","date") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(10000), int64(2), int64(1), int64(0), int64(1), date).
					WillReturnRows(sqlmock.NewRows([]string{"TotalRevenue", "TotalProductSold", "CompletedOrder", "CancelledOrder", "TotalOrder", "Date"}).
						AddRow(10000, 2, 1, 0, 1, date))
				mock.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ar := NewStatisticsRepository(gormdb, nil)
			_, err := ar.Update(context.TODO(), tt.statistic)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
