import { NeoFSSignature, publicKeyBytes } from '../signature';
import { ECDSASigner, ECDSASignerRFC6979 } from '../ecdsa';
import { Scheme } from '../scheme';
import { initializeCrypto } from '../util';

// Initialize crypto schemes
initializeCrypto();

describe('NeoFSSignature', () => {
  const testData = new TextEncoder().encode('Hello, world!');

  describe('ECDSA SHA-512', () => {
    test('should create and verify signature', () => {
      const signer = ECDSASigner.generate();
      const signature = new NeoFSSignature(Scheme.ECDSA_SHA512, new Uint8Array(0), new Uint8Array(0));
      
      signature.calculate(signer, testData);
      
      expect(signature.scheme()).toBe(Scheme.ECDSA_SHA512);
      expect(signature.verify(testData)).toBe(true);
      expect(signature.value().length).toBeGreaterThan(0);
      expect(signature.publicKeyBytes().length).toBeGreaterThan(0);
    });

    test('should create signature from raw key', () => {
      const signer = ECDSASigner.generate();
      const publicKey = signer.public();
      const publicKeyBytes = new Uint8Array(publicKey.maxEncodedSize());
      publicKey.encode(publicKeyBytes);
      
      const signature = NeoFSSignature.fromRawKey(
        Scheme.ECDSA_SHA512,
        publicKeyBytes,
        new Uint8Array([1, 2, 3, 4])
      );
      
      expect(signature.scheme()).toBe(Scheme.ECDSA_SHA512);
      expect(signature.publicKeyBytes()).toEqual(publicKeyBytes);
      expect(signature.value()).toEqual(new Uint8Array([1, 2, 3, 4]));
    });

    test('should create signature from public key', () => {
      const signer = ECDSASigner.generate();
      const publicKey = signer.public();
      const signature = NeoFSSignature.fromPublicKey(
        Scheme.ECDSA_SHA512,
        publicKey,
        new Uint8Array([1, 2, 3, 4])
      );
      
      expect(signature.scheme()).toBe(Scheme.ECDSA_SHA512);
      // Verify the public key bytes match instead of object reference
      const originalBytes = new Uint8Array(publicKey.maxEncodedSize());
      publicKey.encode(originalBytes);
      expect(signature.publicKeyBytes()).toEqual(originalBytes);
    });

    test('should create N3 signature', () => {
      const invocScript = new Uint8Array([1, 2, 3]);
      const verifScript = new Uint8Array([4, 5, 6]);
      
      const signature = NeoFSSignature.fromN3Witness(invocScript, verifScript);
      
      expect(signature.scheme()).toBe(Scheme.N3);
      expect(signature.publicKeyBytes()).toEqual(verifScript);
      expect(signature.value()).toEqual(invocScript);
    });
  });

  describe('ECDSA RFC6979', () => {
    test('should create and verify signature', () => {
      const signer = ECDSASignerRFC6979.generate();
      const signature = new NeoFSSignature(Scheme.ECDSA_DETERMINISTIC_SHA256, new Uint8Array(0), new Uint8Array(0));
      
      signature.calculate(signer, testData);
      
      expect(signature.scheme()).toBe(Scheme.ECDSA_DETERMINISTIC_SHA256);
      expect(signature.verify(testData)).toBe(true);
    });
  });

  describe('Proto message conversion', () => {
    test('should convert to and from proto message', () => {
      const signer = ECDSASigner.generate();
      const signature = new NeoFSSignature(Scheme.ECDSA_SHA512, new Uint8Array(0), new Uint8Array(0));
      signature.calculate(signer, testData);
      
      const proto = signature.toProtoMessage();
      expect(proto.scheme).toBe(Scheme.ECDSA_SHA512);
      expect(proto.key).toBeDefined();
      expect(proto.sign).toBeDefined();
      
      const restoredSignature = NeoFSSignature.fromProtoMessage(proto);
      expect(restoredSignature.scheme()).toBe(signature.scheme());
      expect(restoredSignature.publicKeyBytes()).toEqual(signature.publicKeyBytes());
      expect(restoredSignature.value()).toEqual(signature.value());
    });

    test('should throw error for negative scheme', () => {
      expect(() => {
        NeoFSSignature.fromProtoMessage({ scheme: -1, key: [], sign: [] });
      }).toThrow('negative scheme -1');
    });
  });

  describe('Setters and getters', () => {
    test('should set and get scheme', () => {
      const signature = new NeoFSSignature(Scheme.ECDSA_SHA512, new Uint8Array(0), new Uint8Array(0));
      expect(signature.scheme()).toBe(Scheme.ECDSA_SHA512);
      
      signature.setScheme(Scheme.ECDSA_DETERMINISTIC_SHA256);
      expect(signature.scheme()).toBe(Scheme.ECDSA_DETERMINISTIC_SHA256);
    });

    test('should set and get public key bytes', () => {
      const signature = new NeoFSSignature(Scheme.ECDSA_SHA512, new Uint8Array(0), new Uint8Array(0));
      const newKey = new Uint8Array([1, 2, 3, 4]);
      
      signature.setPublicKeyBytes(newKey);
      expect(signature.publicKeyBytes()).toEqual(newKey);
    });

    test('should set and get value', () => {
      const signature = new NeoFSSignature(Scheme.ECDSA_SHA512, new Uint8Array(0), new Uint8Array(0));
      const newValue = new Uint8Array([5, 6, 7, 8]);
      
      signature.setValue(newValue);
      expect(signature.value()).toEqual(newValue);
    });
  });
});

describe('publicKeyBytes helper', () => {
  test('should extract public key bytes', () => {
    const signer = ECDSASigner.generate();
    const publicKey = signer.public();
    
    const bytes = publicKeyBytes(publicKey);
    expect(bytes.length).toBe(33); // Compressed public key size
  });
});
