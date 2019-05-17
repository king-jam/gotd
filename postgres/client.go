package postgres

import (
	"fmt"
	"net/url"

	"github.com/alexbyk/panicif"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type DBClient struct {
	db *pg.DB
}

type Schema struct {
	gif string
}

type GOTD struct {
	ID  string
	URL string
}

// InitDatabase takes a connection string URL to pass into the Database
func InitDatabase(url *url.URL) (*DBClient, error) {
	options, err := pg.ParseURL(url.String())
	if err != nil {
		return nil, fmt.Errorf("Failure to parse opts: %s", err)
	}

	db := pg.Connect(options)

	// if !db.HasTable(&TokenData{}) {
	// 	db.CreateTable(&TokenData{})
	// }

	//Creating Schema
	for _, model := range []interface{}{Schema{gif: "gifUrl"}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:          false, // create temp table
			FKConstraints: false,
		})
		panicif.Err(err)

	}
	return &DBClient{
		db: db,
	}, nil
}

func (c *DBClient) Insert(gotd GOTD) error {
	err := c.db.Insert(&gotd)
	if err != nil {
		return err
	}
	return nil
}

// Close wraps the db close function for easy cleanup
func (c *DBClient) Close() {
	c.db.Close()
}
