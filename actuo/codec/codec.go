// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package codec

import (
	"bytes"
	"encoding/json"
	"io"
)

type Encoder[Object any] interface {
	Encode(w io.Writer, obj Object) error
}

type Decoder[Object any] interface {
	// Decode a value into a supplied object, mutating it in the process.
	Decode(r io.Reader, into Object) error
}

type Codec[Object any] interface {
	Encoder[Object]
	Decoder[Object]
}

type jsonCodec[Object any] struct{}

func (jsonCodec[Object]) Encode(w io.Writer, obj Object) error {
	return json.NewEncoder(w).Encode(obj)
}

func (jsonCodec[Object]) Decode(r io.Reader, into Object) error {
	if err := json.NewDecoder(r).Decode(into); err != nil {
		return err
	}
	return nil
}

func JSON[Object any]() Codec[Object] {
	return jsonCodec[Object]{}
}

func Encode[Object any](encoder Encoder[Object], obj Object) ([]byte, error) {
	var buf bytes.Buffer
	if err := encoder.Encode(&buf, obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode[Object any](decoder Decoder[Object], data []byte, into Object) error {
	return decoder.Decode(bytes.NewReader(data), into)
}
