package model

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/boggydigital/yt_urls"
	"github.com/fhs/gompd/mpd"
	"gorm.io/gorm"
)

type Playlist struct {
	Id         int64  `gorm:"primaryKey" json:"id"`
	PlaylistId string `gorm:"-" json:"playlist_id"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	Downloaded bool   `json:"downloaded"`
}

func (Playlist) TableName() string {
	return "tplaylist"
}

func (Playlist) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (Playlist) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Playlist
	rs := db.Scopes(scopes...).Find(&data)
	if rs.Error != nil {
		return nil, rs.Error
	}

	retval := []Playlist{}
	for _, item := range data {
		playlist_id, err := yt_urls.PlaylistId(item.Url)
		if err != nil {
			return nil, err
		}

		_, err = os.Stat("/data/icecast/playlists/" + playlist_id)
		if !os.IsNotExist(err) {
			item.Downloaded = true
		}

		item.PlaylistId = playlist_id
		retval = append(retval, item)
	}

	return retval, nil
}

func (Playlist) Get(db *gorm.DB, id int64) (any, error) {
	var data Playlist
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj Playlist) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Playlist{
		Id: id,
	}

	var payload Playlist
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Model(&model).Updates(payload)
	if rs.Error != nil {
		return nil, err
	}

	return obj.Get(db, id)
}

func (Playlist) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Playlist
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return payload, nil
}

func (obj Playlist) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Playlist{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}

func PlaylistPlay(req *Map, res *Map) error {
	playlist_id, err := req.GetString("id")
	if err != nil {
		return err
	}

	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Add("playlists/" + playlist_id)
	if err != nil {
		return err
	}

	err = conn.Play(-1)
	if err != nil {
		return err
	}

	return nil
}

func PlaylistStop(req *Map, res *Map) error {
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Clear()
	if err != nil {
		return err
	}

	return nil
}

func M3u(req *Map, res *Map) error {
	playlist_id, err := req.GetString("id")
	if err != nil {
		return err
	}

	pwd := "/data/icecast/playlists/" + playlist_id
	_, err = os.Stat(pwd)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(pwd)
	if err != nil {
		return err
	}

	f, err := os.Create("/data/icecast/playlist.m3u")
	if err != nil {
		return err
	}
	defer f.Close()

	for _, e := range entries {
		fmt.Fprintln(f, "/data/icecast/playlists/"+playlist_id+"/"+e.Name())
	}

	return nil
}
