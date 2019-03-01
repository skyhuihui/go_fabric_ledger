package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/chainHero/heroes-service/blockchain/blockdata"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric/protos/common"
	futils "github.com/hyperledger/fabric/protos/utils"
)

// 查询详情
func (setup *FabricSetup) QueryInfo(ledger ledger.Client) *fab.BlockchainInfoResponse {
	info, err := ledger.QueryInfo()
	if err != nil {
		fmt.Printf("failed to query info: %s\n", err)
	}
	return info
}

// 查询区块高度
func (setup *FabricSetup) QueryBlockHeight(ledger ledger.Client) int64 {
	info, err := ledger.QueryInfo()
	if err != nil {
		fmt.Printf("failed to query info: %s\n", err)
	}
	return int64(info.BCI.Height)
}

// 查询指定高度 区块数据
func (setup *FabricSetup) QueryBlock(ledger ledger.Client, blockBumber int64) *blockdata.Block {
	blo := &blockdata.Block{}
	blcokHeader := &common.BlockHeader{}
	block, err := ledger.QueryBlock(uint64(blockBumber))
	blcokHeader = (*common.BlockHeader)(block.Header)
	if err != nil {
	}
	blo.BlockHash = hex.EncodeToString(blcokHeader.Hash()) //区块hash
	blo.TransactionNumber = len(block.Data.Data)
	//此处应该遍历block.Data.Data
	var transaction = make([]*blockdata.ChainTransaction, 0)
	for _, data := range block.Data.Data {
		env, err := futils.GetEnvelopeFromBlock(data)
		chainTransaction, err := blockdata.EnvelopeToTrasaction((*common.Envelope)(env))
		chainTransaction.Height = int64(block.Header.Number) //区块高度
		blo.Timestamp = chainTransaction.Timestamp
		blo.Height = chainTransaction.Height
		if chainTransaction.TxID != "" {
			transaction = append(transaction, chainTransaction)
		}
		if err != nil {
			continue
		}
	}
	blo.Transaction = transaction
	return blo
}

//  通过hash查询账本
func (setup *FabricSetup) QueryBlockByHash(ledger ledger.Client, hash string) *blockdata.Block {
	query, err := hex.DecodeString(hash)
	block, err := ledger.QueryBlockByHash(query)
	if block == nil {
		return nil
	}
	blo := &blockdata.Block{}
	blcokHeader := &common.BlockHeader{}
	fmt.Println("查询区块hash:")
	fmt.Println(hash)
	blcokHeader = (*common.BlockHeader)(block.Header)
	if err != nil {
	}
	blo.BlockHash = hex.EncodeToString(blcokHeader.Hash()) //区块hash
	blo.TransactionNumber = len(block.Data.Data)
	//此处应该遍历block.Data.Data
	var transaction = make([]*blockdata.ChainTransaction, 0)
	for _, data := range block.Data.Data {
		env, err := futils.GetEnvelopeFromBlock(data)
		chainTransaction, err := blockdata.EnvelopeToTrasaction((*common.Envelope)(env))
		chainTransaction.Height = int64(block.Header.Number) //区块高度
		blo.Timestamp = chainTransaction.Timestamp
		blo.Height = chainTransaction.Height

		transaction = append(transaction, chainTransaction)
		if err != nil {
			continue
		}
	}
	blo.Transaction = transaction
	return blo
}

//  通过交易id查询账本  缺少交易时间， 通道id
func (setup *FabricSetup) QueryBlockByTx(ledger ledger.Client, txid string) {

	block, err := ledger.QueryBlockByTxID(fab.TransactionID(txid))
	fmt.Println(block)
	if err != nil {
		return
	}
	fmt.Println(hex.EncodeToString(block.Header.DataHash))     //后一个区块hash
	fmt.Println(hex.EncodeToString(block.Header.PreviousHash)) //前一个区块hash

	//此处应该遍历block.Data.Data？
	for _, data := range block.Data.Data {
		env, err := futils.GetEnvelopeFromBlock(data)
		chainTransaction, err := blockdata.EnvelopeToTrasaction((*common.Envelope)(env))
		chainTransaction.Height = int64(block.Header.Number) //区块高度
		fmt.Println(chainTransaction.Height)
		fmt.Println(chainTransaction.Chaincode) //链码名称
		fmt.Println(chainTransaction.Method)
		fmt.Println(chainTransaction.TxID) //交易id
		fmt.Println(chainTransaction.CreatedFlag)
		for _, args := range chainTransaction.TxArgs {
			fmt.Println(string(args)) //交易数据
		}
		if err != nil {
			continue
		}
	}
}

//  获取通道配置
func (setup *FabricSetup) QueryConfig(ledger ledger.Client) {

	cfg, err := ledger.QueryConfig()
	if err != nil {
		fmt.Printf("failed to query config: %s\n", err)
	}

	if cfg != nil {
		fmt.Println("Retrieved channel configuration")
	}
	fmt.Println("cfg.Versions()   测试")
	fmt.Println(cfg.Versions().Channel.Values)
	fmt.Println("cfg.id()   测试")
	fmt.Println(cfg.ID())
	fmt.Println("cfg.BlockNumber()   测试")
	fmt.Println(cfg.BlockNumber())
	fmt.Println("cfg.Orderers()   测试")
	fmt.Println(cfg.Orderers())
	fmt.Println("cfg.AnchorPeers()   测试")
	fmt.Println(cfg.AnchorPeers())
}

//  查询交易信息
func (setup *FabricSetup) QueryTransaction(ledger ledger.Client, hash string) *blockdata.ChainTransaction {

	transaction, err := ledger.QueryTransaction(fab.TransactionID(hash))
	if err != nil {
		fmt.Printf("failed to query transaction: %s\n", err)
		return nil
	}

	chainTransaction, err := blockdata.EnvelopeToTrasaction((*common.Envelope)(transaction.TransactionEnvelope))
	if err != nil {
		return nil
	}
	return chainTransaction
}
