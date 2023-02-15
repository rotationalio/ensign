package pem

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"io"
)

// Writer wraps an io.WriteCloser and allows users to write multiple PEM encoded blocks
// to the underlying writer. Unlike pem.Encode the writer has type specific encoders to
// make it easier to write data from different types. This writer is not thread-safe.
type Writer struct {
	out io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w}
}

// Close the underlying writer if it is a closer, otherwise simply returns nil.
func (w *Writer) Close() error {
	if closer, ok := w.out.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Write data to the underlying writer, bypassing any PEM encoding. Write is not safe
// for concurrent use even if the underlying writer is. Generally speaking this function
// is used with pem.Encode(writer, block) though writer.Encode is also an alias.
func (w *Writer) Write(p []byte) (int, error) {
	if w.out == nil {
		w.out = bytes.NewBuffer(nil)
	}
	return w.out.Write(p)
}

// Encode is shorthand for pem.Encode writing a block fo data to the writer.
func (w *Writer) Encode(block *pem.Block) error {
	return pem.Encode(w, block)
}

// EncodePrivateKey as a PKCS8 ASN.1 DER key and write a PEM block with type "PRIVATE KEY"
func (w *Writer) EncodePrivateKey(key interface{}) error {
	return encodePrivateKeyTo(key, w)
}

// EncodePublicKey as a PKIX ASN1.1 DER key and write a PEM block with type "PUBLIC KEY"
func (w *Writer) EncodePublicKey(key interface{}) error {
	return encodePublicKeyTo(key, w)
}

// EncodeCertificate and write a PEM block with type "CERTIFICATE"
func (w *Writer) EncodeCertificate(c *x509.Certificate) error {
	return encodeCertificateTo(c, w)
}

// EncodeCSR and write a PEM block with type "CERTIFICATE REQUEST"
func (w *Writer) EncodeCSR(c *x509.CertificateRequest) error {
	return encodeCSRTo(c, w)
}

// Reader wraps an io.ReadCloser in order to decode bytes from the underlying data. As
// data is decoded from the underlying reader it is freed from memory and cannot be
// read again. A typical use case is to create a reader to loop over all blocks in the
// underlying data as though it were an iterator.
type Reader struct {
	in      io.Reader // reader to read data from
	data    []byte    // remaining bytes to decode
	hasNext bool      // if there is a next block to decode
}

func NewReader(r io.Reader) (_ *Reader, err error) {
	reader := &Reader{r, nil, false}
	if reader.data, err = io.ReadAll(r); err != nil {
		return nil, err
	}

	reader.hasNext = len(reader.data) > 0
	return reader, nil
}

// Close the underlying writer if it is a closer, otherwise simply returns nil.
func (r *Reader) Close() error {
	if closer, ok := r.in.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Read allows the user to read any remaining data that has not been parsed by the PEM
// decoder. This is useful if there is any trailing data in the reader that is not PEM
// encoded, or if the user would like inspect the data without interrupting the read.
func (r *Reader) Read(p []byte) (int, error) {
	buf := bytes.NewReader(r.data)
	return buf.Read(p)
}

// Next is used to determine if there is a next block of data that can be read and is
// typically used for iteration over all of the blocks in a PEM encoded file.
func (r *Reader) Next() bool {
	return r.hasNext
}

// Decode the next block and move the cursor. If block is nil then no block was able to
// be decoded and the rest of the data can be read using the Read function.
func (r *Reader) Decode() *Block {
	var block *pem.Block
	block, r.data = pem.Decode(r.data)
	r.hasNext = block != nil && len(r.data) > 0

	if block != nil {
		return &Block{*block}
	}
	return nil
}

// Block wraps a pem.Block and adds type-specific decoding functions.
type Block struct {
	pem.Block
}

// DecodePrivateKey from a PEM encoded block. If the block type is "EC PRIVATE KEY",
// then the block is parsed as an EC private key in SEC 1, ASN.1 DER form. If the block
// is "RSA PRIVATE KEY" then it is decoded as a PKCS 1, ASN.1 DER form. If the block
// type is "PRIVATE KEY", the block is decoded as a PKCS 8 ASN.1 DER key, if that fails,
// then the PKCS 1 and EC parsers are tried in that order, before returning an error.
func (b *Block) DecodePrivateKey() (interface{}, error) {
	return ParsePrivateKey(&b.Block)
}

// DecodePublicKey from a PEM encoded block. If the block type is "RSA PUBLIC KEY",
// then it is deocded as a PKCS 1, ASN.1 DER form. If the block is "PUBLIC KEY", then it
// is decoded from PKIX ASN1.1 DER form.
func (b *Block) DecodePublicKey() (interface{}, error) {
	return ParsePublicKey(&b.Block)
}

// DecodeCertificate from PEM encoded block with type "CERTIFICATE"
func (b *Block) DecodeCertificate() (*x509.Certificate, error) {
	if b.Type != BlockCertificate {
		return nil, ErrDecodeCertificate
	}
	return x509.ParseCertificate(b.Bytes)
}

// DecodeCSR from PEM encoded block with type "CERTIFICATE REQUEST"
func (b *Block) DecodeCSR() (*x509.CertificateRequest, error) {
	if b.Type != BlockCertificateRequest {
		return nil, ErrDecodeCSR
	}
	return x509.ParseCertificateRequest(b.Bytes)
}
