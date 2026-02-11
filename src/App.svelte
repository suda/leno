<script lang="ts">
	import store, { type LogMessage } from './store.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import Sidebar from './Sidebar.svelte';

	const STORAGE_KEY = 'leno:visibleKeys';

	let messages = $state<LogMessage[]>([]);
	let filteredMessages = $state<LogMessage[]>([]);
	let keys = $state<string[]>([]);
	let visibleKeys = $state<Record<string, boolean>>(
		JSON.parse(localStorage.getItem(STORAGE_KEY) ?? '{}')
	);
	let searchTerm = $state('');

	$effect(() => {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(visibleKeys));
	});

	function addKeys(newKeys: string[]) {
		keys = [...new Set([...keys, ...newKeys])];
		for (const key of keys) {
			if (!(key in visibleKeys)) visibleKeys[key] = true;
		}
	}

	function selectAll() {
		for (const key of keys) visibleKeys[key] = true;
	}

	function selectNone() {
		for (const key of keys) visibleKeys[key] = false;
	}

	function filterMessage(currentMessage: LogMessage): boolean {
		if (searchTerm === '') return true;
		return keys.some(
			(key) =>
				key in currentMessage &&
				String(currentMessage[key]).toLowerCase().includes(searchTerm.toLowerCase())
		);
	}

	function processMessage(currentMessage: LogMessage | null) {
		if (!currentMessage) return;
		messages = [currentMessage, ...messages];
		addKeys(Object.keys(currentMessage));
		if (filterMessage(currentMessage)) {
			filteredMessages = [currentMessage, ...filteredMessages];
		}
	}

	function applyFilters() {
		filteredMessages = messages.filter(filterMessage);
	}

	$effect(() => {
		return store.subscribe(processMessage);
	});

	function getLevelVariant(level: unknown): 'default' | 'destructive' | 'secondary' | 'outline' {
		switch (String(level).toLowerCase()) {
			case 'error':
			case 'fatal':
				return 'destructive';
			case 'warn':
			case 'warning':
				return 'secondary';
			default:
				return 'outline';
		}
	}
</script>

<div class="flex h-screen w-full overflow-hidden bg-background">
	<Sidebar
		{keys}
		bind:visibleKeys
		bind:searchTerm
		totalMessages={messages.length}
		filteredCount={filteredMessages.length}
		callbacks={{ applyFilters, selectAll, selectNone }}
	/>

	<!-- Main content -->
	<main class="flex-1 overflow-auto">
		<Table.Root>
			<Table.Header class="sticky top-0 z-10 bg-muted/80 backdrop-blur">
				<Table.Row>
					{#each keys as key}
						{#if visibleKeys[key]}
							<Table.Head class="whitespace-nowrap">{key}</Table.Head>
						{/if}
					{/each}
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each filteredMessages as message, i}
					{@const level = String(message.level ?? '').toLowerCase()}
					<Table.Row
						class={level === 'error' || level === 'fatal'
							? 'bg-destructive/10 hover:bg-destructive/20'
							: level === 'warn' || level === 'warning'
								? 'bg-yellow-500/10 hover:bg-yellow-500/20'
								: level === 'info'
									? 'bg-blue-500/10 hover:bg-blue-500/20'
									: level === 'debug'
										? 'bg-violet-500/10 hover:bg-violet-500/20'
										: level === 'trace'
											? 'bg-cyan-500/10 hover:bg-cyan-500/20'
											: i % 2 === 0
												? 'bg-muted/30'
												: ''}
					>
						{#each keys as key}
							{#if visibleKeys[key]}
								<Table.Cell class="text-sm">
									{#if key === 'level' && message[key] !== undefined}
										<Badge variant={getLevelVariant(message[key])}>
											{message[key]}
										</Badge>
									{:else if message[key] !== undefined}
										{message[key]}
									{/if}
								</Table.Cell>
							{/if}
						{/each}
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>

		{#if messages.length === 0}
			<div class="flex h-64 flex-col items-center justify-center gap-2 text-muted-foreground">
				<p class="text-lg">Waiting for log messages...</p>
				<p class="text-sm">Pipe JSON logs to this process to get started.</p>
			</div>
		{/if}
	</main>
</div>
