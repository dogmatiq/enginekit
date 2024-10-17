package config

import "cmp"

// routeKey contains the components of a [Route] that uniquely identify it.
type routeKey struct {
	RouteType       RouteType
	MessageTypeName string
}

func routeKeyOf(r *Route) (routeKey, bool) {
	rt, rtOK := r.AsConfigured.RouteType.TryGet()
	mt, mtOK := r.AsConfigured.MessageTypeName.TryGet()
	return routeKey{rt, mt}, rtOK && mtOK
}

func (k routeKey) Compare(x routeKey) int {
	if c := cmp.Compare(k.RouteType, x.RouteType); c != 0 {
		return c
	}
	return cmp.Compare(k.MessageTypeName, x.MessageTypeName)
}
