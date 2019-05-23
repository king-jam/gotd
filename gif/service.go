package gif

import (
	"fmt"
	"time"
)

type GIF struct {
	ID            uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
	DeactivatedAt time.Time
	GIF           string `json:"url"`
	RequesterID   string
	RequestSrc    string
	Tags          []string
}

type GifService struct {
	repo Repo
}

func NewGifService(repo Repo) *GifService {
	return &GifService{repo: repo}
}

func (g *GifService) StoreGif(gif *GIF) error {
	//Update deactive time for previous gif before storing new gif
	lastGif, err := g.GetMostRecent()
	if err != nil {
		// If there is no previous gif, then store new gif
		if err == ErrRecordNotFound {
			err = g.StoreGif(gif)
			if err != nil {
				return err
			}
			return err
		}
		return err
	}

	//Else, update the deactivate time for previous gif

	lastGif.DeactivatedAt = time.Now()
	fmt.Printf("\n\n%+v\n\n", lastGif)
	err = g.UpdateGif(&lastGif)
	if err != nil {
		return err
	}

	err = g.repo.Insert(gif)
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
