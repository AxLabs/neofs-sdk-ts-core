# neofs-sdk-ts-core

Shared core package for the [NeoFS](https://fs.neo.org/) TypeScript SDK. Provides cryptographic primitives, user identity utilities, type definitions, and helper functions used by platform-specific SDK packages.

- **neofs-sdk-ts-react-native** — uses `grpc-react-native` with custom protobuf serialization
- **neofs-sdk-ts-node** — uses `@grpc/grpc-js` with standard protobuf

## Installation

```bash
npm install neofs-sdk-ts-core
```

## Modules

The package exposes four modules, each available as a subpath export:

### Crypto (`neofs-sdk-ts-core/crypto`)

ECDSA signing and verification with multiple signature schemes supported by NeoFS:

| Scheme | Description |
|--------|-------------|
| `ECDSA_SHA512` | ECDSA with SHA-512 hashing (FIPS 186-3) |
| `ECDSA_DETERMINISTIC_SHA256` | Deterministic ECDSA with SHA-256 (RFC 6979) |
| `ECDSA_WALLETCONNECT` | WalletConnect signature scheme |
| `N3` | Neo N3 witness |

Key exports:

- `Signer` / `SignerV2` — interfaces for signing operations
- `PublicKey` — interface for public key encoding, decoding, and verification
- `ECDSASigner` / `ECDSASignerRFC6979` — concrete ECDSA signer implementations
- `NeoFSSignature` — signature wrapper with scheme-aware verification
- `Scheme` — enum of supported signature schemes
- `registerScheme()` — register custom public key constructors per scheme
- Tillich–Zémor homomorphic hash (`tz` submodule)

```typescript
import { ECDSASigner, Scheme } from 'neofs-sdk-ts-core/crypto';

const signer = new ECDSASigner(privateKeyHex);
const signature = signer.sign(data);
const pubKey = signer.getPublicKey();
```

### User (`neofs-sdk-ts-core/user`)

NeoFS Owner ID (User ID) derivation and validation. An owner ID is a 25-byte value derived from a Neo N3 public key (address version prefix + script hash + checksum).

```typescript
import { ownerIdFromPublicKey, validateOwnerId } from 'neofs-sdk-ts-core/user';

const ownerId = ownerIdFromPublicKey(compressedPublicKey);
const isValid = validateOwnerId(ownerId);
```

### Types (`neofs-sdk-ts-core/types`)

- `Decimal` — arbitrary-precision decimal arithmetic for monetary computations (avoids floating-point issues). Supports protobuf serialization, add/subtract/multiply/divide, and comparison.

```typescript
import { Decimal } from 'neofs-sdk-ts-core/types';

const balance = new Decimal(1000000n, 12);
console.log(balance.toString()); // "0.000001000000"
```

### Utils (`neofs-sdk-ts-core/utils`)

Buffer and cryptographic helper functions:

- `hexToBytes()` / `bytesToHex()` — hex string conversion
- `fromBase64()` / `toBase64()` — Base64 encoding/decoding
- `randomBytes()` — secure random byte generation
- `sha256()` / `sha512()` / `ripemd160()` / `doubleSha256()` — hash functions

```typescript
import { hexToBytes, sha256 } from 'neofs-sdk-ts-core/utils';

const data = hexToBytes('deadbeef');
const hash = sha256(data);
```

## Building

```bash
# Install dependencies
npm install

# Build (compiles TypeScript to dist/)
npm run build

# Clean build artifacts
npm run clean
```

## protoc-gen-grpc-ts

This repository also contains `protoc-gen-grpc-ts`, a Go-based protoc plugin that generates TypeScript code from `.proto` files. It produces:

- Binary protobuf serialization (`BinaryWriter` / `BinaryReader`)
- gRPC service clients for React Native (`grpc-react-native`) and Node.js (`@grpc/grpc-js`)
- Support for unary, server-streaming, client-streaming, and bidirectional streaming RPCs

```bash
# Build the plugin
cd protoc-gen-grpc-ts
make build

# Install to PATH
make install

# Run tests
make test
```

See [`protoc-gen-grpc-ts/README.md`](protoc-gen-grpc-ts/README.md) for full usage documentation.

## License

[Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0)
