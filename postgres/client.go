package postgres

import (
	"fmt"
	"net/url"

	"github.com/alexbyk/panicif"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type DBClient struct {
	database *pg.DB
}

type GOTD struct {
	ID  int
	GIF string
}

// InitDatabase takes a connection string URL to pass into the Database
func InitDatabase(url *url.URL) (*DBClient, error) {
	options, err := pg.ParseURL(url.String())
	if err != nil {
		return nil, fmt.Errorf("Failure to parse opts: %s", err)
	}

	pgdb := pg.Connect(options)

	// if !db.HasTable(&TokenData{}) {
	// 	db.CreateTable(&TokenData{})
	// }

	err = pgdb.CreateTable(&GOTD{}, &orm.CreateTableOptions{
		Temp:          false, // create temp table
		FKConstraints: false,
		IfNotExists:   true,
	})
	panicif.Err(err)

	return &DBClient{
		database: pgdb,
	}, nil
}

func (c *DBClient) Insert(gif GOTD) error {
	err := c.database.Insert(gif)
	if err != nil {
		return err
	}
	return nil
}

// Close wraps the db close function for easy cleanup
func (c *DBClient) Close() {
	c.database.Close()
}
