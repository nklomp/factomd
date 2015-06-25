// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factomapi

import (
	"encoding/hex"

	"github.com/FactomProject/btcd/wire"
	"github.com/FactomProject/FactomCode/common"
	"github.com/FactomProject/FactomCode/database"
    "github.com/FactomProject/FactomCode/process"
    fct "github.com/FactomProject/factoid"
    
    
)

var (
	db     database.Db
	inMsgQ chan wire.FtmInternalMsg
)

func ChainHead(chainid string) (*common.EBlock, error) {
	h, err := atoh(chainid)
	if err != nil {
		return nil, err
	}
	c, err := db.FetchChainByHash(h)
	if err != nil {
		return nil, err
	}
	return c.NextBlock, nil
}

func CommitChain(c *common.CommitChain) error {
	m := wire.NewMsgCommitChain()
	m.CommitChain = c
	inMsgQ <- m
	return nil
}

func CommitEntry(c *common.CommitEntry) error {
	m := wire.NewMsgCommitEntry()
	m.CommitEntry = c
	inMsgQ <- m
	return nil
}

 func FactoidTX(t fct.ITransaction) error {
     m := new(wire.MsgFactoidTX)
     m.SetTransaction(t)
     inMsgQ <- m
     return nil
}

func DBlockByKeyMR(keymr string) (*common.DirectoryBlock, error) {
	key, err := atoh(keymr)
	if err != nil {
		return nil, err
	}
	r, err := db.FetchDBlockByMR(key)
	if err != nil {
		return r, err
	}
	return r, nil
}

func DBlockHead() (*common.DirectoryBlock, error) {
	_, height, err := db.FetchBlockHeightCache()
	if err != nil {
		return nil, err
	}
	block, err := db.FetchDBlockByHeight(uint32(height))
	if err != nil {
		return nil, err
	}
	block.BuildKeyMerkleRoot()
	return block, nil
}

func EBlockByKeyMR(keymr string) (*common.EBlock, error) {
	h, err := atoh(keymr)
	if err != nil {
		return nil, err
	}
	r, err := db.FetchEBlockByMR(h)
	if err != nil {
		return r, err
	}
	return r, nil
}

func ECBalance(eckey string) (uint32, error) {
	key := new([32]byte)
	if p, err := hex.DecodeString(eckey); err != nil {
		return 0, err
	} else {
		copy(key[:], p)
	}
	val, _ := process.GetEntryCreditBalance(key)
	return uint32(val), nil
}

func EntryByHash(hash string) (*common.Entry, error) {
	h, err := atoh(hash)
	if err != nil {
		return nil, err
	}
	r, err := db.FetchEntryByHash(h)
	if err != nil {
		return r, err
	}
	return r, nil
}

func RevealEntry(e *common.Entry) error {
	m := wire.NewMsgRevealEntry()
	m.Entry = e
	inMsgQ <- m
	return nil
}

func SetDB(d database.Db) {
	db = d
}

func SetInMsgQueue(q chan wire.FtmInternalMsg) {
	inMsgQ = q
}

func atoh(a string) (*common.Hash, error) {
	h := common.NewHash()
	p, err := hex.DecodeString(a)
	if err != nil {
		return h, err
	}
	h.SetBytes(p)
	return h, nil
}
