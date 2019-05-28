<template>
	<div class="home">
		<img alt="VinylTap Logo" src="../assets/logo.png">
		<Search />
		<Feed />
	</div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Feed from '@/components/Feed.vue';
import Search from '@/components/Search.vue';

@Component({
	components: {
		Feed,
		Search,
	},
})
export default class Home extends Vue {}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.home {
	width: 80%;
	margin: auto;
}
</style>
