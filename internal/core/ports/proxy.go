package ports

import "context"

// Proxy a proxy interface used to as MITM between Opensearch Dashboards and Opensearch
type Proxy interface {
	Start(ctx context.Context) error
	AddStorage(storage Storage)
}
