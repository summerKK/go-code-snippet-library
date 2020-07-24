package setting

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

var sections = make(map[string]interface{})

func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &Setting{vp: vp}
	// 检测配置文件修改
	s.WatchSettingChange()

	return s, nil
}

func (s Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			err := s.ReadAllSection()
			if err != nil {
				log.Errorf("s.ReadAllSection error:%v", err)
				return
			}
			log.Println("配置文件已更新")
		})
	}()
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}

	return nil
}

func (s *Setting) ReadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
