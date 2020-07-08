package state

// Crappy implementation of verb.StateManager interface, based on GCP Cloud Datastore, using the
// skypies singleton junk. One entity per singleton.

import(
	"context"
	"fmt"

	"github.com/skypies/util/gcp/ds"
	"github.com/skypies/util/gcp/singleton"
)

type GcpStateManager struct {
	SingletonProvider singleton.SingletonProvider
}

func NewGcpStateManager(ctx context.Context, GcpProjectId string) GcpStateManager {
	p,err := ds.NewCloudDSProvider(ctx, GcpProjectId)
	if err != nil {
		panic(fmt.Errorf("NewSM: could not get a clouddsprovider (projectId=%s): %v\n", GcpProjectId, err))
	}

	gsm := GcpStateManager{
		SingletonProvider: singleton.NewProvider(p),
	}

	return gsm
}
	
func (gsm GcpStateManager)ReadState(ctx context.Context, name string, ptr interface{}) error {
	return gsm.SingletonProvider.ReadSingleton(ctx, name, nil, ptr)
}

func (gsm GcpStateManager)WriteState(ctx context.Context, name string, ptr interface{}) error {
	return gsm.SingletonProvider.WriteSingleton(ctx, name, nil, ptr)
}
