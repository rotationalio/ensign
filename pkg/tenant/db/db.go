package db

import (
	"errors"
	"sync"

	"github.com/rotationalio/ensign/pkg/tenant/config"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	mu     sync.RWMutex
	cc     *grpc.ClientConn
	client trtl.TrtlClient
)

var (
	ErrNotConnected = errors.New("not connected to trtl database")
)

func Connect(conf config.DatabaseConfig) (err error) {
	mu.Lock()
	defer mu.Unlock()

	// Check if we're already connected
	if cc != nil || client != nil {
		return nil
	}

	var endpoint string
	if endpoint, err = conf.Endpoint(); err != nil {
		return err
	}

	// Otherwise connect to the trtl database
	opts := make([]grpc.DialOption, 0, 1)
	if conf.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// TODO: connect with mtls
		return errors.New("not implemented: mtls currently not implemented")
	}

	if cc, err = grpc.Dial(endpoint, opts...); err != nil {
		return err
	}

	client = trtl.NewTrtlClient(cc)
	return nil
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()

	err := cc.Close()
	cc = nil
	client = nil
	return err
}

func connected() bool {
	return cc == nil || client == nil
}
