package parser

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestCSVParser_GenerateFile(t *testing.T) {
	// Create a new CSV writer
	writer := csv.NewWriter(os.Stdout)

	// Write the headers to the csv file
	writer.Write([]string{"accountId", "date", "amount"})

	// Create a slices of accountId's
	accountIds := []string{"acc1", "acc2"}

	// Write some random data to the csv file
	for i := 0; i < 20; i++ {
		accountId := accountIds[rand.Intn(len(accountIds))]

		date := randomDate()

		// Generate a random positive int64, then subtract 5000 to get both positive and negative numbers
		amount := rand.Int63n(10000) - 5000

		writer.Write([]string{accountId, date.Format(time.RFC3339), fmt.Sprintf("%+d", amount)})
	}

	// Write any buffered data to the underlying writer (standard output).
	writer.Flush()
}

func randomDate() time.Time {
	daysInTwoMonths := 30 * 2 // Approximation of 2 months in days
	durationInHours := time.Duration(rand.Int63n(int64(daysInTwoMonths*24))) * time.Hour

	pastMonth := time.Now().Add(-30 * 24 * time.Hour)

	return pastMonth.Add(durationInHours)
}
