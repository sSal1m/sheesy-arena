package model

import (
	"encoding/json"
	"fmt"
	"github.com/Team254/cheesy-arena/game"
	"go.etcd.io/bbolt"
	"time"
)

const FuelManual = "manual"
const FuelPlcCorrection = "plc_correction"

// FuelAlliance keeps hardware counts independent of operator corrections.
// Effective scores always use the existing game.Hub representation.
type FuelAlliance struct {
	Raw        game.Hub
	Adjustment [game.ShiftCount]int
	LastCount  int
	Sequence   uint64
	Connection uint64
	Reset      uint64
	Armed      bool
	Suspended  string
}

type FuelMetadata struct {
	AttemptId string
	Mode      string
	Red       FuelAlliance
	Blue      FuelAlliance
}

// FuelCheckpoint is a recovery record, not an officially posted MatchResult.
type FuelCheckpoint struct {
	FuelMetadata
	Match            Match
	Result           MatchResult
	Revision         uint64
	ResultAttemptId  string
	ClockStart       time.Time
	Status           string
	Phase            int
	Started          bool
	Recovered        bool
	AutoWinnerFrozen bool
	UpdatedAt        time.Time
}

type FuelCommand struct {
	RequestId string
	AttemptId string
	MatchId   int
	Alliance  string
	Phase     string
	Delta     int
	Mode      string
	Operator  string
}

type FuelOperation struct {
	ResultAttemptId string
	Mode            string
	PreviousBase    int
	BaseCount       int
	SourceSequence  uint64
	Command         FuelCommand
	Kind            string
	Actor           string
	Time            time.Time
	Before          int
	After           int
	Revision        uint64
}

var fuelStateBucket = []byte("FuelCheckpoint")
var fuelOperationBucket = []byte("FuelOperation")

func (database *Database) GetFuelCheckpoint() (*FuelCheckpoint, error) {
	var checkpoint *FuelCheckpoint
	err := database.bolt.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(fuelStateBucket)
		if bucket == nil || bucket.Get([]byte("current")) == nil {
			return nil
		}
		return json.Unmarshal(bucket.Get([]byte("current")), &checkpoint)
	})
	return checkpoint, err
}

func (database *Database) GetFuelOperation(requestId string) (*FuelOperation, error) {
	var operation *FuelOperation
	err := database.bolt.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(fuelOperationBucket)
		if bucket == nil || bucket.Get([]byte(requestId)) == nil {
			return nil
		}
		return json.Unmarshal(bucket.Get([]byte(requestId)), &operation)
	})
	return operation, err
}

// SaveFuelCheckpoint atomically records an operation and its resulting checkpoint.
// Callers serialize writes; background checkpoints must be captured AFTER obtaining
// the same operation lock, so an old snapshot cannot overwrite a durable command.
func (database *Database) SaveFuelCheckpoint(checkpoint *FuelCheckpoint, operation *FuelOperation) error {
	err := database.bolt.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(fuelStateBucket)
		if err != nil {
			return err
		}
		data, err := json.Marshal(checkpoint)
		if err != nil {
			return err
		}
		if err = bucket.Put([]byte("current"), data); err != nil {
			return err
		}
		if operation != nil {
			operations, err := tx.CreateBucketIfNotExists(fuelOperationBucket)
			if err != nil {
				return err
			}
			key := []byte(operation.Command.RequestId)
			if operations.Get(key) != nil {
				return fmt.Errorf("fuel request already recorded")
			}
			data, err = json.Marshal(operation)
			if err != nil {
				return err
			}
			return operations.Put(key, data)
		}
		return nil
	})
	if err != nil {
		return err
	}
	// The main event DB uses NoSync. Explicitly flush accepted operator actions
	// and lifecycle barriers before acknowledging them.
	return database.bolt.Sync()
}

func (database *Database) GetFuelHistory(attemptId string) ([]FuelOperation, error) {
	operations := []FuelOperation{}
	err := database.bolt.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(fuelOperationBucket)
		if bucket == nil {
			return nil
		}
		return bucket.ForEach(func(_, data []byte) error {
			var operation FuelOperation
			if err := json.Unmarshal(data, &operation); err != nil {
				return err
			}
			if operation.Command.AttemptId == attemptId || operation.ResultAttemptId == attemptId {
				operations = append(operations, operation)
			}
			return nil
		})
	})
	return operations, err
}

// ReconcileFuelMetadata preserves server-owned raw provenance when an effective
// result is edited. Missing metadata on legacy results is intentionally allowed.
func ReconcileFuelMetadata(result *MatchResult, previous *FuelMetadata) error {
	for _, score := range []*game.Score{result.RedScore, result.BlueScore} {
		if score == nil {
			return fmt.Errorf("missing score")
		}
		for _, count := range score.Hub.ShiftCounts {
			if count < 0 || count > 2147483647 {
				return fmt.Errorf("invalid FUEL count")
			}
		}
	}
	result.FuelControl = nil
	if previous == nil {
		return nil
	}
	metadata := *previous
	for index, source := range []*FuelAlliance{&metadata.Red, &metadata.Blue} {
		score := result.RedScore
		if index == 1 {
			score = result.BlueScore
		}
		source.Raw.WonAuto = score.Hub.WonAuto
		for shift, count := range score.Hub.ShiftCounts {
			source.Adjustment[shift] = count
			if metadata.Mode != FuelManual {
				source.Adjustment[shift] -= source.Raw.ShiftCounts[shift]
			}
		}
	}
	result.FuelControl = &metadata
	return nil
}
