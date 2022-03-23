package repositories

import (
	"bufio"
	"context"
	"errors"
	"os"

	"github.com/IllicLanthresh/pack-and-go/internal/models"
)

type cityFileRepo struct {
	filePath string
}

func NewCityFileRepo(filePath string) *cityFileRepo {
	return &cityFileRepo{filePath: filePath}
}

func (t cityFileRepo) ReadById(ctx context.Context, id int64) (*models.City, error) {
	file, err := os.Open("./db/cities.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	city := models.City{ID: id}
	var line int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			line++
			if line == id {
				city.Name = scanner.Text()
				return &city, nil
			}
		}
	}

	if err = scanner.Err(); err == nil {
		return nil, errors.New("could not find record")
	} else {
		return nil, err
	}
}

func (t cityFileRepo) ReadByName(ctx context.Context, name string) (*models.City, error) {
	file, err := os.Open("./db/cities.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	city := models.City{Name: name}
	var line int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			line++
			lineName := scanner.Text()
			if name == lineName {
				city.ID = line
				return &city, nil
			}
		}
	}

	if err = scanner.Err(); err == nil {
		return nil, errors.New("could not find record")
	} else {
		return nil, err
	}
}
