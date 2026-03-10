import { Scheme, PublicKey, getPublicKeyConstructor } from './scheme';
import { Signer, Signature } from './signer';

/**
 * Signature represents a confirmation of data integrity received by the
 * digital signature mechanism.
 * 
 * Signature is mutually compatible with refs.Signature message.
 */
export class NeoFSSignature implements Signature {
  private _scheme: Scheme;
  private _publicKey: Uint8Array;
  private _value: Uint8Array;

  constructor(scheme: Scheme, publicKey: Uint8Array, value: Uint8Array) {
    this._scheme = scheme;
    this._publicKey = new Uint8Array(publicKey);
    this._value = new Uint8Array(value);
  }

  /**
   * Creates a new Signature instance from raw key data.
   */
  static fromRawKey(scheme: Scheme, publicKey: Uint8Array, value: Uint8Array): NeoFSSignature {
    return new NeoFSSignature(scheme, publicKey, value);
  }

  /**
   * Creates a new Signature instance from PublicKey and value.
   */
  static fromPublicKey(scheme: Scheme, publicKey: PublicKey, value: Uint8Array): NeoFSSignature {
    const encodedKey = new Uint8Array(publicKey.maxEncodedSize());
    const written = publicKey.encode(encodedKey);
    if (written < 0) {
      throw new Error('failed to encode public key');
    }
    return new NeoFSSignature(scheme, encodedKey.slice(0, written), value);
  }

  /**
   * Creates signature from given N3 witness scripts.
   */
  static fromN3Witness(invocScript: Uint8Array, verifScript: Uint8Array): NeoFSSignature {
    return new NeoFSSignature(Scheme.N3, verifScript, invocScript);
  }

  /**
   * Creates signature from protobuf message.
   */
  static fromProtoMessage(proto: any): NeoFSSignature {
    if (proto.scheme < 0) {
      throw new Error(`negative scheme ${proto.scheme}`);
    }
    return new NeoFSSignature(
      proto.scheme as Scheme,
      new Uint8Array(proto.key || []),
      new Uint8Array(proto.sign || [])
    );
  }

  /**
   * Converts to protobuf message.
   */
  toProtoMessage(): any {
    return {
      key: Array.from(this._publicKey),
      sign: Array.from(this._value),
      scheme: this._scheme,
    };
  }

  /**
   * Signs data using Signer and encodes public key for subsequent verification.
   */
  calculate(signer: Signer, data: Uint8Array): void {
    try {
      const signature = signer.sign(data);
      const publicKey = signer.public();
      const encodedKey = new Uint8Array(publicKey.maxEncodedSize());
      const written = publicKey.encode(encodedKey);
      
      if (written < 0) {
        throw new Error('failed to encode public key');
      }

      this._scheme = signer.scheme();
      this._publicKey = encodedKey.slice(0, written);
      this._value = signature;
    } catch (error) {
      throw new Error(`signer ${signer.constructor.name} failure: ${error}`);
    }
  }

  /**
   * Verifies data signature using encoded public key.
   */
  verify(data: Uint8Array): boolean {
    try {
      const key = this.decodePublicKey();
      return key.verify(data, this._value);
    } catch {
      return false;
    }
  }

  /**
   * Returns signature scheme used by signer to calculate the signature.
   */
  scheme(): Scheme {
    return this._scheme;
  }

  /**
   * Sets signature scheme used by signer to calculate the signature.
   */
  setScheme(scheme: Scheme): void {
    this._scheme = scheme;
  }

  /**
   * Returns public key of the signer which calculated the signature.
   */
  publicKey(): PublicKey {
    return this.decodePublicKey();
  }

  /**
   * Sets binary-encoded public key of the signer which calculated the signature.
   */
  setPublicKeyBytes(publicKey: Uint8Array): void {
    this._publicKey = new Uint8Array(publicKey);
  }

  /**
   * Returns binary-encoded public key of the signer which calculated the signature.
   */
  publicKeyBytes(): Uint8Array {
    return new Uint8Array(this._publicKey);
  }

  /**
   * Sets calculated digital signature.
   */
  setValue(value: Uint8Array): void {
    this._value = new Uint8Array(value);
  }

  /**
   * Returns calculated digital signature.
   */
  value(): Uint8Array {
    return new Uint8Array(this._value);
  }

  /**
   * Decodes the public key from the stored bytes.
   */
  private decodePublicKey(): PublicKey {
    const constructor = getPublicKeyConstructor(this._scheme);
    if (!constructor) {
      throw new Error(`unsupported scheme ${this._scheme}`);
    }

    const publicKey = constructor();
    publicKey.decode(this._publicKey);
    return publicKey;
  }
}

/**
 * Helper function to get public key bytes from a PublicKey.
 */
export function publicKeyBytes(publicKey: PublicKey): Uint8Array {
  const buf = new Uint8Array(publicKey.maxEncodedSize());
  const written = publicKey.encode(buf);
  if (written < 0) {
    throw new Error('failed to encode public key');
  }
  return buf.slice(0, written);
}
