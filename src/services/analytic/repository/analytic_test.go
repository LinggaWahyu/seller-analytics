package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tokopedia-workshop-2022/seller-analytics-solution/src/services/analytic/domain"
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

func Test_analyticRepository_GetAnalyticByDate(t *testing.T) {
	date := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name    string
		date    time.Time
		want    *domain.Analytic
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			date: date,
			want: &domain.Analytic{
				AverageOrderValue: 100,
				Date:              datatypes.Date(date),
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "analytics" WHERE Date = $1 AND "analytics"."deleted_at" IS NULL ORDER BY "analytics"."id" LIMIT 1`)).
					WithArgs(date).
					WillReturnRows(sqlmock.NewRows([]string{"AverageOrderValue", "Date"}).
						AddRow(100, date))
			},
		},
		{
			name:    "error",
			date:    date,
			want:    nil,
			wantErr: true,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "analytics" WHERE Date = $1 AND "analytics"."deleted_at" IS NULL ORDER BY "analytics"."id" LIMIT 1`)).
					WithArgs(date).WillReturnError(errors.New("mock error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ar := NewAnalyticRepository(gormdb)
			res, err := ar.GetAnalyticByDate(context.TODO(), tt.date)
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

func Test_analyticRepository_CreateAnalytic(t *testing.T) {
	date := datatypes.Date(time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local))
	tests := []struct {
		name     string
		analytic domain.Analytic
		want     *domain.Analytic
		wantErr  bool
		mock     func()
	}{
		{
			name: "success",
			analytic: domain.Analytic{
				AverageOrderValue:     100,
				SalesConvertionRate:   90,
				CancellationOrderRate: 10,
				Date:                  date,
			},
			want: &domain.Analytic{
				AverageOrderValue:     100,
				SalesConvertionRate:   90,
				CancellationOrderRate: 10,
				Date:                  date,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "analytics" ("created_at","updated_at","deleted_at","average_order_value","sales_convertion_rate","cancellation_order_rate","date") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), float64(100), float64(90), float64(10), date).
					WillReturnRows(sqlmock.NewRows([]string{"AvergeOrderValue", "SalesConvertionRate", "CancelationOrderRate", "Date"}).
						AddRow(100, 90, 10, date))
				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			analytic: domain.Analytic{
				AverageOrderValue:     100,
				SalesConvertionRate:   90,
				CancellationOrderRate: 10,
				Date:                  date,
			},
			want:    nil,
			wantErr: true,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "analytics" ("created_at","updated_at","deleted_at","average_order_value","sales_convertion_rate","cancellation_order_rate","date") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), float64(100), float64(90), float64(10), date).
					WillReturnError(errors.New("mock error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ar := NewAnalyticRepository(gormdb)
			res, err := ar.CreateAnalytic(context.TODO(), tt.analytic)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.want == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tt.want.AverageOrderValue, res.AverageOrderValue)
				assert.Equal(t, tt.want.SalesConvertionRate, res.SalesConvertionRate)
				assert.Equal(t, tt.want.CancellationOrderRate, res.CancellationOrderRate)
				assert.Equal(t, tt.want.Date, res.Date)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func Test_analyticRepository_UpdateAnalytic(t *testing.T) {
	date := datatypes.Date(time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local))
	tests := []struct {
		name     string
		analytic domain.Analytic
		wantErr  bool
		mock     func()
	}{
		{
			name: "error get",
			analytic: domain.Analytic{
				AverageOrderValue:     100,
				SalesConvertionRate:   90,
				CancellationOrderRate: 10,
				Date:                  date,
			},
			wantErr: true,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "analytics" WHERE Date = $1 AND "analytics"."deleted_at" IS NULL ORDER BY "analytics"."id" LIMIT 1`)).
					WithArgs(date).WillReturnError(errors.New("mock error"))
			},
		},
		{
			name: "success",
			analytic: domain.Analytic{
				AverageOrderValue:     100,
				SalesConvertionRate:   90,
				CancellationOrderRate: 10,
				Date:                  date,
			},
			wantErr: false,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "analytics" WHERE Date = $1 AND "analytics"."deleted_at" IS NULL ORDER BY "analytics"."id" LIMIT 1`)).
					WithArgs(date).WillReturnRows(sqlmock.NewRows([]string{"AverageOrderValue", "Date"}).
					AddRow(50, date))
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "analytics" ("created_at","updated_at","deleted_at","average_order_value","sales_convertion_rate","cancellation_order_rate","date") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), float64(100), float64(90), float64(10), date).
					WillReturnRows(sqlmock.NewRows([]string{"AvergeOrderValue", "SalesConvertionRate", "CancelationOrderRate", "Date"}).
						AddRow(100, 90, 10, date))
				mock.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			ar := NewAnalyticRepository(gormdb)
			_, err := ar.UpdateAnalytic(context.TODO(), tt.analytic)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
