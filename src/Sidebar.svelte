<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';

	interface Props {
		keys: string[];
		visibleKeys: Record<string, boolean>;
		searchTerm: string;
		totalMessages: number;
		filteredCount: number;
		callbacks: {
			applyFilters: () => void;
			selectAll: () => void;
			selectNone: () => void;
		};
		onToggle: () => void;
	}

	let {
		keys,
		visibleKeys = $bindable(),
		searchTerm = $bindable(),
		totalMessages,
		filteredCount,
		callbacks,
		onToggle,
	}: Props = $props();
</script>

<aside class="flex w-64 shrink-0 flex-col gap-4 border-r border-border bg-card p-4">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-2">
			<svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-foreground" viewBox="0 0 64 64" fill="currentColor">
				<g>
					<polygon points="53.414 37 51.414 35 39 35 39 37 50.586 37 52.586 39 58 39 58 37 53.414 37"/>
					<polygon points="32.586 42 24 42 24 44 33.414 44 35.414 42 41 42 41 40 34.586 40 32.586 42"/>
					<rect x="33" y="30" width="4" height="2"/>
					<path d="M11,30c-2.8,0-5,3.075-5,7s2.2,7,5,7,5-3.075,5-7S13.8,30,11,30Zm0,12c-1.626,0-3-2.29-3-5s1.374-5,3-5,3,2.29,3,5S12.626,42,11,42Z"/>
					<path d="M53,26H50.109a6.433,6.433,0,0,1,1.331-3H55a5.006,5.006,0,0,0,5-5V17H55a4.99,4.99,0,0,0-4.956,4.565A8.545,8.545,0,0,0,48.086,26H43.721l-2-6H34.279l-2,6H25.118a6.284,6.284,0,0,0-3.181-5.62A4.989,4.989,0,0,0,17,16H12v1a5.006,5.006,0,0,0,5,5h3.738a4.284,4.284,0,0,1,2.373,4H11c-4.962,0-9,4.935-9,11s4.038,11,9,11H53c4.962,0,9-4.935,9-11S57.962,26,53,26Zm2-7h2.829A3.006,3.006,0,0,1,55,21H52.171A3.006,3.006,0,0,1,55,19ZM17,20a3.006,3.006,0,0,1-2.829-2H17a3.006,3.006,0,0,1,2.829,2ZM4,37c0-4.963,3.14-9,7-9s7,4.037,7,9-3.14,9-7,9S4,41.963,4,37Zm55.261,4H46v2H58.2A6.434,6.434,0,0,1,53,46H16.158a11.664,11.664,0,0,0,3.8-8h8.455l2-2H34V34H29.586l-2,2H19.959a12.79,12.79,0,0,0-.953-4H29V30H17.937a10.07,10.07,0,0,0-1.779-2H33.721l2-6h4.558l1.334,4H38v2H53a6.1,6.1,0,0,1,4.39,2H43v2H58.816A10.833,10.833,0,0,1,60,37,11.048,11.048,0,0,1,59.261,41Z"/>
				</g>
			</svg>
			<span class="text-2xl font-light tracking-tight">Leño</span>
		</div>
		<button onclick={onToggle} class="text-muted-foreground hover:text-foreground" title="Hide sidebar">
			<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<polyline points="15 18 9 12 15 6" />
			</svg>
		</button>
	</div>

	<Input
		type="search"
		placeholder="Search logs..."
		bind:value={searchTerm}
		oninput={callbacks.applyFilters}
		class="w-full"
	/>

	{#if keys.length > 0}
		<div class="flex flex-col gap-1">
			<div class="flex items-center justify-between">
				<p class="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Fields</p>
				<div class="flex gap-1">
					<button onclick={callbacks.selectAll} class="text-xs text-muted-foreground hover:text-foreground">All</button>
					<span class="text-xs text-muted-foreground">/</span>
					<button onclick={callbacks.selectNone} class="text-xs text-muted-foreground hover:text-foreground">None</button>
				</div>
			</div>
			{#each keys as key}
				<label class="flex cursor-pointer items-center gap-2 rounded px-1 py-0.5 hover:bg-accent">
					<Checkbox bind:checked={visibleKeys[key]} />
					<span class="text-sm">{key}</span>
				</label>
			{/each}
		</div>
	{/if}

	<div class="mt-auto text-xs text-muted-foreground">
		{filteredCount} / {totalMessages} messages
	</div>
</aside>
