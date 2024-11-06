package network

import (
	"crypto"
	"fmt"
	"github.com/peng9808/BlockchainX/core"
	"github.com/sirupsen/logrus"
	"time"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	blockTime   time.Duration
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	s := &Server{
		ServerOpts:  opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	// If we dont got any processor from the server options, we going to use
	// the server as default.
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s
}

func (s *Server) Start() {
	s.initTransports()

	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Error(err)
			}
			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				logrus.Error(err)
			}
		case <-s.quitCh:
			break
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.blockTime)
	for {
		select {
		case <-ticker.C:
			s.createNewBlock()
		}
	}
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}

	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("transaction already in mempool")

		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash":           hash,
		"mempool length": s.memPool.Len(),
	}).Info("adding new tx to the mempool")

	// TODO(@anthdm): broadcast this tx to peers

	return s.memPool.Add(tx)
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
