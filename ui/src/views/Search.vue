<template>
	<div class='search'>
		<Searchbox />
		<Feed :results='this.$store.state.searchResults' :pagination='this.$store.state.searchPagination' />
	</div>
</template>

<script>
// @ is an alias to /src
import Feed from '@/components/Feed.vue';
import Searchbox from '@/components/Searchbox.vue';

export default {
	name: 'Search',
	components: {
		Searchbox,
		Feed
	}
};
</script>
