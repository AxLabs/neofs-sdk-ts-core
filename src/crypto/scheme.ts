/**
 * Digital signature algorithm schemes supported by NeoFS.
 * Negative values are reserved and depend on context (e.g. unsupported scheme).
 */
export enum Scheme {
  ECDSA_SHA512 = 0,               // ECDSA with SHA-512 hashing (FIPS 186-3)
  ECDSA_DETERMINISTIC_SHA256 = 1, // Deterministic ECDSA with SHA-256 hashing (RFC 6979)
  ECDSA_WALLETCONNECT = 2,        // Wallet Connect signature scheme
  N3 = 3,                         // Neo N3 witness
}

/**
 * Error thrown when an incorrect signer is used.
 */
export class IncorrectSignerError extends Error {
  constructor(message: string = 'incorrect signer') {
    super(message);
    this.name = 'IncorrectSignerError';
  }
}

/**
 * Maps scheme to public key constructor.
 */
const publicKeyConstructors = new Map<Scheme, () => PublicKey>();

/**
 * Register a function that returns a new blank PublicKey instance for the given Scheme.
 * This is intended to be called from the init function in packages that implement signature schemes.
 * 
 * @param scheme - The signature scheme
 * @param constructor - Function that creates a new PublicKey instance
 * @throws Error if scheme is already registered
 */
export function registerScheme(scheme: Scheme, constructor: () => PublicKey): void {
  if (publicKeyConstructors.has(scheme)) {
    throw new Error(`scheme ${scheme} is already registered`);
  }
  publicKeyConstructors.set(scheme, constructor);
}

/**
 * Get the public key constructor for a given scheme.
 * 
 * @param scheme - The signature scheme
 * @returns The constructor function or undefined if not found
 */
export function getPublicKeyConstructor(scheme: Scheme): (() => PublicKey) | undefined {
  return publicKeyConstructors.get(scheme);
}

/**
 * PublicKey interface for signature verification.
 */
export interface PublicKey {
  /**
   * Returns maximum size required for binary-encoded public key.
   * Must not return value greater than any return of Encode.
   */
  maxEncodedSize(): number;

  /**
   * Encodes public key into buffer. Returns number of bytes written.
   * Must panic if buffer size is insufficient and less than maxEncodedSize.
   * Must return negative value on any failure except buffer size.
   * 
   * @param buf - Buffer to encode into
   * @returns Number of bytes written, or negative value on failure
   */
  encode(buf: Uint8Array): number;

  /**
   * Decodes binary public key.
   * 
   * @param data - Binary data to decode
   * @throws Error if decoding fails
   */
  decode(data: Uint8Array): void;

  /**
   * Verifies signature of the given data.
   * 
   * @param data - Data that was signed
   * @param signature - Signature to verify
   * @returns True if signature is valid
   */
  verify(data: Uint8Array, signature: Uint8Array): boolean;
}
