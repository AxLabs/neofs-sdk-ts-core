export { Decimal } from './decimal';

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
