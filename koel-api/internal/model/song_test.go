package model_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/summerKK/go-code-snippet-library/koel-api/global"
	boot "github.com/summerKK/go-code-snippet-library/koel-api/init"
	"github.com/summerKK/go-code-snippet-library/koel-api/internal/model"
)

func TestMain(m *testing.M) {
	boot.SetConfig([]string{"../../configs"})
	boot.Boot()

	m.Run()
}

func TestSong_All(t *testing.T) {
	song := &model.Song{}

	songs, err := song.All(global.DBEngine)
	if err != nil {
		t.Errorf("TestSong_All got error:%v", err)
	}

	fmt.Printf("%+v", songs)
}

func TestSong_Create(t *testing.T) {
	album := &model.Album{}
	albumRow, err := album.First(global.DBEngine)
	if err != nil {
		t.Errorf("TestSong_Create got error:%v", err)
	}

	artist := &model.Artist{}
	artistRow, err := artist.First(global.DBEngine)
	if err != nil {
		t.Errorf("TestSong_Create got error:%v", err)
	}
	song := &model.Song{
		AlbumId:  albumRow.ID,
		ArtistId: artistRow.ID,
		Title:    "Taylor Swift",
		Length:   10.08,
		Track:    0,
		Path:     "/Users/summer/Docker/www/summer/koel/tests/songs/subdir/sic.mp3",
		Mtime:    time.Now().Unix(),
	}

	err = song.Create(global.DBEngine)
	if err != nil {
		t.Errorf("TestSong_Create got error:%v", err)
	}
}
