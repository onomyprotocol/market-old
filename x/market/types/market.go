package types

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: change value
const r uint64 = 1

func NewMarket(pairs []*Pair, pairPlusCoins []*PairPlusCoin) *Market {
	// TODO: should we validate that all pairs are unique?

	orderQueues := make([]*OrderQueue, len(pairs))
	drops := make([]*Drop, len(pairs))
	for i, pair := range pairs {
		orderQueues[i] = &OrderQueue{
			Pair:   pair,
			Orders: make([]*Order, 0),
		}
		drops[i] = &Drop{
			Pair:  pair,
			Value: 0,
		}
	}

	books := make([]*Book, len(pairPlusCoins))
	bonds := make([]*Bond, len(pairPlusCoins))
	for i, pairPlusCoin := range pairPlusCoins {
		books[i] = &Book{
			PairPlusCoin: pairPlusCoin,
			Positions:    make([]*Position, 0),
		}
		bonds[i] = &Bond{
			PairPlusCoin: pairPlusCoin,
			Value:        0,
		}
	}

	return &Market{
		OrderQueues: orderQueues,
		Drops:       drops,
		Books:       books,
		Bonds:       bonds,
	}
}

func (m *Market) BondAskAmount(bondAskValue, bondBidValue, bidAmount uint64) uint64 {
	return bondAskValue * bidAmount / bondBidValue
}

// TODO: does it match the spec?
func (m *Market) Stronger(pair *Pair) *types.Coin {
	a := pair.GetA()
	b := pair.GetB()

	if a.IsGTE(*b) {
		return a
	}

	return b
}

// TODO: does it match the spec?
func (m *Market) Weaker(pair *Pair) *types.Coin {
	a := pair.GetA()
	b := pair.GetB()

	if a.IsLT(*b) {
		return a
	}

	return b
}

func (m *Market) MaxBondBid(bookAsk *Book, bondAsk, bondBid *Bond) (uint64, error) {
	positions := bookAsk.GetPositions()
	if len(positions) == 0 {
		return 0, errors.ErrLogic
	}

	exchangeRateFinal := positions[0].GetExchangeRate().Dec.BigInt().Uint64()
	res := bondAsk.GetValue() - (bondBid.GetValue()^2)*(exchangeRateFinal^2)/bondAsk.GetValue()

	return res, nil
}

// TODO: does it match the spec?
func (m *Market) Reconcile(bondAsk, bondBid *Bond, bookAsk, bookBid *Book) error {
	bondAskUpdate := bondAsk
	bondBidUpdate := bondBid

	bookAskUpdate := bookAsk
	bookAskUpdatePositions := bookAskUpdate.GetPositions()
	if len(bookAskUpdatePositions) == 0 {
		return errors.ErrLogic
	}

	bookBidUpdate := bookBid
	bookBidUpdatePositions := bookBidUpdate.GetPositions()
	if len(bookBidUpdatePositions) == 0 {
		return errors.ErrLogic
	}

	bookAskPositions := bookAsk.GetPositions()
	for i := range bookAskPositions {
		if bookAskPositions[i].GetExchangeRate().Dec.GTE(
			types.NewDecFromBigInt(big.NewInt(int64(bondAsk.GetValue() / bondBid.GetValue()))),
		) {
			maxBondBid, err := m.MaxBondBid(bookAskUpdate, bondAskUpdate, bondBidUpdate)
			if err != nil {
				return errors.ErrLogic
			}

			if bookAskUpdatePositions[0].GetAmount().Amount.BigInt().Uint64() >= maxBondBid {
				bondAskUpdate.Value = bondAsk.GetValue() - maxBondBid
				bondBidUpdate.Value = bondBid.GetValue() + maxBondBid

				amount := types.NewInt64DecCoin(
					bookAskUpdatePositions[0].GetAmount().GetDenom(),
					int64(bookAskUpdatePositions[0].GetAmount().Amount.BigInt().Uint64()-maxBondBid),
				)

				bookAskUpdatePositions[0] = &Position{
					Amount:       &amount,
					ExchangeRate: bookAskUpdatePositions[0].GetExchangeRate(),
				}

				break
			} else {
				bondAskUpdate.Value -= maxBondBid
				bondBidUpdate.Value += maxBondBid
				bookAskUpdatePositions = bookAskPositions[1:]
			}
		} else {
			for j := range bookBidUpdatePositions {
				if bookBidUpdatePositions[j].GetExchangeRate().Dec.LTE(
					types.NewDecFromBigInt(big.NewInt(int64(bondBidUpdate.GetValue() / bondAskUpdate.GetValue()))),
				) {
					maxBondBid, err := m.MaxBondBid(bookBidUpdate, bondBidUpdate, bondAskUpdate)
					if err != nil {
						return errors.ErrLogic
					}

					exchangeRate := bookBidUpdatePositions[0].GetExchangeRate()

					if bookBidUpdatePositions[0].GetAmount().Amount.BigInt().Uint64() >= maxBondBid {
						bondBidUpdate.Value = bondBid.GetValue() - maxBondBid
						bondAskUpdate.Value += maxBondBid * exchangeRate.Dec.BigInt().Uint64()

						amount := types.NewInt64DecCoin(
							bookBidUpdatePositions[0].GetAmount().GetDenom(),
							int64(bookBidUpdatePositions[0].GetAmount().Amount.BigInt().Uint64()-maxBondBid),
						)

						bookBidUpdatePositions[0] = &Position{
							Amount:       &amount,
							ExchangeRate: exchangeRate,
						}
					} else {
						bondBidUpdate.Value -= maxBondBid
						bondAskUpdate.Value += maxBondBid * exchangeRate.Dec.BigInt().Uint64()
						bookAskUpdatePositions = bookAskPositions[1:]
					}
				} else {
					*bondAsk = *bondAskUpdate
					*bondBid = *bondBidUpdate
					*bookAsk = *bookAskUpdate
					*bookBid = *bookBidUpdate

					return nil
				}
			}
		}
	}

	return nil
}

func (m *Market) SubmitOrder(order *Order, pair *Pair) error {
	switch o := order.GetOrder().(type) {
	case *Order_LimitOrder:
		bid := o.LimitOrder.GetBid()
		ask := o.LimitOrder.GetAsk()
		if bid.IsEqual(*ask) {
			return errors.ErrLogic
		}

		orderQueue, err := getOrderQueueByPair(m.GetOrderQueues(), pair)
		if err != nil {
			return errors.ErrLogic
		}

		orderQueue.Orders = append(orderQueue.GetOrders(), order)

		return nil
	case *Order_MarketOrder:
		bid := o.MarketOrder.GetBid()
		ask := o.MarketOrder.GetAsk()
		if bid.IsEqual(*ask) {
			return errors.ErrLogic
		}

		orderQueue, err := getOrderQueueByPair(m.GetOrderQueues(), pair)
		if err != nil {
			return errors.ErrLogic
		}

		orderQueue.Orders = append(orderQueue.GetOrders(), order)

		return nil
	default:
		return errors.ErrLogic
	}
}

// TODO: does it match the spec?
func (m *Market) Provision(pair *Pair) error {
	bond, err := getBondByPair(m.GetBonds(), pair)
	if err != nil {
		return err
	}

	drop, err := getDropByPair(m.GetDrops(), pair)
	if err != nil {
		return err
	}

	c := m.Weaker(pair)
	d := m.Stronger(pair)

	d.Amount.Add(d.Amount.MulRaw(int64(r) / int64(bond.Value)))
	c.Amount.AddRaw(int64(r))

	bond.GetPairPlusCoin().Pair = pair
	drop.Value += r

	return nil
}

// TODO: does it match the spec?
func (m *Market) Liquidate(pair *Pair) error {
	drop, err := getDropByPair(m.GetDrops(), pair)
	if err != nil {
		return err
	}

	if !(r < drop.GetValue()) {
		return errors.ErrLogic
	}

	bond, err := getBondByPair(m.GetBonds(), pair)
	if err != nil {
		return err
	}

	c := m.Weaker(pair)
	d := m.Stronger(pair)

	d.Amount.Sub(d.Amount.MulRaw(int64(r) / int64(bond.Value)))
	c.Amount.SubRaw(int64(r))

	bond.GetPairPlusCoin().Pair = pair
	drop.Value -= r

	return nil
}

// TODO: does it match the spec?
func (m *Market) ProcessOrder(pair *Pair) error {
	orderQueue, err := getOrderQueueByPair(m.GetOrderQueues(), pair)
	if err != nil {
		return err
	}

	orders := orderQueue.GetOrders()
	if len(orders) == 0 {
		return err
	}
	order := orders[0]

	// bookAsk == books[pair][o.ask]
	// bookBid == books[pair][o.bid]
	book, err := getBookByPair(m.GetBooks(), pair)
	if err != nil {
		return errors.ErrLogic
	}
	// TODO: how to get bookAsk and bookBid from books?
	_ = book
	var bookAsk, bookBid *Book

	bookAskPositions := bookAsk.GetPositions()
	if len(bookAskPositions) == 0 {
		return errors.ErrLogic
	}

	bookBidPositions := bookBid.GetPositions()
	if len(bookBidPositions) == 0 {
		return errors.ErrLogic
	}

	// bondAsk == bonds[pair][o.ask]
	// bondBid == bonds[pair][o.bid]
	bond, err := getBondByPair(m.GetBonds(), pair)
	if err != nil {
		return errors.ErrLogic
	}
	// TODO: how to get bondAsk and bondBid from bonds?
	_ = bond
	var bondAsk, bondBid *Bond

	var orderAmount *types.DecCoin
	var exchangeRate *types.DecProto
	switch o := order.GetOrder().(type) {
	case *Order_LimitOrder:
		orderAmount = o.LimitOrder.GetAmount()
		exchangeRate = o.LimitOrder.GetExchangeRate()
	case *Order_MarketOrder:
		orderAmount = o.MarketOrder.GetAmount()
		exchangeRate = o.MarketOrder.GetExchangeRate()
	default:
		return errors.ErrLogic
	}

	// TODO: unused
	maxBondBid, err := m.MaxBondBid(bookAsk, bondAsk, bondBid)
	if err != nil {
		return errors.ErrLogic
	}
	_ = maxBondBid

	// Process
	if exchangeRate.Dec.GTE(bookBidPositions[0].GetExchangeRate().Dec) {
		for i := len(bookBidPositions); i > 0; i-- {
			if exchangeRate.Dec.LT(bookBidPositions[i].GetExchangeRate().Dec) {
				insertPositionAt(bookBidPositions, i, &Position{
					Amount:       orderAmount,
					ExchangeRate: exchangeRate,
				})

				break
			}
		}
	} else {
		insertPositionAt(bookBidPositions, 0, &Position{
			Amount:       orderAmount,
			ExchangeRate: exchangeRate,
		})
	}

	if err = m.Reconcile(bondAsk, bondBid, bookAsk, bookBid); err != nil {
		return errors.ErrLogic
	}

	bookAsk.PairPlusCoin = bondAsk.GetPairPlusCoin()
	bookBid.PairPlusCoin = bondBid.GetPairPlusCoin()
	bondAsk.PairPlusCoin = bookAsk.GetPairPlusCoin()
	bondBid.PairPlusCoin = bookBid.GetPairPlusCoin()

	return nil
}

func getOrderQueueByPair(orderQueues []*OrderQueue, pair *Pair) (*OrderQueue, error) {
	for _, orderQueue := range orderQueues {
		// TODO: can we compare this way?
		if orderQueue.GetPair() == pair {
			return orderQueue, nil
		}
	}

	return nil, errors.ErrLogic
}

func getDropByPair(drops []*Drop, pair *Pair) (*Drop, error) {
	for _, drop := range drops {
		// TODO: can we compare this way?
		if drop.GetPair() == pair {
			return drop, nil
		}
	}

	return nil, errors.ErrLogic
}

func getBookByPair(books []*Book, pair *Pair) (*Book, error) {
	for _, book := range books {
		// TODO: can we compare this way?
		if book.GetPairPlusCoin().GetPair() == pair {
			return book, nil
		}
	}

	return nil, errors.ErrLogic
}

func getBondByPair(bonds []*Bond, pair *Pair) (*Bond, error) {
	for _, bond := range bonds {
		// TODO: can we compare this way?
		if bond.GetPairPlusCoin().GetPair() == pair {
			return bond, nil
		}
	}

	return nil, errors.ErrLogic
}

func insertPositionAt(positions []*Position, i int, position *Position) {
	if len(positions) == i {
		positions = append(positions, position)
		return
	}

	positions = append(positions[:i+1], positions[i:]...)
	positions[i] = position
}
