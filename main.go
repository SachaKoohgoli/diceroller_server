package main

import (
	"diceroller_server/helpers"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	urlRouter := mux.NewRouter()

	// urlRouter.HandleFunc("/", homeHandler)

	rollHandlerFunc := func(writer http.ResponseWriter, request *http.Request) {
		wellFormattedKeyRegex := regexp.MustCompile(`(\d)+d(\d)+`)

		diceRequests := make(map[string]int)

		params := request.URL.Query()

		for key, value := range params {
			if key == "token" {
				fmt.Printf("encoded token: %s", value)
				accessToken, err := helpers.DecodeToken(value[0])
				if err != nil {
					renderError(HttpError{
						HttpCode:  400,
						ErrorCode: "MALFORMED_INPUT",
						Msg:       err.Error(),
					}, writer)
				}
				if !accessToken.IsValid(time.Now().UTC(), 3*time.Minute) {
					renderError(HttpError{
						HttpCode:  403,
						ErrorCode: "INVALD_TOKEN",
						Msg:       "Your token is invalid, please generate a new one",
					}, writer)
				}
			}
			if key == "dice" {
				allDice := value
				fmt.Printf("%s: %s\n", key, allDice)

				for _, dieRequest := range allDice {
					match := wellFormattedKeyRegex.FindStringSubmatch(strings.TrimSpace(dieRequest))
					if len(match) < 1 {
						renderError(HttpError{
							HttpCode:  400,
							ErrorCode: "MALFORMED_INPUT",
							Msg:       "Invalid input! Please make sure that your input is well-formed as xDy, where 'x' is the number of dice and 'y' is the number of sides on the dice",
						},
							writer)
						return
					}
					numOfDie, err := strconv.Atoi(match[1])
					if err != nil {
						renderError(HttpError{
							HttpCode:  400,
							ErrorCode: "MALFORMED_INPUT",
							Msg:       fmt.Sprintf("Error in parsing number of die: %s\n", err.Error()),
						},
							writer)
					}
					diceRequests[match[2]] = numOfDie
				}

				var totalRollValue = 0
				dieRolls := make(map[string][]int)

				for dieSide, numOfDie := range diceRequests {
					dieSideVal, _ := strconv.Atoi(dieSide)
					dieRollKey := fmt.Sprintf("d%d", dieSideVal)

					dieRolls[dieRollKey] = make([]int, 0)

					for i := range numOfDie {
						fmt.Printf("Rolling %d of %d for %s\n", i, numOfDie, dieRollKey)
						rollResult := rand.IntN(dieSideVal) + 1 // rand results returns a random int n, 0 <= n < dieSideVal
						dieRolls[dieRollKey] = append(dieRolls[dieRollKey], rollResult)
						totalRollValue += rollResult
					}
				}

				result := HttpDiceRollResult{
					Total:     totalRollValue,
					Breakdown: dieRolls,
				}

				json.NewEncoder(writer).Encode(result)
			}
		}
	}

	generateTokenFunc := func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode(HttpToken{helpers.GenerateToken(time.Now().UTC()).EncodedVal})
	}

	urlRouter.HandleFunc("/diceroller", rollHandlerFunc).Methods(http.MethodGet)
	urlRouter.HandleFunc("/token", generateTokenFunc).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", urlRouter))
}

type HttpDiceRollResult struct {
	Total     int              `json:"total"`
	Breakdown map[string][]int `json:"breakdown"`
}

type HttpToken struct {
	Token string `json:"token"`
}

type HttpError struct {
	HttpCode  int    `json:"http_code"`
	Msg       string `json:"error_message"`
	ErrorCode string `json:"error_code"`
}

func renderError(httpError HttpError, writer http.ResponseWriter) {
	writer.WriteHeader(httpError.HttpCode)
	json.NewEncoder(writer).Encode(httpError)
	fmt.Println(httpError)
}
