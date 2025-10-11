package viewmodel

import "github.com/paulvinueza30/hyprtask/internal/logger"

func (v *ViewModel) handleAction(a ViewAction) {
	v.mu.Lock()
	defer v.mu.Unlock()
	switch a.Type {
	case ActionSetSortKey:
		if key, ok := a.Payload.(SortKey); ok {
			v.setSortKey(key)
		}
	case ActionSetSortOrder:
		if key, ok := a.Payload.(SortOrder); ok {
			v.setSortOrder(key)
		}
	}
	v.buildDisplayData()
}
func (v *ViewModel) setSortKey(sk SortKey) {
	if _, ok := validSortKeys[sk]; !ok {
		logger.Log.Warn("invalid sort key entered", "sort key", sk)
		return
	}
	v.viewOptions.SortBy = sk
	if v.viewOptions.SortOrder == OrderNone {
		// default is ASC
		v.viewOptions.SortOrder = OrderASC
	}
}

func (v *ViewModel) setSortOrder(so SortOrder) {
	if _, ok := validSortOrders[so]; !ok {
		logger.Log.Warn("invalid sort order entered", "sort order", so)
		return
	}
	v.viewOptions.SortOrder = so
}
