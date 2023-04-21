package types

// RequestID is used to uniquely identify an outgoing request.
type RequestID = uint64

// BatchID is a unique identifier of a batch.
// TODO: do these make sense here (see batchdbpb.proto)
type BatchID = []byte

type BatchIDString string
