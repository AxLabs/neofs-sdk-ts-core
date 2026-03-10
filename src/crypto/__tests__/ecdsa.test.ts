import { ECDSASigner, ECDSASignerRFC6979 } from '../ecdsa';
import { Scheme } from '../scheme';

describe('ECDSA Signers', () => {
  const testData = new TextEncoder().encode('Hello, world!');

  describe('ECDSASigner', () => {
    test('should generate new signer', () => {
      const signer = ECDSASigner.generate();
      expect(signer.scheme()).toBe(Scheme.ECDSA_SHA512);
      expect(signer.public()).toBeDefined();
    });

    test('should create signer from private key bytes', () => {
      const originalSigner = ECDSASigner.generate();
      const privateKeyBytes = originalSigner['privateKey'].getPrivate('hex');
      const privateKeyArray = new Uint8Array(privateKeyBytes.length / 2);
      for (let i = 0; i < privateKeyBytes.length; i += 2) {
        privateKeyArray[i / 2] = parseInt(privateKeyBytes.substr(i, 2), 16);
      }
      
      const newSigner = ECDSASigner.fromPrivateKeyBytes(privateKeyArray);
      expect(newSigner.scheme()).toBe(Scheme.ECDSA_SHA512);
    });

    test('should sign and verify data', () => {
      const signer = ECDSASigner.generate();
      const signature = signer.sign(testData);
      
      expect(signature.length).toBe(65); // 1 + 32 + 32
      expect(signature[0]).toBe(0x04); // Uncompressed form marker
      
      const publicKey = signer.public();
      expect(publicKey.verify(testData, signature)).toBe(true);
    });

    test('should fail verification with wrong data', () => {
      const signer = ECDSASigner.generate();
      const signature = signer.sign(testData);
      
      const wrongData = new TextEncoder().encode('Wrong data');
      const publicKey = signer.public();
      expect(publicKey.verify(wrongData, signature)).toBe(false);
    });

    test('should fail verification with wrong signature', () => {
      const signer = ECDSASigner.generate();
      const wrongSignature = new Uint8Array(65);
      wrongSignature[0] = 0x04;
      
      const publicKey = signer.public();
      expect(publicKey.verify(testData, wrongSignature)).toBe(false);
    });
  });

  describe('ECDSASignerRFC6979', () => {
    test('should generate new signer', () => {
      const signer = ECDSASignerRFC6979.generate();
      expect(signer.scheme()).toBe(Scheme.ECDSA_DETERMINISTIC_SHA256);
      expect(signer.public()).toBeDefined();
    });

    test('should create signer from private key bytes', () => {
      const originalSigner = ECDSASignerRFC6979.generate();
      const privateKeyBytes = originalSigner['privateKey'].getPrivate('hex');
      const privateKeyArray = new Uint8Array(privateKeyBytes.length / 2);
      for (let i = 0; i < privateKeyBytes.length; i += 2) {
        privateKeyArray[i / 2] = parseInt(privateKeyBytes.substr(i, 2), 16);
      }
      
      const newSigner = ECDSASignerRFC6979.fromPrivateKeyBytes(privateKeyArray);
      expect(newSigner.scheme()).toBe(Scheme.ECDSA_DETERMINISTIC_SHA256);
    });

    test('should sign and verify data', () => {
      const signer = ECDSASignerRFC6979.generate();
      const signature = signer.sign(testData);
      
      expect(signature.length).toBe(64); // 32 + 32 (r + s, no prefix for RFC6979)
      
      const publicKey = signer.public();
      expect(publicKey.verify(testData, signature)).toBe(true);
    });

    test('should produce deterministic signatures', () => {
      const signer = ECDSASignerRFC6979.generate();
      const signature1 = signer.sign(testData);
      const signature2 = signer.sign(testData);
      
      // RFC6979 should produce deterministic signatures
      expect(signature1).toEqual(signature2);
    });
  });

  describe('Public Key Operations', () => {
    test('should encode and decode public key', () => {
      const signer = ECDSASigner.generate();
      const publicKey = signer.public();
      
      const encoded = new Uint8Array(publicKey.maxEncodedSize());
      const written = publicKey.encode(encoded);
      expect(written).toBe(33);
      
      const decoded = new ECDSASigner(ec.genKeyPair()).public();
      decoded.decode(encoded);
      
      // Both keys should verify the same signature
      const signature = signer.sign(testData);
      expect(decoded.verify(testData, signature)).toBe(true);
    });

    test('should handle invalid public key data', () => {
      const signer = ECDSASigner.generate();
      const publicKey = signer.public();
      
      const invalidData = new Uint8Array([1, 2, 3]); // Too short
      expect(() => {
        publicKey.decode(invalidData);
      }).toThrow();
    });
  });
});

// Import elliptic for key generation
import { ec as EC } from 'elliptic';
const ec = new EC('p256');
