package db

import (
	"CommentTree/comment_tree/config"
	"CommentTree/comment_tree/core"
	"CommentTree/comment_tree/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	log  logger.Interface
	conn *pgx.Conn
}

func NewDB(log logger.Interface, cfg config.DBConfig) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), cfg.Addr)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS comments (
		id        SERIAL   PRIMARY KEY,
		parent_id INTEGER  REFERENCES comments(id) ON DELETE CASCADE,
	    username  TEXT     NOT NULL,
		text      TEXT     NOT NULL
	);`
	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		return nil, err
	}

	return &DB{
		log:  log,
		conn: conn,
	}, nil
}

func (db *DB) Add(ctx context.Context, comment core.Comment) (int, error) {
	if comment.ParentID != nil {
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM comments WHERE id = $1)`
		err := db.conn.QueryRow(ctx, checkQuery, *comment.ParentID).Scan(&exists)
		if err != nil {
			return 0, err
		}
		if !exists {
			return 0, pgx.ErrNoRows
		}
	}

	query := `
	 INSERT INTO comments (parent_id, username, text)
	 VALUES ($1, $2, $3)
	 RETURNING id
	`
	var insertedID int
	err := db.conn.QueryRow(ctx, query, comment.ParentID, comment.Username, comment.Text).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (db *DB) Get(ctx context.Context, id int) ([]core.Comment, error) {
	query := `
	 WITH RECURSIVE comment_tree AS (
	  SELECT id, parent_id, username, text
	  FROM comments
	  WHERE id = $1
	  UNION ALL
	  SELECT c.id, c.parent_id, c.username, c.text
	  FROM comments c
	  INNER JOIN comment_tree ct ON c.parent_id = ct.id
	 )
	 SELECT id, parent_id, username, text FROM comment_tree;
	 `

	rows, err := db.conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []core.Comment
	for rows.Next() {
		var comment core.Comment
		err = rows.Scan(&comment.ID, &comment.ParentID, &comment.Username, &comment.Text)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (db *DB) GetAll(ctx context.Context) ([]core.Comment, error) {
	query := `
		SELECT id, parent_id, username, text
		FROM comments
	`

	rows, err := db.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []core.Comment
	for rows.Next() {
		var comment core.Comment
		err = rows.Scan(&comment.ID, &comment.ParentID, &comment.Username, &comment.Text)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (db *DB) Delete(ctx context.Context, id int) error {
	query := `
        WITH RECURSIVE comment_tree AS (
            SELECT id
            FROM comments
            WHERE id = $1
            UNION ALL
            SELECT c.id
            FROM comments c
            INNER JOIN comment_tree ct ON c.parent_id = ct.id
        )
        DELETE FROM comments
        WHERE id IN (SELECT id FROM comment_tree);
    `

	_, err := db.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
