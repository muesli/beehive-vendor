package gocricket

import (
	"encoding/xml"
	"strconv"
	"io/ioutil"
	"time"
	"net/http"
	"fmt"
)

const (
	EVENT_NO_CHANGE           = 0
	EVENT_OUT                 = 1
	EVENT_MATCH_STATUS_CHANGE = 2
	EVENT_OVER_CHANGED        = 3
	EVENT_RUN_CHANGE          = 4
	CRICBUZZ_URL              = "http://synd.cricbuzz.com/j2me/1.0/livematches.xml"
)

type Cricket struct {
	lastFetchedResult MatchStat
	teamName          string
	event             chan ResponseEvent
}

func NewCricketWatcher(teamName string, event chan ResponseEvent) (c *Cricket) {
	c = new(Cricket)
	c.teamName = teamName
	c.event = event
	return
}

type Response struct {
	BtTeamName  string
	Overs       string
	MatchStatus string
	Runs        string
	Wickets     string
}

type ResponseEvent struct {
	Response
	EventType int
}

type MatchData struct {
	MatchStats []MatchStat `xml:"match"`
}

type MatchStat struct {
	XMLName     xml.Name     `xml:"match"`
	Type        string       `xml:"type,attr"`
	States      State        `xml:"state"`
	Teams       []Team       `xml:"Tm"`
	BattingTeam *BattingTeam `xml:"mscr>btTm"`
}

type State struct {
	MatchState string `xml:"mchState,attr"`
	Status     string `xml:"status,attr"`
}

type Team struct {
	Name string `xml:"Name,attr"`
}

type InningDetails struct {
	Overs string `xml:"noofovers"`
}

type BattingTeam struct {
	Name  string   `xml:"sName,attr"`
	ID    string   `xml:"id,attr"`
	Inngs []Inning `xml:"Inngs"`
}

type Inning struct {
	Description string `xml:"desc,attr"`
	Run         string `xml:"r,attr"`
	Overs       string `xml:"ovrs,attr"`
	Wickets     string `xml:"wkts,attr"`
}

func (m *MatchData) Print() {
	for _, v := range m.MatchStats {
		fmt.Printf("%+v\n", v)
	}
}

func (m *MatchStat) convertToResponse(eventType int) ResponseEvent {
	return ResponseEvent{
		Response: Response{
			Overs: m.BattingTeam.Inngs[0].Overs,
			BtTeamName:m.BattingTeam.Name,
			Runs:m.BattingTeam.Inngs[0].Run,
			Wickets:m.BattingTeam.Inngs[0].Wickets,
		},
		EventType: eventType,
	}
}

func (m *MatchStat) TriggerEvent(lastFetchedStat MatchStat, event chan ResponseEvent) {
	var lastBt *BattingTeam
	var newBt *BattingTeam

	if lastFetchedStat.BattingTeam != nil {
		lastBt = lastFetchedStat.BattingTeam
	}

	if m.BattingTeam != nil {
		newBt = m.BattingTeam
	} else {
		event <- m.convertToResponse(EVENT_NO_CHANGE)
	}
	if newBt.Inngs != nil && len(newBt.Inngs) > 0 {
		inningIndex := len(newBt.Inngs) - 1
		in := newBt.Inngs[inningIndex]
		run, err := strconv.Atoi(in.Run)
		overs, err := strconv.ParseFloat(in.Overs, 32)
		wkts, err := strconv.Atoi(in.Wickets)
		if err != nil {
			event <- m.convertToResponse(EVENT_NO_CHANGE)
		}
		oldRun, _ := strconv.Atoi(lastBt.Inngs[inningIndex].Run)
		oldOvers, _ := strconv.ParseFloat(lastBt.Inngs[inningIndex].Overs, 32)
		oldWkts, _ := strconv.Atoi(lastBt.Inngs[inningIndex].Wickets)

		if oldRun != run {
			event <- m.convertToResponse(EVENT_RUN_CHANGE)
		}
		if int(oldOvers) != int(overs) {
			event <- m.convertToResponse(EVENT_OVER_CHANGED)
		}
		if oldWkts != wkts {
			event <- m.convertToResponse(EVENT_OUT)
		}
	}
}

func (c *Cricket) Start() {
	var temp MatchData
	go func() {
		for {
			var m MatchData
			resp, _ := http.Get(CRICBUZZ_URL)
			data, _ := ioutil.ReadAll(resp.Body)
			err := xml.Unmarshal(data, &m)
			if err != nil {
				fmt.Print("Error is", err)
			}
			matchStat := c.TeamMatchStat(m)
			if matchStat.BattingTeam != nil && len(temp.MatchStats) > 0{
				matchStat.TriggerEvent(c.TeamMatchStat(temp),c.event)
			}
			temp = m
			time.Sleep(time.Second * 60)
		}
	}()
}

func (c *Cricket ) TeamMatchStat(m MatchData) (s MatchStat) {
	for _, k := range m.MatchStats {
		for _, team := range k.Teams {
			if team.Name == c.teamName {
				s = k
				return
			}
		}
	}
	return
}
