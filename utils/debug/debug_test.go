package debug

import (
	"errors"
	"strings"
	"testing"

	"tgragnato.it/goflow/producer"
)

func TestPanicErrorMessageError(t *testing.T) {
	t.Parallel()
	e := &PanicErrorMessage{Inner: "boom"}
	if e.Error() != "boom" {
		t.Fatalf("unexpected Error(): %q", e.Error())
	}
}

func TestPanicErrorMessageUnwrap(t *testing.T) {
	t.Parallel()
	e := &PanicErrorMessage{}
	unwrapped := e.Unwrap()
	if len(unwrapped) != 1 || !errors.Is(unwrapped[0], ErrPanic) {
		t.Fatal("expected ErrPanic in Unwrap")
	}
}

func TestPanicDecoderWrapperNoPanic(t *testing.T) {
	t.Parallel()
	fn := PanicDecoderWrapper(func(msg any) error {
		return nil
	})
	if err := fn("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPanicDecoderWrapperWrappedError(t *testing.T) {
	t.Parallel()
	fn := PanicDecoderWrapper(func(msg any) error {
		return errors.New("decode failed")
	})
	err := fn("hello")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "decoder panic wrapper") {
		t.Fatalf("expected decoder panic wrapper prefix, got %v", err)
	}
}

func TestPanicDecoderWrapperPanic(t *testing.T) {
	t.Parallel()
	fn := PanicDecoderWrapper(func(msg any) error {
		panic("test panic")
	})
	err := fn("hello")
	if err == nil {
		t.Fatal("expected error from panic")
	}
	if err.Error() != "test panic" {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestWrapPanicProducer(t *testing.T) {
	t.Parallel()
	wrapped := &stubProducer{}
	p := WrapPanicProducer(wrapped)
	if p == nil {
		t.Fatal("expected non-nil producer")
	}
}

func TestPanicProducerProduceNoPanic(t *testing.T) {
	t.Parallel()
	p := WrapPanicProducer(&stubProducer{produceErr: nil})
	msgs, err := p.Produce("msg", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
}

func TestPanicProducerProduceWrappedError(t *testing.T) {
	t.Parallel()
	p := WrapPanicProducer(&stubProducer{produceErr: errors.New("produce failed")})
	_, err := p.Produce("msg", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "producer panic wrapper") {
		t.Fatalf("expected producer panic wrapper prefix, got %v", err)
	}
}

func TestPanicProducerProducePanic(t *testing.T) {
	t.Parallel()
	p := WrapPanicProducer(&panicProducer{})
	_, err := p.Produce("msg", nil)
	if err == nil {
		t.Fatal("expected error from panic")
	}
	if err.Error() != "producer panic" {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestPanicProducerCommit(t *testing.T) {
	t.Parallel()
	p := WrapPanicProducer(&stubProducer{})
	p.Commit(nil)
}

func TestPanicProducerClose(t *testing.T) {
	t.Parallel()
	p := WrapPanicProducer(&stubProducer{})
	p.Close()
}

type stubProducer struct {
	produceErr error
}

func (s *stubProducer) Produce(msg any, args *producer.ProduceArgs) ([]producer.ProducerMessage, error) {
	return []producer.ProducerMessage{nil}, s.produceErr
}

func (s *stubProducer) Commit(flowMessageSet []producer.ProducerMessage) {}

func (s *stubProducer) Close() {}

type panicProducer struct{}

func (p *panicProducer) Produce(msg any, args *producer.ProduceArgs) ([]producer.ProducerMessage, error) {
	panic("producer panic")
}

func (p *panicProducer) Commit(flowMessageSet []producer.ProducerMessage) {}

func (p *panicProducer) Close() {}

// Ensure stubProducer implements ProducerInterface
var _ producer.ProducerInterface = (*stubProducer)(nil)
var _ producer.ProducerInterface = (*panicProducer)(nil)
