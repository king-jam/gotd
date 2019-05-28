package gif

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/king-jam/gotd/giphy"
	libgiphy "github.com/sanzaru/go-giphy"
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

// type Giphy struct {
// 	giphy *libgiphy.Giphy
// }

// func NewGiphy(api *libgiphy.Giphy) *Giphy {
// 	return &Giphy{giphy: api}
// }

func BuildGif(gif *GIF) error {
	var url *url.URL
	var err error

	api := libgiphy.NewGiphy("B4LxlW1Av7CPwzIJL7VAIOQE4Lc4wSKm")

	// Check if user input is a url
	validURL := govalidator.IsURL(gif.GIF)
	if validURL {
		fmt.Printf("is a valid url")
		// Reformat the URL
		url, err = url.Parse(gif.GIF)
		if err != nil {
			return err
		}

		//Get Tags From Gif URL
		tags, err := giphy.GetGIFTags(gif.GIF)
		if err != nil {
			return err
		}

		// Normalize the URL
		err = normalizeGiphyURL(url)
		if err != nil {
			return err
		}
		gif.GIF = url.String()
		gif.Tags = tags
	} else {
		res, err := api.GetSearch(gif.GIF, 1, -1, "pg", "", false)
		if err != nil {
			return err
		}
		gif.GIF = res.Data[0].Url
		gif.Tags = strings.Split(gif.GIF, " ")
		fmt.Printf("\n\n%+v", gif)
	}

	return nil

}

func (g *GifService) StoreGif(gif *GIF) error {
	// Add more details onto the gif, such as tags, and reformat the URL
	err := BuildGif(gif)
	if err != nil {
		return err
	}
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

	// If user post the same URL twice, do nothing and return
	if gif.GIF == lastGif.GIF {
		return nil
	}

	//Else, update the deactivate time for previous gif
	lastGif.DeactivatedAt = time.Now()
	err = g.UpdateGif(&lastGif)
	if err != nil {
		return err
	}
	// Insert gif into db
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
		Tags:        []string(dbGif.Tags),
	}
	return gif, nil
}

// validateURL will validate if URL is from giphy.com
func validateURL(url *url.URL) bool {
	// Validate if string is from giphy
	return url.Hostname() == "giphy.com"
}

// normalizeGiphyURL will add /fullscreen to URL
func normalizeGiphyURL(url *url.URL) error {
	if !validateURL(url) {
		return fmt.Errorf("Invalid URL - Use Giphy.com")
	}
	var fullPath string
	// Check if URL has "/fullscreen"
	ok, err := path.Match("/gifs/*/fullscreen", url.Path)
	if err != nil {
		return err
	}
	if !ok {
		fullPath = path.Join(url.Path, "fullscreen")
		url.Path = fullPath
	}
	return nil
}
