<template>
	<v-container>
		<img alt='VinylTap Logo' src='../assets/logo.png' />
		<v-card class='search'>
			<Searchbox />
		</v-card>
	</v-container>
</template>

<script>
// @ is an alias to /src
import Searchbox from '@/components/Searchbox.vue';

export default {
	name: 'Home',
	components: {
		Searchbox
	}
};
</script>
<style lang='scss'>
.search {
	margin: 2em;
}
</style>
