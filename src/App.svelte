<script lang="ts">
	import store, { type LogMessage } from './store.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import Sidebar from './Sidebar.svelte';

	const STORAGE_KEY = 'leno:visibleKeys';
	const FILTERS_STORAGE_KEY = 'leno:fieldFilters';
	const DARK_MODE_KEY = 'leno:darkMode';
	const MAX_DISPLAY = 2000;

	let messages = $state<LogMessage[]>([]);
	let filteredMessages = $state<LogMessage[]>([]);
	let keys = $state<string[]>([]);
	let keysSet = new Set<string>();
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

	// Memoized top values: recomputes when messages array or fieldFilters keys change
	const topValuesCache = $derived.by(() => {
		const result: Record<string, string[]> = {};
		for (const field of Object.keys(fieldFilters)) {
			result[field] = computeTopValues(field);
		}
		return result;
	});

	function computeTopValues(field: string): string[] {
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

	function addKeys(newKeys: string[]) {
		for (const k of newKeys) {
			if (!keysSet.has(k)) {
				keysSet.add(k);
				keys.push(k);
			}
			if (!(k in visibleKeys)) visibleKeys[k] = true;
		}
	}

	function selectAll() {
		for (const key of keys) visibleKeys[key] = true;
	}

	function selectNone() {
		for (const key of keys) visibleKeys[key] = false;
	}

	function addFilter(field: string) {
		if (field in fieldFilters) return;
		const topValues = computeTopValues(field);
		fieldFilters = { ...fieldFilters, [field]: topValues };
		applyFilters();
	}

	function removeFilter(field: string) {
		const { [field]: _removed, ...rest } = fieldFilters;
		fieldFilters = rest;
		applyFilters();
	}

	function filterMessage(currentMessage: LogMessage): boolean {
		if (searchTerm !== '') {
			const lowerSearch = searchTerm.toLowerCase();
			const matchesSearch = keys.some(
				(key) =>
					key in currentMessage && String(currentMessage[key]).toLowerCase().includes(lowerSearch)
			);
			if (!matchesSearch) return false;
		}

		for (const [field, selectedValues] of Object.entries(fieldFilters)) {
			if (selectedValues.length === 0) continue;
			const msgVal =
				currentMessage[field] !== undefined ? String(currentMessage[field]) : undefined;
			if (msgVal === undefined || !selectedValues.includes(msgVal)) return false;
		}

		return true;
	}

	// RAF batching: buffer messages and flush at most once per animation frame
	let pendingMessages: LogMessage[] = [];
	let rafPending = false;

	function scheduleFlush() {
		if (rafPending) return;
		rafPending = true;
		requestAnimationFrame(() => {
			rafPending = false;
			const batch = pendingMessages.splice(0);
			if (batch.length === 0) return;

			for (const msg of batch) {
				addKeys(Object.keys(msg));
			}

			// Prepend batch to messages in one mutation
			messages.unshift(...batch);

			// Prepend matching messages to filteredMessages in one mutation
			const matching = batch.filter(filterMessage);
			if (matching.length > 0) {
				filteredMessages.unshift(...matching);
				// Trim to display cap
				if (filteredMessages.length > MAX_DISPLAY) {
					filteredMessages.length = MAX_DISPLAY;
				}
			}
		});
	}

	function queueMessage(currentMessage: LogMessage | null) {
		if (!currentMessage) return;
		pendingMessages.push(currentMessage);
		scheduleFlush();
	}

	function applyFilters() {
		const result = messages.filter(filterMessage);
		// Replace array contents in place to avoid signal identity change triggering full reconcile
		filteredMessages.length = 0;
		const capped = result.slice(0, MAX_DISPLAY);
		filteredMessages.push(...capped);
	}

	$effect(() => {
		return store.subscribe(queueMessage);
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

	function getLevelClass(level: string): string {
		switch (level) {
			case 'error':
			case 'fatal':
				return 'bg-destructive/8 hover:bg-destructive/15 border-l-2 border-l-destructive/50';
			case 'warn':
			case 'warning':
				return 'bg-yellow-500/8 hover:bg-yellow-500/15 border-l-2 border-l-yellow-500/50';
			case 'info':
				return 'bg-sky-500/8 hover:bg-sky-500/15 border-l-2 border-l-sky-500/40';
			case 'debug':
				return 'bg-violet-500/8 hover:bg-violet-500/15 border-l-2 border-l-violet-500/40';
			case 'trace':
				return 'bg-cyan-500/8 hover:bg-cyan-500/15 border-l-2 border-l-cyan-500/40';
			default:
				return 'hover:bg-muted/50';
		}
	}
</script>

<div class="flex h-screen w-full overflow-hidden bg-background">
	{#if sidebarVisible}
		<Sidebar
			{keys}
			bind:visibleKeys
			bind:searchTerm
			bind:fieldFilters
			{topValuesCache}
			totalMessages={messages.length}
			filteredCount={filteredMessages.length}
			callbacks={{ applyFilters, selectAll, selectNone, addFilter, removeFilter }}
			onToggle={() => (sidebarVisible = false)}
			{darkMode}
			onToggleDarkMode={() => (darkMode = !darkMode)}
		/>
	{:else}
		<button
			onclick={() => (sidebarVisible = true)}
			class="flex h-full w-9 shrink-0 items-start justify-center border-r border-border bg-sidebar pt-[17px] text-sidebar-foreground/40 hover:text-sidebar-foreground transition-colors"
			title="Show sidebar"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-4 w-4"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<polyline points="9 18 15 12 9 6" />
			</svg>
		</button>
	{/if}

	<!-- Main content -->
	<div class="flex flex-1 flex-col overflow-hidden">
		<!-- Top bar -->
		<header
			class="flex h-14 shrink-0 items-center justify-between border-b border-border bg-background px-4"
		>
			<div class="flex items-center gap-3">
				<div class="flex items-center gap-1.5">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-4 w-4 text-muted-foreground"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
						<polyline points="14 2 14 8 20 8" />
						<line x1="16" y1="13" x2="8" y2="13" />
						<line x1="16" y1="17" x2="8" y2="17" />
						<polyline points="10 9 9 9 8 9" />
					</svg>
					<h1 class="text-sm font-medium">Log Stream</h1>
				</div>
				{#if messages.length > 0}
					<div class="flex items-center gap-1.5 rounded-md bg-muted/60 px-2 py-0.5">
						<div class="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse"></div>
						<span class="text-xs text-muted-foreground font-mono">
							{filteredMessages.length.toLocaleString()}
							{#if filteredMessages.length !== messages.length}
								<span class="text-muted-foreground/50">/ {messages.length.toLocaleString()}</span>
							{/if}
						</span>
					</div>
				{/if}
			</div>
			{#if searchTerm || Object.keys(fieldFilters).length > 0}
				<div class="flex items-center gap-2 text-xs text-muted-foreground">
					{#if searchTerm}
						<span class="flex items-center gap-1 rounded-md bg-muted px-2 py-0.5">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-3 w-3"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
							>
								<circle cx="11" cy="11" r="8" /><path d="m21 21-4.35-4.35" />
							</svg>
							"{searchTerm}"
						</span>
					{/if}
					{#if Object.keys(fieldFilters).length > 0}
						<span class="flex items-center gap-1 rounded-md bg-muted px-2 py-0.5">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-3 w-3"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
							>
								<polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
							</svg>
							{Object.keys(fieldFilters).length} filter{Object.keys(fieldFilters).length !== 1
								? 's'
								: ''}
						</span>
					{/if}
				</div>
			{/if}
		</header>

		<!-- Table area -->
		<main class="flex-1 overflow-auto">
			{#if messages.length === 0}
				<div class="flex h-full flex-col items-center justify-center gap-4 text-muted-foreground">
					<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-muted">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6 text-muted-foreground/60"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="1.5"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
							<polyline points="14 2 14 8 20 8" />
							<line x1="16" y1="13" x2="8" y2="13" />
							<line x1="16" y1="17" x2="8" y2="17" />
							<polyline points="10 9 9 9 8 9" />
						</svg>
					</div>
					<div class="text-center">
						<p class="text-sm font-medium text-foreground">Waiting for log messages</p>
						<p class="mt-1 text-xs text-muted-foreground">
							Pipe JSON logs to this process to get started
						</p>
					</div>
					<div class="rounded-md border border-border bg-muted/50 px-3 py-2">
						<code class="text-xs font-mono text-muted-foreground">./myapp | leno</code>
					</div>
				</div>
			{:else}
				<Table.Root>
					<Table.Header
						class="sticky top-0 z-10 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80"
					>
						<Table.Row class="border-b border-border hover:bg-transparent">
							{#each keys as key}
								{#if visibleKeys[key]}
									<Table.Head
										class="whitespace-nowrap h-10 text-xs font-medium text-muted-foreground uppercase tracking-wide px-3"
										>{key}</Table.Head
									>
								{/if}
							{/each}
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each filteredMessages as message}
							{@const level = String(message.level ?? '').toLowerCase()}
							<Table.Row class="border-b border-border/50 {getLevelClass(level)} transition-colors">
								{#each keys as key}
									{#if visibleKeys[key]}
										<Table.Cell class="py-1.5 px-3 text-xs font-mono">
											{#if key === 'level' && message[key] !== undefined}
												<Badge
													variant={getLevelVariant(message[key])}
													class="text-xs font-medium px-1.5 py-0"
												>
													{message[key]}
												</Badge>
											{:else if message[key] !== undefined}
												<span class="text-foreground/80">{message[key]}</span>
											{/if}
										</Table.Cell>
									{/if}
								{/each}
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			{/if}
		</main>
	</div>
</div>
