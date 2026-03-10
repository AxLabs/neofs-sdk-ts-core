// Import Closure Library polyfill first - this must be loaded before any protobuf files
import './utils/closure-polyfill';

export * from './crypto';
export * from './user';
export * from './types';
export * from './utils';

// Note: Platform-specific modules (eacl, bearer, waiter, client) are NOT in core.
// They depend on generated protobuf types which are different for each platform:
// - neofs-sdk-ts-react-native: uses grpc-react-native with custom protobuf serialization
// - neofs-sdk-ts-node: uses @grpc/grpc-js with standard protobuf
