package config

import (
	"os"

	"git.sr.ht/~barveyhirdman/chainkills/common"
	"gopkg.in/yaml.v3"
)

var c *Cfg

type Cfg struct {
	AdminName       string   `json:"admin_name"`
	AdminEmail      string   `json:"admin_email"`
	AppName         string   `json:"app_name"`
	Version         string   `json:"version"`
	RefreshInterval int      `json:"refresh_interval"`
	OnlyWHKills     bool     `json:"only_wh_kills"`
	IgnoreSystems   []string `json:"ignore_systems"`
	Redict          Redict   `json:"redict"`
	Wanderer        Wanderer `json:"wanderer"`
	Discord         Discord  `json:"discord"`
	Friends         Friends  `json:"friends"`
}

type Redict struct {
	Address  string `json:"address"`
	Database int    `json:"database"`
	TTL      int    `json:"ttl"` // Time to live for keys in minutes
}

type Wanderer struct {
	Token string `json:"token"`
	Slug  string `json:"slug"`
	Host  string `json:"host"`
}

type Discord struct {
	Token   string
	Channel string
}

type Friends struct {
	Alliances    []uint64
	Corporations []uint64
	Characters   []uint64
}

func (c *Cfg) IsFriend(allianceID, corpID, CharacterID uint64) bool {
	return common.Contains(c.Friends.Alliances, allianceID) ||
		common.Contains(c.Friends.Corporations, corpID) ||
		common.Contains(c.Friends.Characters, CharacterID)
}

func Read(path string) error {
	fp, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	cfg := Cfg{
		RefreshInterval: 60, // Default value
		Redict: Redict{
			TTL: 60,
		},
	}

	if err := yaml.NewDecoder(fp).Decode(&cfg); err != nil {
		return err
	}

	c = &cfg

	return nil
}

func Get() *Cfg {
	return c
}
