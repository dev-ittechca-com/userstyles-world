package models

import (
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var (
	errorReviewOutOfRange  = errors.New("rating is out of range")
	errorReviewLongComment = errors.New("comment is too long")
	errorReviewEmptyFields = errors.New("fields cannot be empty")
)

type Review struct {
	gorm.Model
	Rating  int
	Comment string

	User   User
	UserID uint

	Style   Style
	StyleID uint
}

func (Review) FindAllForStyle(id any) (q []Review, err error) {
	err = db().
		// Preload(clause.Associations).
		Preload("User").
		Model(modelReview).
		Order("id DESC").
		Find(&q, "style_id = ? ", id).
		Error

	if err != nil {
		return nil, err
	}

	return q, nil
}

func (r *Review) CreateForStyle() error {
	return db().Create(r).Error
}

func (r *Review) FindLastFromUser(styleID, userID any) error {
	return db().
		Preload("User").
		Model(modelReview).
		Order("id DESC").
		Last(&r, "style_id = ? and user_id = ?", styleID, userID).
		Error
}

// Validate verifies user-generated content.
func (r *Review) Validate() error {
	switch {
	case r.Rating < 0 || r.Rating > 5:
		return errorReviewOutOfRange
	case len(r.Comment) > 500:
		return errorReviewLongComment
	case r.Rating == 0 && len(r.Comment) == 0:
		return errorReviewEmptyFields
	default:
		return nil
	}
}

// NewReview is a helper for creating a new review.
func NewReview(uid, sid uint, rating, comment string) *Review {
	i, err := strconv.Atoi(rating)
	if err != nil {
		i = -1
	}

	return &Review{
		UserID:  uid,
		StyleID: sid,
		Comment: strings.TrimSpace(comment),
		Rating:  i,
	}
}
