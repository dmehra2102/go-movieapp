package rating

import (
	"context"
	"errors"

	"movieexample.com/rating/internal/repository"
	model "movieexample.com/rating/pkg"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepostiory interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

// Controller defines a rating service controller.
type Controller struct {
	repo ratingRepostiory
}

// New creates a rating service controller.
func New(repo ratingRepostiory) *Controller {
	return &Controller{repo}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are no ratings for it.
func (c *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

// PutRating writes a rating for a given record.
func (c *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Put(ctx, recordID, recordType, rating)
}
