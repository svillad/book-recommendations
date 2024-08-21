package stores

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/book-recommendations/service/models"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const (
	tableBook = "book"
)

// BookStore specifies the methods to get books
type BookStore interface {
	GetBooks(ctx context.Context, req models.BookRequest) ([]models.Book, error)
}

type bookStore struct {
	logger *log.Entry
	db     *sqlx.DB
}

func NewBookStore(logger *log.Entry, db *sqlx.DB) BookStore {
	return &bookStore{
		logger: logger,
		db:     db,
	}
}

func (s *bookStore) GetBooks(ctx context.Context, req models.BookRequest) ([]models.Book, error) {
	var joins string
	var wheres []string
	var values []interface{}
	var limit string

	joins = joins + " JOIN author AS au ON au.id = bo.author_id"
	if len(req.Authors) > 0 {
		values = append(values, req.Authors)
		wheres = append(wheres, fmt.Sprintf("au.id IN (%s)", req.Authors))
	}

	joins = joins + " JOIN genre AS ge ON ge.id = bo.genre_id"
	if len(req.Genres) > 0 {
		values = append(values, req.Genres)
		wheres = append(wheres, fmt.Sprintf("ge.id IN (%s)", req.Genres))
	}

	if req.MinYear != "" || req.MaxYear != "" {
		if req.MinYear == "" {
			req.MinYear = strconv.Itoa(models.MinYear)
		}
		if req.MaxYear == "" {
			req.MaxYear = strconv.Itoa(models.MaxYear)
		}
		wheres = append(wheres, fmt.Sprintf(`year_published BETWEEN %s AND %s`, req.MinYear, req.MaxYear))
	}

	if req.MinPages != "" || req.MaxPages != "" {
		if req.MinPages == "" {
			req.MinPages = strconv.Itoa(models.MinPages)
		}
		if req.MaxPages == "" {
			req.MaxPages = strconv.Itoa(models.MaxPages)
		}
		wheres = append(wheres, fmt.Sprintf(`pages BETWEEN %s AND %s`, req.MinPages, req.MaxPages))
	}

	if req.Limit != "" {
		limit = " LIMIT " + req.Limit
	}

	var whereConditions string
	if len(wheres) > 0 {
		whereConditions = ` WHERE ` + strings.Join(wheres[:], " AND ")
	}

	query := `SELECT
	bo.id, bo.title, bo.year_published, bo.rating, bo.pages,
	au.id, au.first_name, au.last_name,
	ge.id, ge.title
	FROM book AS bo` + joins + whereConditions + ` ORDER BY rating DESC` + limit

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error while building query: %w", err)
	}
	defer func() {
		errClose := rows.Close()
		errRows := rows.Err()
		if errClose != nil || errRows != nil {
			s.logger.WithFields(log.Fields{
				"errClose": errClose,
				"errRows":  errRows,
			}).Error("something went wrong while closing rows")
		}
	}()
	books := make([]models.Book, 0)
	count := 0
	for rows.Next() {
		count++
		var (
			id            int64
			title         string
			yearPublished int64
			rating        float64
			pages         int64
			authorID      int64
			authorFirst   string
			authorLast    string
			genreID       int64
			genreTitle    string
		)
		if err := rows.Scan(&id, &title, &yearPublished, &rating, &pages, &authorID, &authorFirst, &authorLast, &genreID, &genreTitle); err != nil {
			return nil, fmt.Errorf("error getting books: %w", err)
		}
		books = append(
			books,
			models.Book{
				ID:            int64(count),
				Title:         title,
				YearPublished: yearPublished,
				Rating:        rating,
				Pages:         pages,
				Genre: models.Genre{
					ID:    genreID,
					Title: genreTitle,
				},
				Author: models.Author{
					ID:        authorID,
					FirstName: authorFirst,
					LastName:  authorLast,
				},
			},
		)
	}

	return books, nil
}
