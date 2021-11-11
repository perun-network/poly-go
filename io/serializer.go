// Copyright 2019 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package io

import (
	"io"
)

type (
	// Serializer objects can be serialized into and from streams.
	Serializer interface {
		Encoder
		Decoder
	}

	// An Encoder can encode itself into a stream.
	Encoder interface {
		// Encode writes itself to a stream.
		// If the stream fails, the underlying error is returned.
		Encode(io.Writer) error
	}

	// A Decoder can decode itself from a stream.
	Decoder interface {
		// Decode reads an object from a stream.
		// If the stream fails, the underlying error is returned.
		Decode(io.Reader) error
	}
)
