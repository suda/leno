<script lang="ts">
	import store, { type LogMessage } from './store.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import Sidebar from './Sidebar.svelte';

	const STORAGE_KEY = 'leno:visibleKeys';
	const FILTERS_STORAGE_KEY = 'leno:fieldFilters';
	const DARK_MODE_KEY = 'leno:darkMode';

	let messages = $state<LogMessage[]>([]);
	let filteredMessages = $state<LogMessage[]>([]);
	let keys = $state<string[]>([]);
	let visibleKeys = $state<Record<string, boolean>>(
		JSON.parse(localStorage.getItem(STORAGE_KEY) ?? '{}')
	);
	let searchTerm = $state('');
	let sidebarVisible = $state(true);
	let darkMode = $state(localStorage.getItem(DARK_MODE_KEY) === 'true');
	// fieldFilters: map of field name → set of selected values (null means "all")
	let fieldFilters = $state<Record<string, string[]>>(
		JSON.parse(localStorage.getItem(FILTERS_STORAGE_KEY) ?? '{}')
	);

	$effect(() => {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(visibleKeys));
	});

	$effect(() => {
		localStorage.setItem(DARK_MODE_KEY, String(darkMode));
		document.documentElement.classList.toggle('dark', darkMode);
	});

	$effect(() => {
		localStorage.setItem(FILTERS_STORAGE_KEY, JSON.stringify(fieldFilters));
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

	function getTopValues(field: string): string[] {
		const counts = new Map<string, number>();
		for (const msg of messages) {
			if (field in msg && msg[field] !== undefined) {
				const val = String(msg[field]);
				if (val.trim() === '') continue;
				counts.set(val, (counts.get(val) ?? 0) + 1);
			}
		}
		return [...counts.entries()]
			.sort((a, b) => b[1] - a[1])
			.slice(0, 5)
			.map(([val]) => val);
	}

	function addFilter(field: string) {
		if (field in fieldFilters) return;
		const topValues = getTopValues(field);
		fieldFilters = { ...fieldFilters, [field]: topValues };
		applyFilters();
	}

	function removeFilter(field: string) {
		const { [field]: _, ...rest } = fieldFilters;
		fieldFilters = rest;
		applyFilters();
	}

	function filterMessage(currentMessage: LogMessage): boolean {
		if (searchTerm !== '') {
			const matchesSearch = keys.some(
				(key) =>
					key in currentMessage &&
					String(currentMessage[key]).toLowerCase().includes(searchTerm.toLowerCase())
			);
			if (!matchesSearch) return false;
		}

		for (const [field, selectedValues] of Object.entries(fieldFilters)) {
			if (selectedValues.length === 0) continue;
			const msgVal = currentMessage[field] !== undefined ? String(currentMessage[field]) : undefined;
			if (msgVal === undefined || !selectedValues.includes(msgVal)) return false;
		}

		return true;
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
	{#if sidebarVisible}
		<Sidebar
			{keys}
			{messages}
			bind:visibleKeys
			bind:searchTerm
			bind:fieldFilters
			totalMessages={messages.length}
			filteredCount={filteredMessages.length}
			callbacks={{ applyFilters, selectAll, selectNone, addFilter, removeFilter, getTopValues }}
			onToggle={() => (sidebarVisible = false)}
			{darkMode}
			onToggleDarkMode={() => (darkMode = !darkMode)}
		/>
	{:else}
		<button
			onclick={() => (sidebarVisible = true)}
			class="flex h-full w-8 shrink-0 items-start justify-center border-r border-border bg-card pt-3 text-muted-foreground hover:text-foreground"
			title="Show sidebar"
		>
			<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<polyline points="9 18 15 12 9 6" />
			</svg>
		</button>
	{/if}

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
									? 'bg-sky-500/15 hover:bg-sky-500/25'
									: level === 'debug'
										? 'bg-violet-500/20 hover:bg-violet-500/30'
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
