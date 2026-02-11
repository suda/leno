import { writable } from 'svelte/store';

export type LogMessage = Record<string, unknown>;

const messageStore = writable<LogMessage | null>(null);

const source = new EventSource('/events');

source.addEventListener('message', (event: MessageEvent) => {
  try {
    const message = JSON.parse(event.data as string) as LogMessage;
    messageStore.set(message);
  } catch (error) {
    console.error(`"${event.data}" is not a valid JSON`, error);
  }
});

export default {
  subscribe: messageStore.subscribe,
};
