package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

type UserId int
type UserMap map[UserId]*User

type Users struct {
	userMap UserMap

	allAges     []int
	allPayments []int
}

type User struct {
	id UserId

	ageIndex       int
	paymentIndexes []int
}

func AverageAge(users Users) float64 {
	average, count := 0.0, 0.0
	for _, age := range users.allAges {
		count++
		average += (float64(age) - average) / count
	}

	return average
}

func AveragePaymentAmount(users Users) float64 {
	average, count := 0.0, 0.0
	for _, p := range users.allPayments {
		count++
		amount := float64(p)
		average += (amount - average) / count
	}

	return average / 100
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(users Users) float64 {
	mean := AveragePaymentAmount(users)

	squaredDiffs, count := 0.0, 0.0
	for _, p := range users.allPayments {
		count++
		amount := float64(p / 100)
		diff := amount - mean
		squaredDiffs += diff * diff
	}

	return math.Sqrt(squaredDiffs / count)
}

func LoadData() Users {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	users := Users{
		userMap: make(UserMap, len(userLines)),
	}

	for _, line := range userLines {
		id, _ := strconv.Atoi(line[0])

		age, _ := strconv.Atoi(line[2])
		users.allAges = append(users.allAges, age)

		users.userMap[UserId(id)] = &User{
			id:       UserId(id),
			ageIndex: len(users.allAges) - 1,
		}
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	for _, line := range paymentLines {
		userID, _ := strconv.Atoi(line[2])
		paymentCents, _ := strconv.Atoi(line[0])

		users.allPayments = append(users.allPayments, paymentCents)
		id := UserId(userID)

		user := users.userMap[id]
		user.paymentIndexes = append(user.paymentIndexes, len(users.allPayments)-1)
	}

	return users
}
