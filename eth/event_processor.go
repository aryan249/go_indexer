package eth

import (
	"go-indexer/abis/Vault"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
)

type EventProcessor struct {
	eventCh <-chan []EventData
	vault   *Vault.Vault
	auth    *bind.TransactOpts
}

func NewEventProcessor(eventCh <-chan []EventData, vault *Vault.Vault, auth *bind.TransactOpts) *EventProcessor {
	return &EventProcessor{
		eventCh: eventCh,
		vault:   vault,
		auth:    auth,
	}
}

func (ep *EventProcessor) Start() {
	aggregate := big.NewInt(0)
	for {
		select {
		case events, ok := <-ep.eventCh:
			if !ok {
				// Channel closed, process aggregate if needed and exit
				if aggregate.Cmp(big.NewInt(0)) > 0 {
					// _, err := ep.vault.SomeSetterMethod(aggregate)
				}
				return
			}
			for _, event := range events {
				log.Printf("Processing event: %s | User: %s | Amount: %s", event.EventName, event.User.Hex(), event.Amount.String())
				aggregate.Add(aggregate, event.Amount)
				// Example: call a setter or process as needed using ep.vault
				// _, err := ep.vault.SomeSetterMethod(...)
			}
		}
	}
}
