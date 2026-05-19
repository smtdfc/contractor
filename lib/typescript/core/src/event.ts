export interface EventMetadata {
    name: string;
    method: string;
}

export type EventPayload<T> = T;
