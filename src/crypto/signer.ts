import { Scheme, PublicKey } from './scheme';

/**
 * Signer interface for entities that can be used for signing operations in NeoFS.
 * Unites secret and public parts. For example, an ECDSA private key or external auth service.
 */
export interface Signer {
  /**
   * Returns corresponding signature scheme.
   */
  scheme(): Scheme;

  /**
   * Signs digest of the given data. Implementations encapsulate data
   * hashing that depends on Scheme. For example, if scheme uses SHA-256, then
   * Sign signs SHA-256 hash of the data.
   * 
   * @param data - Data to sign
   * @returns Signature bytes
   * @throws Error if signing fails
   */
  sign(data: Uint8Array): Uint8Array;

  /**
   * Returns the public key associated with this signer.
   * 
   * @returns Public key bytes
   */
  getPublicKey(): Uint8Array;

  /**
   * Returns the public key corresponding to the Signer.
   */
  public(): PublicKey;
}

/**
 * SignerV2 allows to sign authorized data.
 */
export interface SignerV2 {
  /**
   * Signs data and returns a Signature object.
   * 
   * @param data - Data to sign
   * @returns Signature object
   * @throws Error if signing fails
   */
  signData(data: Uint8Array): Signature;
}

/**
 * Signature interface for representing digital signatures.
 */
export interface Signature {
  /**
   * Returns the signature scheme used.
   */
  scheme(): Scheme;

  /**
   * Returns the public key bytes.
   */
  publicKeyBytes(): Uint8Array;

  /**
   * Returns the signature value.
   */
  value(): Uint8Array;

  /**
   * Verifies the signature against the given data.
   * 
   * @param data - Data that was signed
   * @returns True if signature is valid
   */
  verify(data: Uint8Array): boolean;
}

/**
 * OverlapSigner returns a Signer which implements SignerV2 for functions
 * trying to switch to SignerV2 usage. It is unsafe to use OverlapSigner for
 * any other case.
 * 
 * @param signerV2 - SignerV2 instance
 * @returns Signer that wraps SignerV2
 */
export function overlapSigner(signerV2: SignerV2): Signer {
  return {
    scheme: () => signerV2.signData(new Uint8Array(0)).scheme(),
    sign: (data: Uint8Array) => signerV2.signData(data).value(),
    getPublicKey: () => {
      // This is a simplified implementation - in practice you'd need to extract
      // the public key from the SignerV2 implementation
      throw new Error('public key extraction not implemented for SignerV2');
    },
    public: () => {
      // This is a simplified implementation - in practice you'd need to extract
      // the public key from the SignerV2 implementation
      throw new Error('public key extraction not implemented for SignerV2');
    }
  };
}
