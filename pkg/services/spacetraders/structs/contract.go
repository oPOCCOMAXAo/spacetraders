package ss

import "time"

// Contract is a faction mission with a deadline and delivery terms.
type Contract struct {
	ID                string        `json:"id"`
	Type              ContractType  `json:"type"`
	FactionSymbol     string        `json:"factionSymbol"`
	Terms             ContractTerms `json:"terms"`
	Accepted          bool          `json:"accepted"`
	Fulfilled         bool          `json:"fulfilled"`
	DeadlineToAccept  time.Time     `json:"deadlineToAccept"`
	DeadlineToFulfill time.Time     `json:"deadlineToFulfill"`
}

type ContractTerms struct {
	Deadline time.Time             `json:"deadline"`
	Payment  ContractPayment       `json:"payment"`
	Deliver  []ContractDeliverGood `json:"deliver"`
}

type ContractPayment struct {
	OnAccepted  int64 `json:"onAccepted"`
	OnFulfilled int64 `json:"onFulfilled"`
}

type ContractDeliverGood struct {
	TradeSymbol       string `json:"tradeSymbol"`
	DestinationSymbol string `json:"destinationSymbol"`
	UnitsRequired     int    `json:"unitsRequired"`
	UnitsFulfilled    int    `json:"unitsFulfilled"`
}
