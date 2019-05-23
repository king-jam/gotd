package gif

import (
	"time"
)

type GIF struct {
	ID            uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
	DeactivatedAt time.Time
	GIF           string
	RequesterID   string
	RequestSrc    string
	Tags          []string
	db            DB
}
type GifService struct {
	repo Repo
}

func NewGifService(repo Repo) *GifService {
	return &GifService{repo: repo}
}

func (g *GifService) StoreGif(gif *GIF) error {
	err := g.repo.Insert(gif)
	if err != nil {
		return err
	}
	return nil
}

func (g *GifService) UpdateGif(gif *GIF) error {

	err := g.repo.Update(gif)
	if err != nil {
		return err
	}
	return nil
}

func (g *GifService) RemoveGifById(id int) error {
	err := g.repo.DeleteGIFByID(id)
	if err != nil {
		return err
	}
	return nil
}

func (g *GifService) GetGIFById(id uint) (GIF, error) {
	gif, err := g.repo.FindGIFByID(id)
	if err != nil {
		return GIF{}, err
	}
	object := GIF{
		GIF:         gif.GIF,
		RequesterID: gif.RequesterID,
		RequestSrc:  gif.RequestSrc,
		Tags:        gif.Tags,
	}
	return object, nil
}

func (g *GifService) GetAllGifs() ([]GIF, error) {
	var gifList []GIF
	gifs, err := g.repo.FindAllGifs()
	if err != nil {
		return gifList, err
	}
	for i := range gifs {
		gif := GIF{
			ID:          gifs[i].ID,
			GIF:         gifs[i].GIF,
			RequesterID: gifs[i].RequesterID,
			RequestSrc:  gifs[i].RequestSrc,
			Tags:        gifs[i].Tags,
		}
		gifList = append(gifList, gif)
	}
	return gifList, nil
}

func (g *GifService) GetMostRecent() (GIF, error) {
	dbGif, err := g.repo.LatestGIF()
	if err != nil {
		return GIF{}, err
	}
	gif := GIF{
		GIF:         dbGif.GIF,
		RequesterID: dbGif.RequesterID,
		RequestSrc:  dbGif.RequestSrc,
	}
	return gif, nil
}
