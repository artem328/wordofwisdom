package proto

import (
	"fmt"
	"strconv"
)

type MsgType byte

const (
	MsgTypeQuoteRequest      MsgType = 1
	MsgTypeQuote             MsgType = 2
	MsgTypeError             MsgType = 253
	MsgTypeChallenge         MsgType = 254
	MsgTypeChallengeSolution MsgType = 255
)

func (mt MsgType) String() string {
	switch mt {
	case MsgTypeQuoteRequest:
		return "QuoteRequest"
	case MsgTypeQuote:
		return "Quote"
	case MsgTypeError:
		return "Error"
	case MsgTypeChallenge:
		return "Challenge"
	case MsgTypeChallengeSolution:
		return "ChallengeSolution"
	default:
		return "unknown(" + strconv.Itoa(int(mt)) + ")"
	}
}

func Decode(msgType MsgType, body []byte) (Msg, error) {
	var msg Msg

	switch msgType {
	case MsgTypeQuoteRequest:
		msg = new(QuoteRequest)
	case MsgTypeQuote:
		msg = new(Quote)
	case MsgTypeError:
		msg = new(Error)
	case MsgTypeChallenge:
		msg = new(Challenge)
	case MsgTypeChallengeSolution:
		msg = new(ChallengeSolution)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownMsgType, msgType)
	}

	if err := msg.Decode(body); err != nil {
		return nil, fmt.Errorf("decode %s: %w", msgType, err)
	}

	return msg, nil
}

func Encode(msg Msg) (msgType MsgType, body []byte, err error) {
	switch msg.(type) {
	case *QuoteRequest:
		msgType = MsgTypeQuoteRequest
	case *Quote:
		msgType = MsgTypeQuote
	case *Error:
		msgType = MsgTypeError
	case *Challenge:
		msgType = MsgTypeChallenge
	case *ChallengeSolution:
		msgType = MsgTypeChallengeSolution
	default:
		return 0, nil, fmt.Errorf("%w: %T", ErrUnknownMsgType, msg)
	}

	body, err = msg.Encode()
	if err != nil {
		return 0, nil, fmt.Errorf("encode %s: %w", msgType, err)
	}

	return msgType, body, nil
}
