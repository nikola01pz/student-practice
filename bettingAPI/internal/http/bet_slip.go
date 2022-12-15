package http

import "bettingAPI/internal/mysql"

func isStakeSufficient(stake int) bool {
	return stake >= 2
}

func hasSufficientFunds(balance int) bool {
	return balance >= 2
}

func (h *handler) isPayoutLimitReached(betSlip mysql.BetSlipRequest) (float32, bool) {
	offerTips := h.db.GetOfferTipCoefficients(betSlip.Bets)
	var payout float32 = 0.0
	for i := range offerTips {
		payout += betSlip.Stake * float32(offerTips[i].Coefficient)
	}
	var maxPayout float32 = 10000
	return payout, payout > maxPayout
}

func hasOnlyOneTipPerOffer(bets []mysql.Bet) bool {
	for i := range bets {
		for j := range bets {
			if i != j {
				if bets[i].OfferID == bets[j].OfferID {
					return false
				}
			}
		}
	}
	return true
}
