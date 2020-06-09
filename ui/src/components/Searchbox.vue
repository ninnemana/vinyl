<template>
	<form class='search' @submit.prevent='handleSubmit'>
		<input type='search' v-model='term' :placeholder='label' />
	</form>
</template>

<script lang='ts'>
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Searchbox extends Vue {
	@Prop() private query!: string;
	private term!: string;
	private label!: string;

	private handleSubmit() {
		this.$store.dispatch('search', this.term);
		this.$router.push(`/search?query=${this.term}`);
	}

	private data() {
		return { 
			term: '',
			label: 'Discover something new ...',
		};
	}
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped lang='scss'>
.search{	
	input[type=search]{
		border: none;
		background: none;
		border-bottom: .4rem solid gray;
		padding: .5rem .2rem;
		margin: 1rem 0;
		width: 100%;
		font-size: 2rem;

		&:focus{
			outline: none;
		}
	}
}
</style>
