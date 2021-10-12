package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

type UserID int
type UserMap map[UserID]*User

type Users struct {
	userMap UserMap

	allAges     []int
	allPayments []uint32
}

type User struct {
	id UserID

	// I keep these indices around to keep the general
	// datastructure somewhat useful.
	ageIndex       int
	paymentIndexes []int
}

func AverageAge(users Users) float64 {
	avg0, avg1, avg2, avg3 := 0, 0, 0, 0

	ages := users.allAges
	count := len(ages)
	limit := count - 3

	i := 0
	for ; i < limit; i += 4 {
		avg3 += ages[i+3]
		avg2 += ages[i+2]
		avg1 += ages[i+1]
		avg0 += ages[i]
	}

	for ; i < len(ages); i++ {
		avg0 += ages[i]
	}

	return float64(avg0+avg1+avg2+avg3) / float64(count)
}

func AveragePaymentAmount(users Users) float64 {
	avg0, avg1, avg2, avg3 := uint64(0), uint64(0), uint64(0), uint64(0)

	payments := users.allPayments
	count := len(payments)
	limit := count - 3

	i := 0
	for ; i < limit; i += 4 {
		avg3 += uint64(payments[i+3])
		avg2 += uint64(payments[i+2])
		avg1 += uint64(payments[i+1])
		avg0 += uint64(payments[i])
	}

	for ; i < len(payments); i++ {
		avg0 += uint64(payments[i])
	}

	return (float64(avg0 + avg1 + avg2 + avg3)) / (float64(count * 100))
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(users Users) float64 {
	mean := AveragePaymentAmount(users)

	count := len(users.allPayments)
	limit := count - 1

	squaredDiffs := 0.0
	squaredDiffs2 := 0.0

	payments := users.allPayments

	i := 0
	for ; i < limit; i += 2 {
		p1 := payments[i]
		p2 := payments[i+1]

		amount1 := float64(p1) * .01
		diff1 := amount1 - mean
		squaredDiffs += diff1 * diff1

		amount2 := float64(p2) * .01
		diff2 := amount2 - mean
		squaredDiffs2 += diff2 * diff2

	}

	for ; i < count; i++ {
		p1 := payments[i]

		amount1 := float64(p1) * .01
		diff1 := amount1 - mean
		squaredDiffs += diff1 * diff1
	}

	return math.Sqrt((squaredDiffs + squaredDiffs2) / float64(count))
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
		users.allAges = append(users.allAges, int(age))

		users.userMap[UserID(id)] = &User{
			id:       UserID(id),
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

		users.allPayments = append(users.allPayments, uint32(paymentCents))
		id := UserID(userID)

		user := users.userMap[id]
		user.paymentIndexes = append(user.paymentIndexes, len(users.allPayments)-1)
	}

	return users
}
