package repositories

import (
	"context"
	"database/sql"

	"github.com/IllicLanthresh/pack-and-go/internal/models"
)

type tripDbRepo struct {
	db *sql.DB
}

func NewTripDbRepo(db *sql.DB) *tripDbRepo {
	return &tripDbRepo{db: db}
}

func (t tripDbRepo) ListAll(ctx context.Context) ([]*models.Trip, error) {
	rows, err := t.db.QueryContext(
		ctx,
		`SELECT id,
                       origin_id,
                       destination_id,
                       dates_bitmask,
                       price
                FROM Trip`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trips := make([]*models.Trip, 0)
	for rows.Next() {
		var trip = models.Trip{
			Origin:      models.City{Incomplete: true},
			Destination: models.City{Incomplete: true},
		}

		if err = rows.Scan(
			&trip.Id,
			&trip.Origin.ID,
			&trip.Destination.ID,
			&trip.Dates,
			&trip.Price,
		); err != nil {
			return nil, err
		}

		trips = append(trips, &trip)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trips, nil
}

func (t tripDbRepo) Create(ctx context.Context, trip *models.Trip) (int64, error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // nolint: errcheck

	var r sql.Result
	r, err = tx.ExecContext(
		ctx,
		`INSERT INTO Trip (origin_id, destination_id, dates_bitmask, price) VALUES (?, ?, ?, ?)`,
		trip.Origin.ID,
		trip.Destination.ID,
		trip.Dates,
		trip.Price,
	)
	if err != nil {
		return 0, err
	}
	rowId, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return rowId, err
}

func (t tripDbRepo) ReadById(ctx context.Context, id int64) (*models.Trip, error) {
	var trip = models.Trip{
		Origin:      models.City{Incomplete: true},
		Destination: models.City{Incomplete: true},
	}

	err := t.db.QueryRowContext(
		ctx,
		`SELECT id,
                       origin_id,
                       destination_id,
                       dates_bitmask,
                       price
                FROM Trip
                WHERE id = ?`,
		id,
	).Scan(
		&trip.Id,
		&trip.Origin.ID,
		&trip.Destination.ID,
		&trip.Dates,
		&trip.Price,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &trip, nil
}
