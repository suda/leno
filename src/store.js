import { writable } from 'svelte/store';

const messageStore = writable('');

const socket = new WebSocket('ws://localhost:3000');

// Listen for messages
socket.addEventListener('message', function (event) {
	try {
		const message = JSON.parse(event.data);
		messageStore.set(message);
	} catch (error) {
		console.error(`"${event.data}" is not a valid JSON`, error);
	}
});

export default {
	subscribe: messageStore.subscribe
}

