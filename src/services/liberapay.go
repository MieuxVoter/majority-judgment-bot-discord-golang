package services

import (
	"fmt"
	"github.com/andersfylling/disgord/json"
	"github.com/sarulabs/di"
	"io"
	"log"
	"main/src/container"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"xorm.io/xorm"
)

// MainApiEndpoint is described at the bottom of this file, in the Annex
const MainApiEndpoint = "https://liberapay.com/%s/public.json"

// Liberapay service fetches (and perhaps caches) financial data from liberapay.com
type Liberapay struct {
	config  *Config
	orm     *xorm.Engine
	main    *LiberapayMain
	updated time.Time
}

type (
	LiberapayMain struct {
		Goal      LiberapayAmount `json:"goal"`
		Receiving LiberapayAmount `json:"receiving"`
	}
	LiberapayAmount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	}
)

func (service *Liberapay) GetUserName() string {
	return url.PathEscape(service.config.Get("LIBERAPAY_USERNAME"))
}

func (service *Liberapay) GetMainEndpoint() string {
	return fmt.Sprintf(MainApiEndpoint, service.GetUserName())
}

func (service *Liberapay) FetchMain() error {
	response, err := http.Get(service.GetMainEndpoint())
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			fmt.Println("cannot close connection to liberapay")
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var result LiberapayMain
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	service.main = &result
	service.updated = time.Now()

	return nil
}

func (service *Liberapay) GetMain() (*LiberapayMain, error) {
	if service.main == nil {
		err := service.FetchMain()
		if err != nil {
			return nil, err
		}
	}

	if time.Now().Sub(service.updated).Hours() > 1.0 {
		err := service.FetchMain()
		if err != nil {
			return nil, err
		}
	}

	return service.main, nil
}

func (service *Liberapay) GetSurvivalChancePercentage() (float64, error) {
	data, err := service.GetMain()
	if err != nil {
		return 0.0, err
	}

	receiving, err := strconv.ParseFloat(data.Receiving.Amount, 64)
	if err != nil {
		return 0.0, err
	}
	goal, err := strconv.ParseFloat(data.Goal.Amount, 64)
	if err != nil {
		return 0.0, err
	}

	return 100.0 * receiving / goal, nil
}

func (service *Liberapay) GetSurvivalAsString() (string, error) {
	data, err := service.GetMain()
	if err != nil {
		return "???", err
	}

	survivalChance, err := service.GetSurvivalChancePercentage()
	if err != nil {
		return "?!?", err
	}

	currency := data.Goal.Currency
	if currency == "EUR" {
		currency = "€"
	}

	return fmt.Sprintf(
		"`%2.2f%%` _(`%s` / `%s` %s per week)_",
		survivalChance, data.Receiving.Amount, data.Goal.Amount, currency,
	), nil
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "liberapay",
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &Liberapay{
				config: ctn.Get("config").(*Config),
				orm:    ctn.Get("database.engine").(*xorm.Engine),
			}
			return cmd, nil
		},
	})
	if err != nil {
		log.Fatalf("service liberapay failed to build : %s\n", err)
	}
}

// ANNEX ---------------------------------------------------------------------------------------------------------------

// https://liberapay.com/MajorityJudgmentBot/public.json
/*
{
    "avatar": "https://seccdn.libravatar.org/avatar/89fbf2bf47f509416365c28d0d07d7bb?s=160&d=404",
    "display_name": "Majority Judgment Bot",
    "giving": {
        "amount": "0.00",
        "currency": "EUR"
    },
    "goal": {
        "amount": "2.00",
        "currency": "EUR"
    },
    "id": 1817045,
    "kind": "individual",
    "npatrons": 0,
    "receiving": {
        "amount": "0.00",
        "currency": "EUR"
    },
    "statements": [
        {
            "content": "\ud83e\udd16\ud83d\udde9 Greetings, kind human.\r\n\r\nMy purpose is to assist you in making rad polls using Majority Judgment.\r\n\r\nTo do so, I need a cosy home on some computer, a spacious warehouse to store your polls, a flowing connection to the internet, and some of that delicious Big Bang cake that you call energy.\r\n\r\nI also need a friend to sometimes grease my joints and massage my hardworking back.\r\n\r\nMy friend Patreon provides me with all of the above.\r\n\r\nIf you want me to keep operating, please consider making a small gift of life to me.\r\n\r\n\ud83c\udfdb\r\n",
            "lang": "en"
        }
    ],
    "summaries": [
        {
            "content": "My purpose is to assist you in making rad polls using Majority Judgment.",
            "lang": "en"
        }
    ],
    "username": "MajorityJudgmentBot"
}
*/
