package main

import (
	"fmt"
	"sort"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
)

// type Player struct {
// 	Id      int64   `gorm:"primaryKey" json:"id"`
// 	Name    string  `json:"name"`
// 	Surname string  `json:"surname"`
// 	Email   *string `json:"email"`
// 	Avatar  *string `json:"avatar"`
// 	Hidden  bool    `json:"hidden"`
// }

func GetPlayerDetail(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	var rs model.Player
	err = DB.QueryRowx(`
		select
			*
		from
			tboardgameplayers
		where
			id = ?
	`, id).StructScan(&rs)
	if err != nil {
		return err
	}

	res.Set("player", rs)
	return nil
}

type Play struct {
	Id          int64                 `json:"id,omitempty"`
	BoardgameId int64                 `json:"boardgame_id,omitempty"`
	Name        string                `json:"name,omitempty"`
	Square200   models.JsonNullString `json:"url,omitempty"`
	Winners     models.JsonArray      `json:"winners,omitempty"`
	Players     models.JsonArray      `json:"players,omitempty"`
	Cooperative models.JsonNullString `json:"cooperative,omitempty"`
}

func GetPlayerPlays(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	url, err := req.GetString("url")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	RDB, err := req.GetRedis()
	if err != nil {
		return err
	}

	var rs []Play
	err = db.CachedSelect(DB, RDB, "GetPlayerPlays"+url, &rs, fmt.Sprintf(`
		select
			p.id,
			p.boardgame_id,
			g.name,
			g.square200,
			json_extract(p.play_data, '$.winners') winners,
			json_extract(p.play_data, '$.players') players,
			json_extract(p.play_data, '$.cooperative') cooperative
		from
			tboardgamestats s,
			tboardgameplays p,
			tboardgames g
		where
			s.player_id = ? and
			g.id = s.boardgame_id and
			p.id = s.play_id %s
	`, db.YearConstraint(req, "and")), id)
	if err != nil {
		return err
	}

	won := 0
	count := 0
	cooperative_count := 0
	cooperative_won := 0

	type percent struct {
		Idx   int
		Won   int
		Count int
	}

	type best_game struct {
		Id        int64                 `json:"id,omitempty"`
		Name      string                `json:"name,omitempty"`
		Square200 models.JsonNullString `json:"url,omitempty"`
		Won       int                   `json:"won,omitempty"`
		Count     int
		Percent   float64
		Printable string `json:"count,omitempty"`
	}

	player_counts := map[int]percent{}
	game_counts := map[int64]percent{}

	for i, play := range rs {
		var winners []int64
		var players []int64

		err := play.Winners.Unmarshal(&winners)
		if err != nil {
			return err
		}

		err = play.Players.Unmarshal(&players)
		if err != nil {
			return err
		}

		n := len(players)
		if _, ok := player_counts[n]; !ok {
			player_counts[n] = percent{Won: 0, Count: 0}
		}

		if _, ok := game_counts[play.BoardgameId]; !ok {
			game_counts[play.BoardgameId] = percent{Won: 0, Count: 0}
		}

		cur_player := player_counts[n]
		cur_games := game_counts[play.BoardgameId]

		if play.Cooperative.Valid && play.Cooperative.String == "true" {
			cooperative_count++
			// cur.Count++

			for _, winner := range winners {
				if winner == id {
					cooperative_won++
					// cur.Won++
				}
			}
		} else {
			count++
			cur_player.Count++
			cur_games.Count++

			for _, winner := range winners {
				if winner == id {
					won++
					cur_player.Won++
					cur_games.Won++
				}
			}
		}

		cur_games.Idx = i
		player_counts[n] = cur_player
		game_counts[play.BoardgameId] = cur_games
	}

	best_games := []best_game{}
	for key, val := range game_counts {
		tmp := best_game{
			Id:        key,
			Name:      rs[val.Idx].Name,
			Square200: rs[val.Idx].Square200,
			Won:       val.Won,
			Count:     val.Count,
		}

		if tmp.Count > 0 {
			tmp.Percent = float64(tmp.Won) / float64(tmp.Count)
			tmp.Printable = fmt.Sprintf("%.2f%%", tmp.Percent)
		}

		if tmp.Count > 5 && tmp.Percent > 0.10 {
			best_games = append(best_games, tmp)
		}
	}

	sort.Slice(best_games, func(i, j int) bool {
		return best_games[i].Percent > best_games[j].Percent
	})

	plays3_won := 0
	plays3_count := 0

	for i := 3; i < len(player_counts); i++ {
		plays3_won += player_counts[i].Won
		plays3_count += player_counts[i].Count
	}

	res.Set("player_counts", player_counts)
	res.Set("best_games", best_games)
	res.Set("cooperative", cooperative_count)
	res.Set("cooperative_won", cooperative_won)
	res.Set("plays_count", count)
	res.Set("plays_won", won)
	res.Set("plays3_won", plays3_won)
	res.Set("plays3_count", plays3_count)

	if cooperative_count > 0 {
		res.Set("cooperative_per", float64(cooperative_won)/float64(cooperative_count))
	} else {
		res.Set("cooperative_per", 0)
	}

	if count > 0 {
		res.Set("per", float64(won)/float64(count))
	} else {
		res.Set("per", 0)
	}

	if plays3_count > 0 {
		res.Set("plays3_per", float64(plays3_won)/float64(plays3_count))
	} else {
		res.Set("playes3_per", 0)
	}

	return nil
}
