<template>
	<div class='feed'>
		<div class='results' v-if='results'>
			<ListView class='result' v-for='result in results' v-bind:key='result.id' :entry='result' />
		</div>
		<div v-else>No registered records. Maybe explore for something new?</div>
		<Pagination :onChange='handlePage' :settings='paging()' />
	</div>
</template>

<script lang='ts'>
import { Component, Prop, Vue } from 'vue-property-decorator';
import ListView from './ListView.vue';
import Pagination from './Pagination.vue';

@Component({
	components: {
		ListView,
		Pagination,
	}
})
export default class Feed extends Vue {
	@Prop() private results!: {};
	@Prop() private pagination!: {};

	private data() {
		return {
			paging: () => {
				return this.pagination;
			}
		}
	}

	private handlePage(i: any) {
		switch (i) {
		case "previous":
			this.$store.dispatch("setSearchPage", this.$store.state.searchPagination.page-1);
			break;
		case "next":
			this.$store.dispatch("setSearchPage", this.$store.state.searchPagination.page+1);
			break;
		default:
			this.$store.dispatch("setSearchPage", i);
			break;
		}
	}
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped lang='scss'>
.feed {
	max-width: 1335px;
	margin: 0 auto;

	.results {
		display: flex;
		flex-flow: row wrap;
		justify-content: flex-start;

		.result {
			height: 350px;
			flex-basis: 20%;
			width: 259px;
			position: relative;
			padding: 10px;
			box-sizing: border-box;
			margin: 10px;
		}
	}
}
</style>
