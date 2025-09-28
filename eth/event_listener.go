package eth

import (
	"context"
	"log"
	"math/big"

	"go-indexer/abis/Vault"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventData struct {
	EventName string
	User      common.Address
	Amount    *big.Int
}

type EventListener struct {
	ctx     context.Context
	client  *ethclient.Client
	vault   *Vault.Vault
	eventCh chan<- []EventData
}

func NewEventListener(ctx context.Context, client *ethclient.Client, vault *Vault.Vault, eventCh chan<- []EventData) *EventListener {
	return &EventListener{
		ctx:     ctx,
		client:  client,
		vault:   vault,
		eventCh: eventCh,
	}
}

func (el *EventListener) Listen() {
	headers := make(chan *types.Header)
	sub, err := el.client.SubscribeNewHead(el.ctx, headers)
	if err != nil {
		log.Fatalf("Header subscription failed: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("Header subscription error: %v", err)
			return
		case header := <-headers:
			el.fetchAndSendEvents(header.Number)
		case <-el.ctx.Done():
			return
		}
	}
}

func (el *EventListener) fetchAndSendEvents(blockNum *big.Int) {
	var events []EventData

	// Deposit events
	iter, err := el.vault.FilterDeposit(nil, nil, nil)
	if err == nil {
		for iter.Next() {
			ev := iter.Event
			events = append(events, EventData{
				EventName: "Deposit",
				User:      ev.Owner,
				Amount:    ev.Assets,
			})
		}
		iter.Close()
	}
	// Withdraw events
	iterW, err := el.vault.FilterWithdraw(nil, nil, nil, nil)
	if err == nil {
		for iterW.Next() {
			ev := iterW.Event
			events = append(events, EventData{
				EventName: "Withdraw",
				User:      ev.Owner,
				Amount:    ev.Assets,
			})
		}
		iterW.Close()
	}
	if len(events) > 0 {
		el.eventCh <- events
	}
}
