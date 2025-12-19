package repository_test

import (
	"auth-service/model"
	"auth-service/repository"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAssignmentRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewAssignmentsRepository(db)

	start := time.Now()
	end := start.Add(24 * time.Hour)

	assignment := &model.DriverAssignment{
		VehicleID:  1,
		StartDate:  start,
		EndDate:    end,
		TotalTrips: 1,
		DriverName: "John",
		Status:     "active",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO driver_assignments
		(vehicle_id, start_date, end_date, total_trips, driver_name, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`)).
		WithArgs(assignment.VehicleID, assignment.StartDate, assignment.EndDate, assignment.TotalTrips, assignment.DriverName, assignment.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.Create(assignment)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), assignment.ID)
}

func TestAssignmentRepository_FindByVehicle(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewAssignmentsRepository(db)

	start := time.Now()
	end := start.Add(24 * time.Hour)

	rows := sqlmock.NewRows([]string{"id", "vehicle_id", "start_date", "end_date", "total_trips", "driver_name", "status"}).
		AddRow(1, 1, start, end, 1, "John", "active")

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, vehicle_id, start_date, end_date, total_trips, driver_name, status
		FROM driver_assignments
		WHERE vehicle_id = $1
	`)).WithArgs(1).WillReturnRows(rows)

	result, err := repo.FindByVehicle(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "John", result[0].DriverName)
}

func TestAssignmentRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewAssignmentsRepository(db)

	start := time.Now()
	end := start.Add(24 * time.Hour)

	assignment := &model.DriverAssignment{
		ID:         1,
		VehicleID:  1,
		StartDate:  start,
		EndDate:    end,
		TotalTrips: 10,
		DriverName: "Jane",
		Status:     "inactive",
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE driver_assignments
		SET vehicle_id = $1,
			start_date = $2,
			end_date = $3,
			total_trips = $4,
			driver_name = $5,
			status = $6
		WHERE id = $7
	`)).
		WithArgs(assignment.VehicleID, assignment.StartDate, assignment.EndDate, assignment.TotalTrips, assignment.DriverName, assignment.Status, assignment.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(assignment)
	assert.NoError(t, err)
}

func TestFindByVehicle_QueryError1(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewAssignmentsRepository(db)

	mock.ExpectQuery("FROM driver_assignments").
		WithArgs(uint(1)).
		WillReturnError(errors.New("query error"))

	result, err := repo.FindByVehicle(1)

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestFindByVehicle_ScanError2(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewAssignmentsRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "vehicle_id", "start_date",
	}).AddRow(1, 1, time.Now())

	mock.ExpectQuery("FROM driver_assignments").
		WillReturnRows(rows)

	result, err := repo.FindByVehicle(1)

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestFindByVehicle_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewAssignmentsRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "vehicle_id", "start_date", "end_date",
		"total_trips", "driver_name", "status",
	}).AddRow(
		1,
		1,
		time.Now(),
		time.Now(),
		10,
		"Driver A",
		"active",
	)

	mock.ExpectQuery("FROM driver_assignments").
		WithArgs(uint(1)).
		WillReturnRows(rows)

	result, err := repo.FindByVehicle(1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestAssignmentRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewAssignmentsRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM driver_assignments WHERE id = $1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(1)
	assert.NoError(t, err)
}
