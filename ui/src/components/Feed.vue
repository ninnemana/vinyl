<template>
	<div class="feed">
		<ListView v-for="result in results" v-bind:key="result.id" entry={result} />
		<div v-if="!results">No registered records. Maybe explore for something new?</div>
	</div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import ListView from '@/components/ListView.vue';

@Component
export default class Feed extends Vue {
	@Prop() private results!: {};
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.feed {
	background: #ffffff;
	border-radius: .3rem;
	margin: auto;
	display: flex;
	flex-direction: column;
	box-shadow: .7rem .7rem 2.5rem 1.5rem rgba(203, 197, 193, 1);

	> div {
		padding: 1rem;
		box-shadow: inset 0 0 .02rem rgba(203, 197, 193, 1);
	}

	a {
		color: #42b983;
	}
}
</style>
