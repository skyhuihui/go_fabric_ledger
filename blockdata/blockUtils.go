package blockdata

import (
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	futils "github.com/hyperledger/fabric/protos/utils"
)

func EnvelopeToTrasaction(env *common.Envelope) (*ChainTransaction, error) {
	ta := &ChainTransaction{}
	//	//====Transaction====== []
	// Height                  int64 `json:",string"`
	// TxID, Chaincode, Method string
	// CreatedFlag             bool
	// TxArgs                  [][]byte `json:"-"`

	var err error
	if env == nil {
		return ta, errors.New("common.Envelope is nil")
	}
	payl := &common.Payload{}
	err = proto.Unmarshal(env.Payload, payl)
	if err != nil {
		return nil, err
	}
	tx := &pb.Transaction{}
	err = proto.Unmarshal(payl.Data, tx)
	if err != nil {
		return nil, err
	}

	taa := &pb.TransactionAction{}
	taa = tx.Actions[0]

	cap := &pb.ChaincodeActionPayload{}
	err = proto.Unmarshal(taa.Payload, cap)
	if err != nil {
		return nil, err
	}

	pPayl := &pb.ChaincodeProposalPayload{}
	proto.Unmarshal(cap.ChaincodeProposalPayload, pPayl)

	prop := &pb.Proposal{}
	pay, err := proto.Marshal(pPayl)
	if err != nil {
		return nil, err
	}

	prop.Payload = pay
	h, err := proto.Marshal(payl.Header)
	if err != nil {
		return nil, err
	}
	prop.Header = h

	invocation := &pb.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(pPayl.Input, invocation)
	if err != nil {
		return nil, err
	}

	spec := invocation.ChaincodeSpec

	// hdr := &common.Header{}
	// hdr = payl.Header
	channelHeader := &common.ChannelHeader{}
	proto.Unmarshal(payl.Header.ChannelHeader, channelHeader)

	ta.Timestamp = channelHeader.Timestamp.Seconds
	ta.ChannelId = channelHeader.ChannelId
	ta.TxID = channelHeader.TxId
	if len(spec.GetInput().GetArgs()) != 0 {
		ta.TxArgs = spec.GetInput().GetArgs()[1:]
		ta.Method = string(spec.GetInput().GetArgs()[0])
	}
	ta.Chaincode = spec.GetChaincodeId().GetName()
	ta.CreatedFlag = channelHeader.GetType() == int32(common.HeaderType_ENDORSER_TRANSACTION)
	return ta, nil
}

// blockToChainCodeEvents parses block events for chaincode events associated with individual transactions
func BlockToChainCodeEvents(block *common.Block) []*pb.ChaincodeEvent {
	if block == nil || block.Data == nil || block.Data.Data == nil || len(block.Data.Data) == 0 {
		return nil
	}
	events := make([]*pb.ChaincodeEvent, 0)
	//此处应该遍历block.Data.Data？
	for _, data := range block.Data.Data {
		event, err := GetChainCodeEventsByByte(data)
		if err != nil {
			continue
		}
		events = append(events, event)
	}
	return events
}

func GetChainCodeEventsByByte(data []byte) (*pb.ChaincodeEvent, error) {
	// env := &common.Envelope{}
	// if err := proto.Unmarshal(data, env); err != nil {
	// 	return nil, fmt.Errorf("error reconstructing envelope(%s)", err)
	// }

	env, err := futils.GetEnvelopeFromBlock(data)
	if err != nil {
		return nil, fmt.Errorf("error reconstructing envelope(%s)", err)
	}
	// get the payload from the envelope
	payload, err := futils.GetPayload(env)
	if err != nil {
		return nil, fmt.Errorf("Could not extract payload from envelope, err %s", err)
	}

	chdr, err := futils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, fmt.Errorf("Could not extract channel header from envelope, err %s", err)
	}

	if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {

		tx, err := futils.GetTransaction(payload.Data)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling transaction payload for block event: %s", err)
		}
		//此处应该遍历tx.Actions？
		chaincodeActionPayload, err := futils.GetChaincodeActionPayload(tx.Actions[0].Payload)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling transaction action payload for block event: %s", err)
		}
		propRespPayload, err := futils.GetProposalResponsePayload(chaincodeActionPayload.Action.ProposalResponsePayload)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling proposal response payload for block event: %s", err)
		}

		caPayload, err := futils.GetChaincodeAction(propRespPayload.Extension)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling chaincode action for block event: %s", err)
		}
		ccEvent, err := futils.GetChaincodeEvents(caPayload.Events)
		if ccEvent != nil {
			return ccEvent, nil
		}

	}
	return nil, errors.New("no HeaderType_ENDORSER_TRANSACTION type ")
}

func EventConvert(event *pb.ChaincodeEvent) *ChainTxEvents {
	if event == nil {
		return nil
	}
	clientEvent := &ChainTxEvents{}
	clientEvent.Chaincode = event.ChaincodeId
	clientEvent.Name = event.EventName
	clientEvent.Payload = event.Payload
	clientEvent.TxID = event.TxId
	return clientEvent
}
