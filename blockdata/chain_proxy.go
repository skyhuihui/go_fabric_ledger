package blockdata

type Chain struct {
	Height int64
}

type Block struct {
	BlockHash         string
	Timestamp         int64
	Height            int64 `json:",string"`
	TransactionNumber int
	Transaction       []*ChainTransaction `json:"-"`
}

type ChainTransaction struct {
	Height                             int64 `json:",string"`
	Timestamp                          int64
	TxID, Chaincode, Method, ChannelId string
	CreatedFlag                        bool
	TxArgs                             [][]byte `json:"-"`
}

const (
	TxStatus_Success = 0
	TxStatus_Fail    = 1
)

type ChainTxEvents struct {
	TxID, Chaincode, Name string
	Status                int
	Payload               []byte `json:"-"`
}

type ChainBlock struct {
	Height       int64 `json:",string"`
	Hash         string
	TimeStamp    string              `json:",omitempty"`
	Transactions []*ChainTransaction `json:"-"`
	TxEvents     []*ChainTxEvents    `json:"-"`
}
