<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';

	interface Props {
		keys: string[];
		visibleKeys: Record<string, boolean>;
		searchTerm: string;
		fieldFilters: Record<string, string[]>;
		topValuesCache: Record<string, string[]>;
		totalMessages: number;
		filteredCount: number;
		callbacks: {
			applyFilters: () => void;
			selectAll: () => void;
			selectNone: () => void;
			addFilter: (field: string) => void;
			removeFilter: (field: string) => void;
		};
		onToggle: () => void;
		darkMode: boolean;
		onToggleDarkMode: () => void;
	}

	let {
		keys,
		visibleKeys = $bindable(),
		searchTerm = $bindable(),
		fieldFilters = $bindable(),
		topValuesCache,
		totalMessages,
		filteredCount,
		callbacks,
		onToggle,
		darkMode,
		onToggleDarkMode,
	}: Props = $props();

	let version = $state('');
	fetch('/version')
		.then((r) => r.text())
		.then((v) => (version = v));

	function toggleFilterValue(field: string, value: string) {
		const current = fieldFilters[field] ?? [];
		if (current.includes(value)) {
			fieldFilters = { ...fieldFilters, [field]: current.filter((v) => v !== value) };
		} else {
			fieldFilters = { ...fieldFilters, [field]: [...current, value] };
		}
		callbacks.applyFilters();
	}

	function selectAllValues(field: string) {
		fieldFilters = { ...fieldFilters, [field]: topValuesCache[field] ?? [] };
		callbacks.applyFilters();
	}

	function selectNoneValues(field: string) {
		fieldFilters = { ...fieldFilters, [field]: [] };
		callbacks.applyFilters();
	}
</script>

<aside
	class="flex w-64 shrink-0 flex-col border-r border-sidebar-border bg-sidebar text-sidebar-foreground overflow-hidden"
>
	<!-- Header -->
	<div class="flex h-14 shrink-0 items-center justify-between border-b border-sidebar-border px-4">
		<div class="flex items-center gap-2.5">
			<div
				class="flex h-7 w-7 items-center justify-center rounded-md bg-sidebar-primary text-sidebar-primary-foreground"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-4 w-4"
					viewBox="0 0 64 64"
					fill="currentColor"
				>
					<g>
						<polygon
							points="53.414 37 51.414 35 39 35 39 37 50.586 37 52.586 39 58 39 58 37 53.414 37"
						/>
						<polygon
							points="32.586 42 24 42 24 44 33.414 44 35.414 42 41 42 41 40 34.586 40 32.586 42"
						/>
						<rect x="33" y="30" width="4" height="2" />
						<path
							d="M11,30c-2.8,0-5,3.075-5,7s2.2,7,5,7,5-3.075,5-7S13.8,30,11,30Zm0,12c-1.626,0-3-2.29-3-5s1.374-5,3-5,3,2.29,3,5S12.626,42,11,42Z"
						/>
						<path
							d="M53,26H50.109a6.433,6.433,0,0,1,1.331-3H55a5.006,5.006,0,0,0,5-5V17H55a4.99,4.99,0,0,0-4.956,4.565A8.545,8.545,0,0,0,48.086,26H43.721l-2-6H34.279l-2,6H25.118a6.284,6.284,0,0,0-3.181-5.62A4.989,4.989,0,0,0,17,16H12v1a5.006,5.006,0,0,0,5,5h3.738a4.284,4.284,0,0,1,2.373,4H11c-4.962,0-9,4.935-9,11s4.038,11,9,11H53c4.962,0,9-4.935,9-11S57.962,26,53,26Zm2-7h2.829A3.006,3.006,0,0,1,55,21H52.171A3.006,3.006,0,0,1,55,19ZM17,20a3.006,3.006,0,0,1-2.829-2H17a3.006,3.006,0,0,1,2.829,2ZM4,37c0-4.963,3.14-9,7-9s7,4.037,7,9-3.14,9-7,9S4,41.963,4,37Zm55.261,4H46v2H58.2A6.434,6.434,0,0,1,53,46H16.158a11.664,11.664,0,0,0,3.8-8h8.455l2-2H34V34H29.586l-2,2H19.959a12.79,12.79,0,0,0-.953-4H29V30H17.937a10.07,10.07,0,0,0-1.779-2H33.721l2-6h4.558l1.334,4H38v2H53a6.1,6.1,0,0,1,4.39,2H43v2H58.816A10.833,10.833,0,0,1,60,37,11.048,11.048,0,0,1,59.261,41Z"
						/>
					</g>
				</svg>
			</div>
			<span class="text-sm font-semibold tracking-tight">Leño</span>
		</div>
		<div class="flex items-center gap-0.5">
			<button
				onclick={onToggleDarkMode}
				class="flex h-7 w-7 items-center justify-center rounded-md text-sidebar-foreground/60 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground transition-colors"
				title={darkMode ? 'Switch to light mode' : 'Switch to dark mode'}
			>
				{#if darkMode}
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
						<circle cx="12" cy="12" r="4" />
						<path
							d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"
						/>
					</svg>
				{:else}
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
						<path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
					</svg>
				{/if}
			</button>
			<button
				onclick={onToggle}
				class="flex h-7 w-7 items-center justify-center rounded-md text-sidebar-foreground/60 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground transition-colors"
				title="Hide sidebar"
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
					<polyline points="15 18 9 12 15 6" />
				</svg>
			</button>
		</div>
	</div>

	<!-- Scrollable content -->
	<div class="flex flex-1 flex-col gap-0 overflow-y-auto">
		<!-- Search -->
		<div class="px-3 py-3 border-b border-sidebar-border">
			<div class="relative">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-sidebar-foreground/40 pointer-events-none"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<circle cx="11" cy="11" r="8" />
					<path d="m21 21-4.35-4.35" />
				</svg>
				<Input
					type="search"
					placeholder="Search logs..."
					bind:value={searchTerm}
					oninput={callbacks.applyFilters}
					class="w-full pl-8 h-8 text-sm bg-sidebar-accent/50 border-sidebar-border focus-visible:ring-sidebar-ring"
				/>
			</div>
		</div>

		{#if keys.length > 0}
			<!-- Fields section -->
			<div class="px-3 py-3 border-b border-sidebar-border">
				<div class="flex items-center justify-between mb-2">
					<span class="text-xs font-medium text-sidebar-foreground/50 uppercase tracking-wider"
						>Columns</span
					>
					<div class="flex items-center gap-1">
						<button
							onclick={callbacks.selectAll}
							class="text-xs text-sidebar-foreground/50 hover:text-sidebar-foreground transition-colors px-1 rounded hover:bg-sidebar-accent"
							>All</button
						>
						<span class="text-xs text-sidebar-foreground/30">/</span>
						<button
							onclick={callbacks.selectNone}
							class="text-xs text-sidebar-foreground/50 hover:text-sidebar-foreground transition-colors px-1 rounded hover:bg-sidebar-accent"
							>None</button
						>
					</div>
				</div>
				<div class="flex flex-col gap-0.5">
					{#each keys as key}
						<div
							class="group flex cursor-pointer items-center gap-2 rounded-md px-2 py-1 hover:bg-sidebar-accent transition-colors"
						>
							<label class="flex flex-1 cursor-pointer items-center gap-2 min-w-0">
								<Checkbox bind:checked={visibleKeys[key]} />
								<span class="text-sm truncate">{key}</span>
							</label>
							{#if !(key in fieldFilters)}
								<button
									onclick={() => callbacks.addFilter(key)}
									class="invisible shrink-0 text-sidebar-foreground/40 hover:text-sidebar-foreground group-hover:visible transition-colors"
									title="Filter by {key}"
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										class="h-3.5 w-3.5"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
									</svg>
								</button>
							{:else}
								<span class="shrink-0 text-sidebar-primary" title="Filter active">
									<svg
										xmlns="http://www.w3.org/2000/svg"
										class="h-3.5 w-3.5"
										viewBox="0 0 24 24"
										fill="currentColor"
										stroke="currentColor"
										stroke-width="1"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
									</svg>
								</span>
							{/if}
						</div>
					{/each}
				</div>
			</div>
		{/if}

		{#if Object.keys(fieldFilters).length > 0}
			<!-- Filters section -->
			<div class="px-3 py-3">
				<div class="mb-2">
					<span class="text-xs font-medium text-sidebar-foreground/50 uppercase tracking-wider"
						>Active Filters</span
					>
				</div>
				<div class="flex flex-col gap-2">
					{#each Object.entries(fieldFilters) as [field, selectedValues]}
						{@const availableValues = topValuesCache[field] ?? []}
						<div
							class="rounded-lg border border-sidebar-border bg-sidebar-accent/30 overflow-hidden"
						>
							<div class="flex items-center justify-between px-2.5 py-1.5 bg-sidebar-accent/50">
								<span class="text-xs font-medium">{field}</span>
								<div class="flex items-center gap-1">
									<button
										onclick={() => selectAllValues(field)}
										class="text-xs text-sidebar-foreground/50 hover:text-sidebar-foreground transition-colors"
										>All</button
									>
									<span class="text-xs text-sidebar-foreground/30">/</span>
									<button
										onclick={() => selectNoneValues(field)}
										class="text-xs text-sidebar-foreground/50 hover:text-sidebar-foreground transition-colors"
										>None</button
									>
									<button
										onclick={() => callbacks.removeFilter(field)}
										class="ml-0.5 text-sidebar-foreground/40 hover:text-destructive transition-colors"
										title="Remove filter"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											class="h-3.5 w-3.5"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											stroke-width="2"
											stroke-linecap="round"
											stroke-linejoin="round"
										>
											<line x1="18" y1="6" x2="6" y2="18" />
											<line x1="6" y1="6" x2="18" y2="18" />
										</svg>
									</button>
								</div>
							</div>
							<div class="px-2 py-1.5 flex flex-col gap-0.5">
								{#each availableValues as value}
									<label
										class="flex cursor-pointer items-center gap-2 rounded px-1 py-0.5 hover:bg-sidebar-accent transition-colors"
									>
										<Checkbox
											checked={selectedValues.includes(value)}
											onCheckedChange={() => toggleFilterValue(field, value)}
										/>
										<span class="text-xs truncate text-sidebar-foreground/80" title={value}
											>{value}</span
										>
									</label>
								{/each}
								{#if availableValues.length === 0}
									<p class="text-xs text-sidebar-foreground/40 italic px-1 py-0.5">No values yet</p>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>

	<!-- Footer -->
	<div class="shrink-0 border-t border-sidebar-border px-4 py-3">
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-1.5">
				<div class="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse"></div>
				<span class="text-xs text-sidebar-foreground/60">
					{filteredCount.toLocaleString()} / {totalMessages.toLocaleString()}
				</span>
			</div>
			{#if version}
				<div class="flex items-center gap-1.5 text-xs text-sidebar-foreground/40">
					<span>v{version}</span>
					<span>·</span>
					<a
						href="https://github.com/suda/leno"
						target="_blank"
						rel="noopener noreferrer"
						class="hover:text-sidebar-foreground transition-colors">GitHub</a
					>
				</div>
			{/if}
		</div>
	</div>
</aside>
