<template>
	<div class="auth-link">
		<div v-if='this.$store.state.authorized'>
			<router-link to='/account'>My Account</router-link>
		</div>
		<div v-if='!this.$store.state.authorized'>
			<router-link to='/login'>Login</router-link>
		</div>
	</div>
</template>

<script lang='ts'>
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component({
	components: {}
})
export default class AuthorizationLink extends Vue {
	// @Prop() private results!: {};
	// @Prop() private pagination!: {};

	// private data() {
	// 	return {
	// 		paging: () => {
	// 			return this.pagination;
	// 		}
	// 	}
	// }

	// private handlePage(i: any) {
	// 	switch (i) {
	// 	case "previous":
	// 		this.$store.dispatch("setSearchPage", this.$store.state.searchPagination.page-1);
	// 		break;
	// 	case "next":
	// 		this.$store.dispatch("setSearchPage", this.$store.state.searchPagination.page+1);
	// 		break;
	// 	default:
	// 		this.$store.dispatch("setSearchPage", i);
	// 		break;
	// 	}
	// }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped lang='scss'>
</style>
