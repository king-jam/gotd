package postgres

import (
	"net/url"

	"github.com/go-pg/pg"
)
type DBClient struct {
	db *pg.DB
}

type Schema struct {
	gif string
}

type GOTD struct {
	ID string,
	url string
}

// InitDatabase takes a connection string URL to pass into the Database
func InitDatabase(url *url.URL) (*DBClient, error) {
	options, err := pg.ParseURL(url.String())
	if err != nil {
		return nil, fmt.Errorf("Failure to parse opts: %s," err)
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
		panicIf(err)
		
	}
	return &DBClient{
		Db: db,
	}, nil
}

func (c *DBClient) Insert(db, gotd GOTD) error {
	err := c.db.Insert(&gotd)
	if err != nil {
		return err
	}
	return nil
}
// Close wraps the db close function for easy cleanup
func (c *DBClient) Close() {
	b.db.Close()
}
