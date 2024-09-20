package http

import (
	accesstoken "diceroller_server/access_token"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// +HttpDiceRollResult+ is the web-only response data format for the dice roll handler
// This should not be used internally
type HttpDiceRollResult struct {
	Total     int              `json:"total"`
	Breakdown map[string][]int `json:"breakdown"`
}

func RollDice(writer http.ResponseWriter, request *http.Request) {
	wellFormattedKeyRegex := regexp.MustCompile(`(\d)+d(\d)+`)

	diceRequests := make(map[string]int)

	params := request.URL.Query()

	if !params.Has("token") || !params.Has("dice") {
		renderError(HttpError{
			HttpCode:  400,
			ErrorCode: "MALFORMED_INPUT",
			Msg:       "Request must have 'token' and 'dice'",
		}, writer)
		return
	}

	diceRollResult := HttpDiceRollResult{}

	for key, value := range params {
		// Handle auth token
		if key == "token" {
			accessToken, err := accesstoken.DecodeToken(value[0])
			if err != nil {
				renderError(HttpError{
					HttpCode:  400,
					ErrorCode: "MALFORMED_INPUT",
					Msg:       err.Error(),
				}, writer)
				return
			}
			if !accessToken.IsValid(time.Now().UTC(), 3*time.Minute) {
				renderError(HttpError{
					HttpCode:  403,
					ErrorCode: "INVALD_TOKEN",
					Msg:       "Your token is invalid, please generate a new one",
				}, writer)
				return
			}
		}

		// Handle dice rolling
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
					return
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
			diceRollResult.Total = totalRollValue
			diceRollResult.Breakdown = dieRolls
		}
	}

	json.NewEncoder(writer).Encode(diceRollResult)
}
