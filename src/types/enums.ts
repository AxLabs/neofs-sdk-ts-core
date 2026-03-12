/**
 * NeoFS protocol enums derived from the protobuf definitions.
 * These are re-exported by platform-specific SDKs (node, react-native).
 */

/**
 * Search filter match types for object attribute queries.
 *
 * Values correspond to the NeoFS protobuf `MatchType` enum.
 * IMPORTANT: `UNSPECIFIED` (0) is the protobuf default and is NOT serialized
 * on the wire — using it in a search filter effectively disables the filter.
 * Always use an explicit match type like `STRING_EQUAL`.
 */
export enum MatchType {
  UNSPECIFIED = 0,
  STRING_EQUAL = 1,
  STRING_NOT_EQUAL = 2,
  NOT_PRESENT = 3,
  COMMON_PREFIX = 4,
  NUM_GT = 5,
  NUM_GE = 6,
  NUM_LT = 7,
  NUM_LE = 8,
}

/**
 * Checksum algorithm identifiers used in object headers.
 *
 * Values correspond to the NeoFS protobuf `ChecksumType` enum.
 */
export enum ChecksumType {
  CHECKSUM_TYPE_UNSPECIFIED = 0,
  /** Tillich-Zémor homomorphic hash */
  TZ = 1,
  /** SHA-256 */
  SHA256 = 2,
}

/**
 * Object types in NeoFS.
 *
 * Values correspond to the NeoFS protobuf `ObjectType` enum.
 */
export enum ObjectType {
  REGULAR = 0,
  TOMBSTONE = 1,
  STORAGE_GROUP = 2,
  LOCK = 3,
  LINK = 4,
}
