package repositories

import (
	"database/sql"
	"fmt"

	"github.com/ntk221/refactor_notion_backend/models"
)

const (
	articleNumPerPage = 5
)

func InserArticle(db *sql.DB, article models.Article) (models.Article, error) {
	const sqlInsertArticle = `
		insert into articles (title, contents, username)
		values (?, ?, ?);
	`

	var newArticle models.Article
	newArticle.Title, newArticle.Contents, newArticle.UserName = article.Title, article.Contents, article.UserName

	result, err := db.Exec(sqlInsertArticle, newArticle.Title, newArticle.Contents, newArticle.UserName)
	if err != nil {
		return models.Article{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Article{}, err
	}

	newArticle.ID = uint(id)

	return newArticle, nil
}

func SelectArticleList(db *sql.DB, page uint) ([]models.Article, error) {
	const sqlStr = `
	select article_id, title, contents, username, nice
	from articles
	limit ? offset ?;
	`

	rows, err := db.Query(sqlStr, articleNumPerPage, ((page - 1) * articleNumPerPage))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Contents, &article.UserName, &article.NiceNum)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func GetArticleByID(db *sql.DB, articleID uint) (models.Article, error) {
	const sqlStr = `
		select *
		from articles
		where article_id = ?;
	`

	fmt.Println(articleID)

	row := db.QueryRow(sqlStr, articleID)
	if err := row.Err(); err != nil {
		return models.Article{}, err
	}

	var article models.Article
	var createdTime sql.NullTime
	err := row.Scan(&article.ID, &article.Title, &article.Contents, &article.UserName, &article.NiceNum, &createdTime)
	if err != nil {
		return models.Article{}, err
	}

	if createdTime.Valid {
		article.CreatedAt = createdTime.Time
	}

	return article, nil
}

func UpdateArticleNice(db *sql.DB, articleID uint) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	const sqlGetNice = `
	select nice
	from articles
	where article_id = ?;
	`

	row := tx.QueryRow(sqlGetNice, articleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	var nicenum int
	err = row.Scan(&nicenum)
	if err != nil {
		tx.Rollback()
		return err
	}

	const sqlUpdateNice = `update articles set nice = ? where article_id = ?;`
	_, err = tx.Exec(sqlUpdateNice, nicenum+1, articleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
