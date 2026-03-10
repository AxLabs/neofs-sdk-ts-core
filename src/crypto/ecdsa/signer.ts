import { sha256, sha512 } from '../../utils/crypto';
import { randomBytes } from '../../utils/buffer';
import { ec as EC } from 'elliptic';
import { Scheme, PublicKey } from '../scheme';
import { Signer } from '../signer';
import { ECDSAPublicKey, ECDSAPublicKeyRFC6979 } from './public-key';

const ec = new EC('p256');

/**
 * ECDSA Signer with SHA-512 hashing.
 * Implements neofscrypto.Signer interface.
 */
export class ECDSASigner implements Signer {
  private privateKey: EC.KeyPair;

  constructor(privateKey: EC.KeyPair) {
    this.privateKey = privateKey;
  }

  /**
   * Creates a new ECDSA signer from a private key.
   */
  static fromPrivateKey(privateKey: EC.KeyPair): ECDSASigner {
    return new ECDSASigner(privateKey);
  }

  /**
   * Creates a new ECDSA signer from private key bytes.
   */
  static fromPrivateKeyBytes(privateKeyBytes: Uint8Array): ECDSASigner {
    const privateKey = ec.keyFromPrivate(privateKeyBytes);
    return new ECDSASigner(privateKey);
  }

  /**
   * Creates a new ECDSA signer from WIF (Wallet Import Format) string.
   * Based on C# implementation: LoadWif
   */
  static fromWIF(wif: string): ECDSASigner {
    const privateKeyBytes = this.decodeWIF(wif);
    return this.fromPrivateKeyBytes(privateKeyBytes);
  }

  /**
   * Decodes WIF (Wallet Import Format) string to private key bytes.
   * Based on C# implementation: GetPrivateKeyFromWIF
   */
  static decodeWIF(wif: string): Uint8Array {
    if (!wif) throw new Error('WIF cannot be null');
    
    // Decode Base58
    const data = this.base58Decode(wif);
    
    // Neo WIF format: 0x80 + 32 bytes private key + checksum
    if (data.length < 33 || data[0] !== 0x80) {
      throw new Error(`Invalid WIF format: length=${data.length}, first=${data[0]}`);
    }
    
    // Extract private key (32 bytes) - take the first 32 bytes after the 0x80 prefix
    const privateKey = new Uint8Array(32);
    privateKey.set(data.slice(1, 33));
    
    return privateKey;
  }

  /**
   * Base58 decoding implementation.
   */
  private static base58Decode(encoded: string): Uint8Array {
    const alphabet = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz';
    let value = 0n;
    
    // Decode Base58 string to BigInteger
    for (let i = 0; i < encoded.length; i++) {
      const digit = alphabet.indexOf(encoded[i]);
      if (digit < 0) {
        throw new Error(`Invalid Base58 character '${encoded[i]}' at position ${i}`);
      }
      value = value * 58n + BigInt(digit);
    }
    
    // Convert BigInteger to byte array
    const leadingZeros = encoded.match(/^1+/)?.[0].length || 0;
    const bytes = [];
    
    while (value > 0n) {
      bytes.unshift(Number(value & 0xFFn));
      value = value >> 8n;
    }
    
    // Add leading zeros
    const result = new Uint8Array(leadingZeros + bytes.length);
    result.set(bytes, leadingZeros);
    
    return result;
  }

  /**
   * Generates a new random ECDSA key pair.
   */
  static generate(): ECDSASigner {
    const privateKey = ec.genKeyPair();
    return new ECDSASigner(privateKey);
  }

  scheme(): Scheme {
    return Scheme.ECDSA_SHA512;
  }

  sign(data: Uint8Array): Uint8Array {
    // Hash the data with SHA-512
    const hash = sha512(data);
    
    // Sign the hash
    const signature = this.privateKey.sign(hash);
    
    // Format signature as per NeoFS specification
    const r = signature.r.toArray('be', 32);
    const s = signature.s.toArray('be', 32);
    
    // Create signature buffer: [0x04][r][s]
    const sigBuffer = new Uint8Array(1 + 32 + 32);
    sigBuffer[0] = 0x04; // uncompressed form marker
    sigBuffer.set(r, 1);
    sigBuffer.set(s, 1 + 32);
    
    return sigBuffer;
  }

  public(): PublicKey {
    return new ECDSAPublicKey(this.privateKey);
  }

  getPublicKey(): Uint8Array {
    return new Uint8Array(this.privateKey.getPublic().encode('array', false).slice(1)); // Remove the 0x04 prefix
  }
}

/**
 * ECDSA Signer with deterministic SHA-256 hashing (RFC 6979).
 * Implements neofscrypto.Signer interface.
 */
export class ECDSASignerRFC6979 implements Signer {
  private privateKey: EC.KeyPair;

  constructor(privateKey: EC.KeyPair) {
    this.privateKey = privateKey;
  }

  /**
   * Creates a new ECDSA RFC6979 signer from a private key.
   */
  static fromPrivateKey(privateKey: EC.KeyPair): ECDSASignerRFC6979 {
    return new ECDSASignerRFC6979(privateKey);
  }

  /**
   * Creates a new ECDSA RFC6979 signer from private key bytes.
   */
  static fromPrivateKeyBytes(privateKeyBytes: Uint8Array): ECDSASignerRFC6979 {
    const privateKey = ec.keyFromPrivate(privateKeyBytes);
    return new ECDSASignerRFC6979(privateKey);
  }

  /**
   * Creates a new ECDSA RFC6979 signer from a WIF (Wallet Import Format) string.
   */
  static fromWIF(wif: string): ECDSASignerRFC6979 {
    const privateKeyBytes = ECDSASigner.decodeWIF(wif);
    return ECDSASignerRFC6979.fromPrivateKeyBytes(privateKeyBytes);
  }

  /**
   * Generates a new random ECDSA key pair.
   */
  static generate(): ECDSASignerRFC6979 {
    const privateKey = ec.genKeyPair();
    return new ECDSASignerRFC6979(privateKey);
  }

  scheme(): Scheme {
    return Scheme.ECDSA_DETERMINISTIC_SHA256;
  }

  sign(data: Uint8Array): Uint8Array {
    // Hash the data with SHA-256
    const hash = sha256(data);
    
    // Sign the hash using deterministic ECDSA (RFC 6979)
    const signature = this.privateKey.sign(hash, { canonical: true });
    
    // Format signature as per NeoFS specification (64 bytes: r + s)
    const r = signature.r.toArray('be', 32);
    const s = signature.s.toArray('be', 32);
    
    // Create signature buffer: [r][s] (64 bytes total, no prefix)
    const sigBuffer = new Uint8Array(32 + 32);
    sigBuffer.set(r, 0);
    sigBuffer.set(s, 32);
    
    return sigBuffer;
  }

  public(): PublicKey {
    return new ECDSAPublicKeyRFC6979(this.privateKey);
  }

  getPublicKey(): Uint8Array {
    return new Uint8Array(this.privateKey.getPublic().encode('array', false).slice(1)); // Remove the 0x04 prefix
  }
}
