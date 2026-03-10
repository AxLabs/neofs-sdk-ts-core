import { sha256, sha512 } from '../../utils/crypto';
import { ec as EC } from 'elliptic';
import { PublicKey } from '../scheme';

const ec = new EC('p256');

/**
 * ECDSA PublicKey with SHA-512 hashing.
 * Implements neofscrypto.PublicKey interface.
 */
export class ECDSAPublicKey implements PublicKey {
  private publicKey: EC.KeyPair;

  constructor(publicKey: EC.KeyPair) {
    this.publicKey = publicKey;
  }

  /**
   * Creates a new ECDSA public key from a key pair.
   */
  static fromKeyPair(publicKey: EC.KeyPair): ECDSAPublicKey {
    return new ECDSAPublicKey(publicKey);
  }

  /**
   * Creates a new ECDSA public key from compressed bytes.
   */
  static fromBytes(data: Uint8Array): ECDSAPublicKey {
    const publicKey = ec.keyFromPublic(data, 'array');
    return new ECDSAPublicKey(publicKey);
  }

  maxEncodedSize(): number {
    return 33; // Compressed public key size
  }

  encode(buf: Uint8Array): number {
    if (buf.length < 33) {
      throw new Error(`too short buffer ${buf.length}`);
    }

    const compressed = this.publicKey.getPublic(true, 'array');
    buf.set(compressed, 0);
    return 33;
  }

  decode(data: Uint8Array): void {
    try {
      this.publicKey = ec.keyFromPublic(data, 'array');
    } catch (error) {
      throw new Error(`failed to decode ECDSA public key: ${error}`);
    }
  }

  verify(data: Uint8Array, signature: Uint8Array): boolean {
    try {
      // Hash the data with SHA-512
      const hash = sha512(data);
      
      // Parse signature
      if (signature.length !== 65 || signature[0] !== 0x04) {
        return false;
      }

      const r = signature.slice(1, 33);
      const s = signature.slice(33, 65);

      // Verify signature
      return this.publicKey.verify(hash, { r, s });
    } catch {
      return false;
    }
  }
}

/**
 * ECDSA PublicKey with deterministic SHA-256 hashing (RFC 6979).
 * Implements neofscrypto.PublicKey interface.
 */
export class ECDSAPublicKeyRFC6979 implements PublicKey {
  private publicKey: EC.KeyPair;

  constructor(publicKey: EC.KeyPair) {
    this.publicKey = publicKey;
  }

  /**
   * Creates a new ECDSA RFC6979 public key from a key pair.
   */
  static fromKeyPair(publicKey: EC.KeyPair): ECDSAPublicKeyRFC6979 {
    return new ECDSAPublicKeyRFC6979(publicKey);
  }

  /**
   * Creates a new ECDSA RFC6979 public key from compressed bytes.
   */
  static fromBytes(data: Uint8Array): ECDSAPublicKeyRFC6979 {
    const publicKey = ec.keyFromPublic(data, 'array');
    return new ECDSAPublicKeyRFC6979(publicKey);
  }

  maxEncodedSize(): number {
    return 33; // Compressed public key size
  }

  encode(buf: Uint8Array): number {
    if (buf.length < 33) {
      throw new Error(`too short buffer ${buf.length}`);
    }

    const compressed = this.publicKey.getPublic(true, 'array');
    buf.set(compressed, 0);
    return 33;
  }

  decode(data: Uint8Array): void {
    try {
      this.publicKey = ec.keyFromPublic(data, 'array');
    } catch (error) {
      throw new Error(`failed to decode ECDSA RFC6979 public key: ${error}`);
    }
  }

  verify(data: Uint8Array, signature: Uint8Array): boolean {
    try {
      // Hash the data with SHA-256
      const hash = sha256(data);
      
      // Parse the signature format: [r][s] (64 bytes)
      if (signature.length !== 64) {
        return false;
      }
      
      const r = signature.slice(0, 32);
      const s = signature.slice(32, 64);
      
      // Create signature object for elliptic library
      // Convert to hex without Buffer
      const toHex = (bytes: Uint8Array) => Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join('');
      const sigObj = {
        r: toHex(r),
        s: toHex(s)
      };
      
      return this.publicKey.verify(hash, sigObj);
    } catch {
      return false;
    }
  }
}
