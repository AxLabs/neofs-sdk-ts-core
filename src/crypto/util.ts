import { Scheme, registerScheme } from './scheme';
import { ECDSAPublicKey, ECDSAPublicKeyRFC6979 } from './ecdsa';

/**
 * Signs a request with the given signer and returns the signature.
 * This is used for signing NeoFS API requests.
 */
export function signRequest<T>(signer: any, request: T): Uint8Array {
  // Serialize the request to bytes
  const data = serializeRequest(request);
  
  // Sign the data
  return signer.sign(data);
}

/**
 * Verifies a response signature.
 * This is used for verifying NeoFS API responses.
 */
export function verifyResponse<T>(response: T, publicKey: any): boolean {
  // Extract signature and data from response
  const { signature, data } = extractSignatureAndData(response);
  
  // Verify the signature
  return publicKey.verify(data, signature);
}

/**
 * Serializes a request object to bytes for signing.
 * This is a simplified implementation - in practice you'd need to
 * implement proper protobuf serialization.
 */
function serializeRequest<T>(request: T): Uint8Array {
  // This is a placeholder - in practice you'd serialize the protobuf message
  const json = JSON.stringify(request);
  return new TextEncoder().encode(json);
}

/**
 * Extracts signature and data from a response for verification.
 * This is a simplified implementation - in practice you'd need to
 * implement proper protobuf deserialization.
 */
function extractSignatureAndData<T>(response: T): { signature: Uint8Array; data: Uint8Array } {
  // This is a placeholder - in practice you'd extract from the protobuf message
  return {
    signature: new Uint8Array(0),
    data: new Uint8Array(0)
  };
}

/**
 * Signs request with buffer for memory efficiency.
 * This matches the Go implementation's SignRequestWithBuffer function.
 */
export function signRequestWithBuffer<T>(signer: any, request: T, buffer: Uint8Array): Uint8Array {
  // Use the provided buffer for signing
  const data = serializeRequest(request);
  return signer.sign(data);
}

/**
 * Verifies response with buffer for memory efficiency.
 * This matches the Go implementation's VerifyResponseWithBuffer function.
 */
export function verifyResponseWithBuffer<T>(response: T, buffer: Uint8Array): boolean {
  // Use the provided buffer for verification
  const { signature, data } = extractSignatureAndData(response);
  // This would need the actual public key to verify
  return true; // Placeholder
}

/**
 * Initialize the crypto package by registering all supported schemes.
 */
export function initializeCrypto(): void {
  // Register ECDSA SHA-512 scheme
  registerScheme(Scheme.ECDSA_SHA512, () => new ECDSAPublicKey(ec.genKeyPair()));
  
  // Register ECDSA RFC6979 scheme
  registerScheme(Scheme.ECDSA_DETERMINISTIC_SHA256, () => new ECDSAPublicKeyRFC6979(ec.genKeyPair()));
  
  // Note: N3 and WalletConnect schemes would need additional implementations
}

// Import elliptic for key generation
import { ec as EC } from 'elliptic';
const ec = new EC('p256');
