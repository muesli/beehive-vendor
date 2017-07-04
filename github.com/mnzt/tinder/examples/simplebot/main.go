package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/mnzt/tinder"
)

// --------------------------------------------------------------------------------
// Globals used to log in. Change these to use the bot.
//
// See the link below for more details.
// https://gist.github.com/rtt/10403467
// --------------------------------------------------------------------------------

// myLatitude is the current latitude for you.
var myLatitude = float32(12.3456)

// myLongitude is the current longitude for you.
var myLongitude = float32(-12.3456)

// myPID is your Facebook PID.
var myPID = "pid"

// myToken is your Facebook auth token.
var myToken = "token"

// --------------------------------------------------------------------------------

// diff is used to find user age. It finds the difference between two times.
func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// pingAndSwipe sets your location and then swipes right on 20 singles around
// you.
func pingAndSwipe(tin *tinder.Tinder) error {
	// Make the location move around a little.
	randLong := rand.Float32() * 0.0001
	if rand.Int()%2 == 0 {
		randLong = -randLong
	}
	randLat := rand.Float32() * 0.0001
	if rand.Int()%2 == 0 {
		randLat = -randLat
	}

	localLon := myLongitude + randLong
	localLat := myLatitude + randLat

	fmt.Printf("Pinging and setting location to %v,%v...\n", localLon, localLat)
	err := tin.Ping(localLon, localLat)
	if err != nil {
		if !strings.Contains(err.Error(), "position change not significant") {
			return err
		}
	}

	time.Sleep(2 * time.Second)

	fmt.Printf("Swiping on all local matches\n")
	recs, err := tin.GetRecommendations(20)
	if err != nil {
		if strings.Contains(err.Error(), "recs timeout") {
			fmt.Printf("Got automated bot drop, waiting 30 minutes...\n")
			time.Sleep(30 * time.Minute)
			return nil
		}

		return err
	}

	for i := range recs.Results {
		time.Sleep(1 * time.Second)
		usr := recs.Results[i]
		years, _, _, _, _, _ := diff(usr.Birth, time.Now())

		fmt.Printf("Swiping right on %v (UID %v, age %v)\n", usr.Name, usr.ID,
			years)
		_, err := tin.Like(usr.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Random seed to use in fudging time and location.
	rand.Seed(time.Now().UnixNano())

	tin := tinder.Init(myPID, myToken)

	// Ctrl-C to kill. This gets stuck if the program is currently sleeping.
	exitChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			exitChan <- struct{}{}
		}
	}()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
mainLoop:
	for {
		select {
		case <-ticker.C:
			err := tin.Auth()
			if err != nil {
				fmt.Printf("Auth failed: %s\n", err.Error())
				os.Exit(1)
			}

			userInfo, err := tin.GetUser(tin.Me.User.ID)
			if userInfo.Results.Name == "" {
				fmt.Printf("Token expired. Quitting.\n")
				os.Exit(0)
			}
			if err != nil {
				fmt.Printf("Failed to get user info: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Logged in, current username is %v\n",
				userInfo.Results.Name)

			err = pingAndSwipe(tin)
			if err != nil {
				fmt.Printf("Received error on ping/swipe: %s\n", err.Error())
			}

			// Do pings/swipes at random times to make this look more like a
			// real user.
			time.Sleep(60 * time.Second)
			randomSeconds := rand.Intn(120)
			time.Sleep(time.Duration(randomSeconds) * time.Second)

		case <-exitChan:
			break mainLoop
		}
	}
}
