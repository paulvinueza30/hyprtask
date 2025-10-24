package viewmodel

import "github.com/paulvinueza30/hyprtask/internal/logger"

func (v *ViewModel) handleAction(a ViewAction) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.setSortKey(a.NewSortKey)
	v.setSortOrder(a.NewSortOrder)
}
func (v *ViewModel) setSortKey(sk SortKey) {
	if _, ok := validSortKeys[sk]; !ok {
		logger.Log.Warn("invalid sort key entered", "sort key", sk)
		return
	}
	v.viewOptions.SortKey = sk
}

func (v *ViewModel) setSortOrder(so SortOrder) {
	if _, ok := validSortOrders[so]; !ok {
		logger.Log.Warn("invalid sort order entered", "sort order", so)
		return
	}
	v.viewOptions.SortOrder = so
}
