package model_test

import (
	"fmt"
	"testing"

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
