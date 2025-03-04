package schemahandler

import (
	"io"

	"github.com/logward/omniparser/customfuncs"
	"github.com/logward/omniparser/errs"
	"github.com/logward/omniparser/header"
	"github.com/logward/omniparser/transformctx"
)

// CreateCtx is a context object for CreateFunc.
type CreateCtx struct {
	Name         string
	Header       header.Header
	Content      []byte
	CustomFuncs  customfuncs.CustomFuncs
	CreateParams interface{}
}

// CreateFunc is a function that checks if a given schema is supported by its associated
// schema handler or not. And, if yes, it parses the schema content, creates and initializes
// a new instance of its associated schema handler.
// If a given schema is not supported, errs.ErrSchemaNotSupported should be returned.
// Any other error returned will cause omniparser to fail entirely.
// Note, any non errs.ErrSchemaNotSupported error returned here is errs.CtxAwareErr
// formatted (i.e. error contains schema name and if possible error line number).
type CreateFunc func(ctx *CreateCtx) (SchemaHandler, error)

// SchemaHandler is an interface representing a schema handler responsible for ingesting,
// processing and transforming input stream based on its given schema.
type SchemaHandler interface {
	// NewIngester returns an Ingester for an input stream.
	// Omniparser will not call NewIngester unless CreateSchemaHandler has returned supported.
	// Omniparser calls NewIngester when client supplies an input stream and is ready
	// for the parser to ingest/process/transform the input.
	NewIngester(ctx *transformctx.Ctx, input io.Reader) (Ingester, error)
}

// RawRecord represents a raw record ingested from the input.
type RawRecord interface {
	// Raw returns the actual raw record that is version specific to each of the schema handlers.
	Raw() interface{}
	// Checksum returns a UUIDv3 (MD5) stable hash of the raw record.
	Checksum() string
}

// Ingester is an interface of ingestion and transformation for a given input stream.
type Ingester interface {
	// Read is called repeatedly during the processing of an input stream. Each call it should return
	// the raw record (type of `interface{}`) and its transformed record (type of `[]byte`). It's
	// entirely up to the implementation of this interface/method to decide whether internally it does
	// all the processing all at once (as processes the entire input the very first call of Read) and
	// only hands out one record per Read call, OR, processes and returns one record for each call.
	// However, the overall design principle of omniparser is to have streaming processing capability
	// so memory won't be a constraint when dealing with large input file. All built-in ingesters are
	// implemented this way.
	Read() (RawRecord, []byte, error)

	// IsContinuableError is called to determine if an error returned by Read is fatal or not. After
	// Read is called, a result record or an error will be returned to caller. After caller consumes
	// the record or the error, omniparser needs to decide whether to continue the transform operation
	// or not, based on whether the last err, if present, is "continuable" or not.
	IsContinuableError(error) bool

	// CtxAwareErr interface is embedded to provide omniparser and custom functions a way to provide
	// context aware (such as input file name + line number) error formatting.
	errs.CtxAwareErr
}
