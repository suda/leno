<script>
	import { onMount } from 'svelte';
	import store from './store.js';
	let messages = [];
	let filteredMessages = [];
	let keys = [];
	let visibleKeys = {};
	let searchTerm = '';

	function addKeys(newKeys) {
		keys.push(...newKeys);
		keys = [...new Set(keys)];

		for (let key of keys) {
			if (!visibleKeys.hasOwnProperty(key)) visibleKeys[key] = true;
		}
	}

	function processMessage(currentMessage) {
		messages = [currentMessage, ...messages];
		addKeys(Object.keys(currentMessage));
		if (filterMessage(currentMessage)) {
			filteredMessages = [currentMessage, ...filteredMessages];
		}
	}

	function filterMessage(currentMessage) {
		let pass = false;
		if (searchTerm !== '') {
			pass = keys.reduce((acc, key) => {
				if (currentMessage.hasOwnProperty(key)) {
					return acc || currentMessage[key].toString().toLowerCase().indexOf(searchTerm) > -1;
				}
				return acc;
			}, pass)
		} else {
			pass = true;
		}
		return pass;
	}

	function applyFilters() {
		filteredMessages = messages.filter(filterMessage);
	}

	onMount(() => {
		store.subscribe(processMessage);
	});
</script>

<main>
	<sidebar>
		<ul class="flex flex-col w-full">
			<li class="my-px">
				<span class="flex flex-row items-center justify-center h-12 px-4 text-gray-600">
					<span class="flex items-center justify-center text-lg text-gray-600">
						<svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16" viewBox="0 0 64 64" x="0px" y="0px" fill="currentColor">
							<g>
								<polygon points="53.414 37 51.414 35 39 35 39 37 50.586 37 52.586 39 58 39 58 37 53.414 37"/>
								<polygon points="32.586 42 24 42 24 44 33.414 44 35.414 42 41 42 41 40 34.586 40 32.586 42"/>
								<rect x="33" y="30" width="4" height="2"/>
								<path d="M11,30c-2.8,0-5,3.075-5,7s2.2,7,5,7,5-3.075,5-7S13.8,30,11,30Zm0,12c-1.626,0-3-2.29-3-5s1.374-5,3-5,3,2.29,3,5S12.626,42,11,42Z"/>
								<path d="M53,26H50.109a6.433,6.433,0,0,1,1.331-3H55a5.006,5.006,0,0,0,5-5V17H55a4.99,4.99,0,0,0-4.956,4.565A8.545,8.545,0,0,0,48.086,26H43.721l-2-6H34.279l-2,6H25.118a6.284,6.284,0,0,0-3.181-5.62A4.989,4.989,0,0,0,17,16H12v1a5.006,5.006,0,0,0,5,5h3.738a4.284,4.284,0,0,1,2.373,4H11c-4.962,0-9,4.935-9,11s4.038,11,9,11H53c4.962,0,9-4.935,9-11S57.962,26,53,26Zm2-7h2.829A3.006,3.006,0,0,1,55,21H52.171A3.006,3.006,0,0,1,55,19ZM17,20a3.006,3.006,0,0,1-2.829-2H17a3.006,3.006,0,0,1,2.829,2ZM4,37c0-4.963,3.14-9,7-9s7,4.037,7,9-3.14,9-7,9S4,41.963,4,37Zm55.261,4H46v2H58.2A6.434,6.434,0,0,1,53,46H16.158a11.664,11.664,0,0,0,3.8-8h8.455l2-2H34V34H29.586l-2,2H19.959a12.79,12.79,0,0,0-.953-4H29V30H17.937a10.07,10.07,0,0,0-1.779-2H33.721l2-6h4.558l1.334,4H38v2H53a6.1,6.1,0,0,1,4.39,2H43v2H58.816A10.833,10.833,0,0,1,60,37,11.048,11.048,0,0,1,59.261,41Z"/>
							</g>
						</svg>
					</span>
					<span class="ml-3 text-4xl font-light">Leño</span>
				</span>
			</li>
			<li class="my-px px-4">
				<div class="pt-2 relative mx-auto text-gray-600">
					<input class="border-2 border-gray-300 bg-white h-10 pl-5 pr-8 w-48 rounded-lg text-sm focus:outline-none"
						type="search" name="search" placeholder="Search" bind:value={searchTerm} on:input={applyFilters}>
					<button type="submit" class="absolute right-0 top-0 mt-5 mr-2 border-none">
						<svg class="text-gray-600 h-4 w-4 fill-current" xmlns="http://www.w3.org/2000/svg"
							xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" id="Capa_1" x="0px" y="0px"
							viewBox="0 0 56.966 56.966" style="enable-background:new 0 0 56.966 56.966;" xml:space="preserve"
							width="512px" height="512px">
							<path
							d="M55.146,51.887L41.588,37.786c3.486-4.144,5.396-9.358,5.396-14.786c0-12.682-10.318-23-23-23s-23,10.318-23,23  s10.318,23,23,23c4.761,0,9.298-1.436,13.177-4.162l13.661,14.208c0.571,0.593,1.339,0.92,2.162,0.92  c0.779,0,1.518-0.297,2.079-0.837C56.255,54.982,56.293,53.08,55.146,51.887z M23.984,6c9.374,0,17,7.626,17,17s-7.626,17-17,17  s-17-7.626-17-17S14.61,6,23.984,6z" />
						</svg>
					</button>
				</div>
			</li>
			<li class="my-px">
				<span class="flex font-medium text-sm text-gray-400 px-4 mt-4 uppercase">Selected fields</span>
			</li>
			{#each keys as key}
				<li class="my-px">
					<span class="flex flex-row items-center h-6 px-5 rounded-lg text-gray-600 hover:bg-gray-100">
						<span class="flex items-center justify-center text-lg text-gray-400">
							<input type="checkbox" bind:checked={visibleKeys[key]} />
						</span>
						<span class="ml-3">{key}</span>
					</span>
				</li>
			{/each}
		</ul>
	</sidebar>
	<article>
		<table class="w-full">
			<thead>
				<tr>
					{#each keys as key}
						{#if visibleKeys[key]}
							<th class="text-left px-2">{key}</th>
						{/if}
					{/each}
				</tr>
			</thead>
			<tbody>
			{#each filteredMessages as message}
				{#if message}
				<tr class={`message ${message.level}`}>
					{#each keys as key}
						{#if visibleKeys[key]}
							<td>
								<span class="value">
									{#if message.hasOwnProperty(key)}
										{message[key]}
									{/if}
								</span>
							</td>
						{/if}
					{/each}
				</tr>
				{/if}
			{/each}
			</tbody>
		</table>
	</article>
</main>

<style global>
	/* purgecss start ignore */
	@tailwind base;
	@tailwind components;
	/* purgecss end ignore */
	@tailwind utilities;

	html {
		font-family: system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,"Noto Sans",sans-serif,"Apple Color Emoji","Segoe UI Emoji","Segoe UI Symbol","Noto Color Emoji";
	}

	body {
		padding: 0;
		height: 100%;
		width: 100%;
	}

	main {
		padding: 0;
		margin: 0 auto;
		@apply flex flex-row;
	}

	sidebar {
		@apply w-64 h-screen flex p-2 bg-white;
	}

	article {
		overflow-y: scroll;
		/* height: 100%; */
		width: 100%;
		@apply h-screen bg-gray-800 text-gray-100 p-4;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}

	.message .key {
		@apply px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800;
	}
 
	.message .value {
		@apply px-2 text-sm;
	}

	.message.error .key {
		@apply bg-red-100 text-red-800;
	}

	.message.warn .key, .message.warning .key {
		@apply bg-orange-100 text-orange-800;
	}

	.message.info .key {
		@apply bg-green-100 text-green-800;
	}

	.message.debug .key {
		@apply bg-blue-100 text-blue-800;
	}

	.message.trace .key {
		@apply bg-indigo-100 text-indigo-800;
	}
</style>