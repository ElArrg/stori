package parser

import (
	"context"
	encodingCsv "encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"

	"github.com/elarrg/stori/ledger/internal/models"
)

type CSVParser struct {
	expectedFields []string
}

type record struct {
	line int
	data []string
}

func NewCSVParser(expectedFields []string) *CSVParser {
	return &CSVParser{
		expectedFields: expectedFields,
	}
}

func (c *CSVParser) Parse(ctx context.Context, r io.Reader) (records []models.Transaction, err error) {
	csv := encodingCsv.NewReader(r)
	fieldsPosition, err := c.mapFieldPosition(csv)
	if err != nil {
		return nil, err
	}

	linesCounter := 2
	for {
		var data []string
		data, err = csv.Read()
		if err != nil {
			if err == io.EOF {
				// finished reading file
				break
			}

			return nil, fmt.Errorf("[trans-csv-parser]: error parsing record, %v", err)
		}

		row := &record{
			line: linesCounter,
			data: data,
		}

		var trans *models.Transaction
		trans, err = c.mapRecordToModel(row, fieldsPosition)
		if err != nil {
			return nil, err
		}

		records = append(records, *trans)
		linesCounter++
	}

	return records, nil
}

func (c *CSVParser) mapFieldPosition(csv *encodingCsv.Reader) (map[string]int, error) {
	// Get the first line where field names are specified
	headers, err := csv.Read()
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("[trans-csv-parser]: file is empty")
		}

		return nil, fmt.Errorf("[trans-csv-parser]: couldn't read headers line, %v", err)
	}

	// initialize required missing fields
	missingFields := make(map[string]bool, len(c.expectedFields))
	for _, field := range c.expectedFields {
		if v := missingFields[field]; v {
			return nil, fmt.Errorf("[trans-csv-parser]: required field %s is duplicated", field)
		}
		missingFields[field] = true
	}

	// map fields to index position
	fieldsPosition := make(map[string]int, len(headers))
	for i, field := range headers {
		if _, ok := fieldsPosition[field]; ok {
			return nil, fmt.Errorf("[trans-csv-parser]: duplicated field '%s' at column %d", field, i+1)
		}

		fieldsPosition[field] = i

		// if is a required field remove it from missing
		if v := missingFields[field]; v {
			delete(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return nil, fmt.Errorf("[trans-csv-parser]: there are missing required fields %v", missingFields)
	}

	return fieldsPosition, nil
}

func (c *CSVParser) mapRecordToModel(r *record, fieldPosition map[string]int) (*models.Transaction, error) {
	trans := models.Transaction{}

	trans.ID = uuid.NewString() // assign new ID
	trans.AccountID = r.data[(fieldPosition)["accountId"]]

	err := trans.Date.UnmarshalText([]byte(r.data[(fieldPosition)["date"]]))
	if err != nil {
		return nil, fmt.Errorf("[trans-csv-parser] (row: %d): couldn't parse date, %v", r.line, err)
	}

	trans.Year, trans.Month, _ = trans.Date.Date()

	amount, err := strconv.ParseInt(r.data[(fieldPosition)["amount"]], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("[trans-csv-parser] (row: %d): couldnt parse amount, %v", r.line, err)
	}

	trans.Amount = amount
	if amount >= 0 {
		trans.Type = models.CreditTransactionType
	} else {
		trans.Type = models.DebitTransactionType
	}

	return &trans, nil
}

/*
// ParseErrHandler a function that handles parsing errors.
// It takes an error as input and returns a boolean value indicating whether the error was handled successfully.
// A ParseErrHandler function should implement a custom logic to handle parsing errors in a specific way.
// It should return true if the error was handled successfully and false otherwise.
// If an error is not handled successfully it will stop the processing.
type ParseErrHandler func(err error) bool

type ProcessBatchFunc func(batch []*models.Transaction) bool

type WorkerBatchResultsChan <-chan []*models.Transaction

type recordsBatchChan <-chan []*record

func (c *CSVParser) ParseConcurrent(ctx context.Context, r io.Reader, processFn ProcessBatchFunc, errHandler ParseErrHandler, workers int, batchSize int) {
	ctx, cancel := context.WithCancel(ctx)
	csv := encodingCsv.NewReader(r)
	errsCh := make(chan error)
	defer close(errsCh)

	fieldsPositionMap, err := c.mapFieldPosition(csv)
	if err != nil {
		errsCh <- err
	}
	batchesCh := c.readRows(ctx, csv, fieldsPositionMap, batchSize, errsCh)

	workerChs := make([]WorkerBatchResultsChan, workers)
	for i := 0; i < workers; i++ {
		workerChs[i] = c.parseRows(ctx, fieldsPositionMap, batchesCh, errsCh)
	}

	c.processRecords(ctx, cancel, processFn, errHandler, errsCh, workerChs...)
}

func (c *CSVParser) readRows(ctx context.Context, csv *encodingCsv.Reader, batchSize int, errsCh chan<- error) recordsBatchChan {
	rowCount := 2 // the first record starts at this line number
	batchesCh := make(chan []*record)

	go func() {
		defer close(batchesCh)

		currentBatch := make([]*record, 0, batchSize)
		for {
			select {
			case <-ctx.Done(): // processing should be cancelled now!
				return

			default:
				r, err := csv.Read()
				if err != nil {
					if errors.Is(err, io.EOF) {
						batchesCh <- currentBatch
						return
					}

					errsCh <- err // send err and send the partial record
				}

				row := record{
					line: rowCount,
					data: r,
				}
				rowCount++
				currentBatch = append(currentBatch, &row)

				if len(currentBatch) == batchSize {
					batchesCh <- currentBatch
					currentBatch = make([]*record, 0, batchSize)
				}
			}
		}
	}()

	return batchesCh
}

func (c *CSVParser) parseRows(ctx context.Context, fieldPosition map[string]int, batchesCh recordsBatchChan, errsCh chan<- error) WorkerBatchResultsChan {
	transBatches := make(chan []*models.Transaction)

	go func() {
		defer close(transBatches)

		select {
		case <-ctx.Done(): // process should be cancelled now
			return

		default:
			for b := range batchesCh {
				batch := make([]*models.Transaction, 0, len(b))

				for _, r := range b {
					trans, err := c.mapRecordToModel(r, fieldPosition)
					if err != nil {
						errsCh <- err
					}
					batch = append(batch, trans)
				}

				transBatches <- batch
			}
		}
	}()

	return transBatches
}

func (c *CSVParser) processRecords(ctx context.Context, cancelCtx context.CancelFunc, fn ProcessBatchFunc, errFn ParseErrHandler, errsCh <-chan error, parsedBatchesCh ...WorkerBatchResultsChan) {
	wg := sync.WaitGroup{}
	wg.Add(len(parsedBatchesCh))

	for _, p := range parsedBatchesCh {
		go func(ch WorkerBatchResultsChan) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return

			case err := <-errsCh:
				if !errFn(err) {
					cancelCtx()
				}

			case b := <-ch:
				if !fn(b) {
					cancelCtx()
				}
			}
		}(p)
	}

	wg.Wait()
}
*/
