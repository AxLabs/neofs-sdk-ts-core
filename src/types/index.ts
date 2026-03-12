export { Decimal } from './decimal';
export * from './enums';

export interface ContainerID {
  value: Uint8Array;
}

export interface ObjectID {
  value: Uint8Array;
}

export interface Address {
  containerId: ContainerID;
  objectId: ObjectID;
}

/** Object header attribute as returned by the NeoFS API. */
export interface ObjectAttribute {
  key: string;
  value: string;
}

/** Object header as returned by head / get operations. */
export interface ObjectHeader {
  containerId?: ContainerID;
  ownerId?: Uint8Array;
  attributes?: ObjectAttribute[];
  payloadLength?: number;
  payloadHash?: { type: number; sum: Uint8Array };
  homomorphicHash?: { type: number; sum: Uint8Array };
  objectType?: number;
  version?: { major: number; minor: number };
  creationEpoch?: number;
}
